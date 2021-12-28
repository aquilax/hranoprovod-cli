package app

import (
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/aquilax/hranoprovod-cli/v2/parser"
	"github.com/aquilax/hranoprovod-cli/v2/reporter"
	"github.com/aquilax/hranoprovod-cli/v2/resolver"
	"github.com/aquilax/hranoprovod-cli/v2/shared"
)

// Hranoprovod is the main app type
type Hranoprovod struct {
	options *Options
	output  io.Writer
}

// NewHranoprovod creates new application
func NewHranoprovod(options *Options) Hranoprovod {
	return Hranoprovod{options, os.Stdout}
}

// Register generates report
func (hr Hranoprovod) Register(gc GlobalConfig, pc parser.Config, rc resolver.Config, rpc reporter.Config, fc FilterConfig) error {
	parser := parser.NewParser(pc)
	nl, err := loadDatabase(parser, gc.DbFileName)
	if err != nil {
		return err
	}
	err = resolver.NewResolver(nl, rc).Resolve()
	if err != nil {
		return err
	}
	r := reporter.NewRegReporter(rpc, nl, os.Stdout)
	return hr.walkNodes(gc.LogFileName, parser, getNodeFilter(fc), r)
}

// Register generates report
func (hr Hranoprovod) Print(gc GlobalConfig, pc parser.Config, rpc reporter.Config, fc FilterConfig) error {
	parser := parser.NewParser(pc)
	r := reporter.NewPrintReporter(rpc, os.Stdout)
	return hr.walkNodes(gc.LogFileName, parser, getNodeFilter(fc), r)
}

// Balance generates balance report
func (hr Hranoprovod) Balance(gc GlobalConfig, pc parser.Config, rc resolver.Config, rpc reporter.Config, fc FilterConfig) error {
	parser := parser.NewParser(pc)
	nl, err := loadDatabase(parser, gc.DbFileName)
	if err != nil {
		return err
	}
	resolver.NewResolver(nl, rc).Resolve()
	r := reporter.NewBalanceReporter(rpc, nl, os.Stdout)
	return hr.walkNodes(gc.LogFileName, parser, getNodeFilter(fc), r)
}

// Lint lints file
func Lint(fileName string, pc parser.Config) error {
	p := parser.NewParser(pc)
	go p.ParseFile(fileName)
	return func() error {
		for {
			select {
			case <-p.Nodes:
			case err := <-p.Errors:
				return err
			case <-p.Done:
				return nil
			}
		}
	}()
}

// ReportElement generates report for single element
func ReportElement(dbFileName string, elementName string, ascending bool, pc parser.Config, rc resolver.Config, rpc reporter.Config) error {
	p := parser.NewParser(pc)
	nl, err := loadDatabase(p, dbFileName)
	if err != nil {
		return err
	}
	resolver.NewResolver(nl, rc).Resolve()
	var list []shared.Element
	for name, node := range nl {
		for _, el := range node.Elements {
			if el.Name == elementName {
				list = append(list, shared.NewElement(name, el.Value))
			}
		}
	}
	if ascending {
		sort.SliceStable(list, func(i, j int) bool {
			return list[i].Value > list[j].Value
		})
	} else {
		sort.SliceStable(list, func(i, j int) bool {
			return list[i].Value < list[j].Value
		})
	}
	return reporter.NewElementReporter(rpc, list).Flush()
}

// ReportQuantity Generates a quantity report
func (hr Hranoprovod) ReportQuantity(gc GlobalConfig, ascending bool, pc parser.Config, fc FilterConfig) error {
	parser := parser.NewParser(pc)
	r := reporter.NewQuantityReporter(ascending, os.Stdout)
	return hr.walkNodes(gc.LogFileName, parser, getNodeFilter(fc), r)
}

// ReportUnresolved generates report for unresolved elements
func (hr Hranoprovod) ReportUnresolved(gc GlobalConfig, pc parser.Config, rc resolver.Config, rpc reporter.Config, fc FilterConfig) error {
	parser := parser.NewParser(pc)
	nl, err := loadDatabase(parser, gc.DbFileName)
	if err != nil {
		return err
	}
	resolver.NewResolver(nl, rc).Resolve()
	r := reporter.NewUnsolvedReporter(rpc, nl, os.Stdout)

	return hr.walkNodes(gc.LogFileName, parser, getNodeFilter(fc), r)
}

// CSV generates CSV export
func (hr Hranoprovod) CSV(gc GlobalConfig, pc parser.Config, rpc reporter.Config, fc FilterConfig) error {
	parser := parser.NewParser(pc)
	r := reporter.NewCSVReporter(rpc, os.Stdout)
	return hr.walkNodes(gc.LogFileName, parser, getNodeFilter(fc), r)
}

// Summary generates summary
func (hr Hranoprovod) Summary(gc GlobalConfig, pc parser.Config, rc resolver.Config, rpc reporter.Config, fc FilterConfig) error {
	parser := parser.NewParser(pc)
	nl, err := loadDatabase(parser, gc.DbFileName)
	if err != nil {
		return err
	}
	resolver.NewResolver(nl, rc).Resolve()
	r := reporter.NewSummaryReporterTemplate(rpc, nl, os.Stdout)
	return hr.walkNodes(gc.LogFileName, parser, getNodeFilter(fc), r)
}

// Stats generates statistics report
func (hr Hranoprovod) Stats(gc GlobalConfig, pc parser.Config, rpc reporter.Config) error {
	var err error
	fLog, err := os.Open(gc.LogFileName)
	if err != nil {
		return parser.NewErrorIO(err, gc.LogFileName)
	}
	defer fLog.Close()

	fDb, err := os.Open(gc.DbFileName)
	if err != nil {
		return parser.NewErrorIO(err, gc.DbFileName)
	}
	defer fDb.Close()

	countLog := 0
	var firstLogDate time.Time
	var lastLogDate time.Time
	parser.ParseStreamCallback(fLog, pc, func(n *shared.ParserNode, _ error) (stop bool) {
		lastLogDate, err = time.Parse(gc.DateFormat, n.Header)
		if err == nil {
			if firstLogDate.IsZero() {
				firstLogDate = lastLogDate
			}
		}
		countLog++
		return false
	})

	countDb := 0
	parser.ParseStreamCallback(fDb, pc, func(n *shared.ParserNode, err error) (stop bool) {
		countDb++
		return false
	})
	return reporter.NewStatsReporter(rpc, []string{
		fmt.Sprintf("  Database file:      %s\n", gc.DbFileName),
		fmt.Sprintf("  Database records:   %d\n", countDb),
		fmt.Sprintln(""),
		fmt.Sprintf("  Log file:           %s\n", gc.LogFileName),
		fmt.Sprintf("  Log records:        %d\n", countLog),
		fmt.Sprintf("  First record:       %s (%d days ago)\n", firstLogDate.Format(rpc.DateFormat), int(time.Since(firstLogDate).Hours()/24)),
		fmt.Sprintf("  Last record:        %s (%d days ago)\n", lastLogDate.Format(rpc.DateFormat), int(time.Since(lastLogDate).Hours()/24)),
	}).Flush()
}

func loadDatabase(p parser.Parser, fileName string) (shared.DBNodeList, error) {
	nodeList := shared.NewDBNodeList()
	// Database must be optional. If the default file name is used and the file is not found,
	// return empty node list
	if fileName == DefaultDbFilename {
		if exists, _ := fileExists(fileName); !exists {
			return nodeList, nil
		}
	}
	go p.ParseFile(fileName)
	return func() (shared.DBNodeList, error) {
		for {
			select {
			case node := <-p.Nodes:
				nodeList.Push(shared.NewDBNodeFromNode(node))
			case error := <-p.Errors:
				return nodeList, error
			case <-p.Done:
				return nodeList, nil
			}
		}
	}()
}

type NodeFilter = func(t time.Time, node *shared.ParserNode) (bool, error)

func getNodeFilter(fc FilterConfig) *NodeFilter {
	if fc.BeginningTime == nil && fc.EndTime == nil {
		// no filter if beginning and end time are nil
		return nil
	}

	inInterval := func(t time.Time) bool {
		if (fc.BeginningTime != nil && !isGoodDate(t, *fc.BeginningTime, dateBeginning)) || (fc.EndTime != nil && !isGoodDate(t, *fc.EndTime, dateEnd)) {
			return false
		}
		return true
	}

	filter := func(t time.Time, node *shared.ParserNode) (bool, error) {
		return inInterval(t), nil
	}
	return &filter
}

func (hr Hranoprovod) walkNodes(logFileName string, p parser.Parser, filter *NodeFilter, r reporter.Reporter) error {
	var node *shared.ParserNode
	var ln *shared.LogNode
	var err error
	var t time.Time
	var ok bool

	go p.ParseFile(logFileName)
	for {
		select {
		case node = <-p.Nodes:
			t, err = time.Parse(hr.options.GlobalConfig.DateFormat, node.Header)
			if err != nil {
				return err
			}
			ok = true
			if filter != nil {
				ok, err = (*filter)(t, node)
				if err != nil {
					return err
				}
			}
			if ok {
				ln, err = shared.NewLogNodeFromElements(t, node.Elements, node.Metadata)
				if err != nil {
					return err
				}
				r.Process(ln)
			}
		case err = <-p.Errors:
			return err
		case <-p.Done:
			r.Flush()
			return nil
		}
	}
}

// compareType identifies the type of date comparison
type compareType bool

const (
	dateBeginning compareType = true
	dateEnd       compareType = false
)

func isGoodDate(time, compareTime time.Time, ct compareType) bool {
	if time.Equal(compareTime) {
		return true
	}
	if ct == dateBeginning {
		return time.After(compareTime)
	}
	return time.Before(compareTime)
}
