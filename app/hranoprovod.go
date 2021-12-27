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
func (hr Hranoprovod) Register(pc parser.Config) error {
	parser := parser.NewParser(pc)
	nl, err := loadDatabase(parser, hr.options.Global.DbFileName)
	if err != nil {
		return err
	}
	resolver.NewResolver(nl, hr.options.Resolver.ResolverMaxDepth).Resolve()
	r := reporter.NewRegReporter(hr.options.Reporter, nl, os.Stdout)
	return hr.walkNodes(parser, r)
}

// Register generates report
func (hr Hranoprovod) Print(pc parser.Config) error {
	parser := parser.NewParser(pc)
	r := reporter.NewPrintReporter(hr.options.Reporter, os.Stdout)
	return hr.walkNodes(parser, r)
}

// Balance generates balance report
func (hr Hranoprovod) Balance(pc parser.Config) error {
	parser := parser.NewParser(pc)
	nl, err := loadDatabase(parser, hr.options.Global.DbFileName)
	if err != nil {
		return err
	}
	resolver.NewResolver(nl, hr.options.Resolver.ResolverMaxDepth).Resolve()
	r := reporter.NewBalanceReporter(hr.options.Reporter, nl, os.Stdout)
	return hr.walkNodes(parser, r)
}

// Lint lints file
func (hr Hranoprovod) Lint(fileName string, pc parser.Config) error {
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
func (hr Hranoprovod) ReportElement(elementName string, ascending bool, pc parser.Config) error {
	p := parser.NewParser(pc)
	nl, err := loadDatabase(p, hr.options.Global.DbFileName)
	if err != nil {
		return err
	}
	resolver.NewResolver(nl, hr.options.Resolver.ResolverMaxDepth).Resolve()
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
	for _, el := range list {
		fmt.Fprintf(hr.output, "%0.2f\t%s\n", el.Value, el.Name)
	}
	return nil
}

// ReportQuantity Generates a quantity report
func (hr Hranoprovod) ReportQuantity(ascending bool, pc parser.Config) error {
	parser := parser.NewParser(pc)
	r := reporter.NewQuantityReporter(ascending, os.Stdout)
	return hr.walkNodes(parser, r)
}

// ReportUnresolved generates report for unresolved elements
func (hr Hranoprovod) ReportUnresolved(pc parser.Config) error {
	parser := parser.NewParser(pc)
	nl, err := loadDatabase(parser, hr.options.Global.DbFileName)
	if err != nil {
		return err
	}
	resolver.NewResolver(nl, hr.options.Resolver.ResolverMaxDepth).Resolve()
	r := reporter.NewUnsolvedReporter(hr.options.Reporter, nl, os.Stdout)

	return hr.walkNodes(parser, r)
}

// CSV generates CSV export
func (hr Hranoprovod) CSV(pc parser.Config) error {
	parser := parser.NewParser(pc)
	r := reporter.NewCSVReporter(hr.options.Reporter, os.Stdout)
	return hr.walkNodes(parser, r)
}

// Stats generates statistics report
func (hr Hranoprovod) Stats() error {
	var err error
	fLog, err := os.Open(hr.options.Global.LogFileName)
	if err != nil {
		return parser.NewErrorIO(err, hr.options.Global.LogFileName)
	}
	defer fLog.Close()

	fDb, err := os.Open(hr.options.Global.DbFileName)
	if err != nil {
		return parser.NewErrorIO(err, hr.options.Global.LogFileName)
	}
	defer fDb.Close()

	countLog := 0
	var firstLogDate time.Time
	var lastLogDate time.Time
	parser.ParseStreamCallback(fLog, '#', func(n *shared.ParserNode, _ error) (stop bool) {
		lastLogDate, err = shared.ParseTime(n.Header, hr.options.Global.DateFormat)
		if err == nil {
			if firstLogDate.IsZero() {
				firstLogDate = lastLogDate
			}
		}

		countLog++
		return false
	})

	countDb := 0
	parser.ParseStreamCallback(fDb, '#', func(n *shared.ParserNode, err error) (stop bool) {
		countDb++
		return false
	})

	fmt.Fprintf(hr.output, "  Database file:      %s\n", hr.options.Global.DbFileName)
	fmt.Fprintf(hr.output, "  Database records:   %d\n", countDb)
	fmt.Fprintln(hr.output, "")
	fmt.Fprintf(hr.output, "  Log file:           %s\n", hr.options.Global.LogFileName)
	fmt.Fprintf(hr.output, "  Log records:        %d\n", countLog)
	fmt.Fprintf(hr.output, "  First record:       %s (%d days ago)\n", firstLogDate.Format(hr.options.Reporter.DateFormat), int(time.Since(firstLogDate).Hours()/24))
	fmt.Fprintf(hr.output, "  Last record:        %s (%d days ago)\n", lastLogDate.Format(hr.options.Reporter.DateFormat), int(time.Since(lastLogDate).Hours()/24))
	return nil
}

// Summary generates summary
func (hr Hranoprovod) Summary(pc parser.Config) error {
	parser := parser.NewParser(pc)
	nl, err := loadDatabase(parser, hr.options.Global.DbFileName)
	if err != nil {
		return err
	}
	resolver.NewResolver(nl, hr.options.Resolver.ResolverMaxDepth).Resolve()
	r := reporter.NewSummaryReporterTemplate(hr.options.Reporter, nl, os.Stdout)
	return hr.walkNodes(parser, r)
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

func (hr Hranoprovod) walkNodes(p parser.Parser, r reporter.Reporter) error {
	var node *shared.ParserNode
	var ln *shared.LogNode
	var err error
	var t time.Time

	go p.ParseFile(hr.options.Global.LogFileName)
	for {
		select {
		case node = <-p.Nodes:
			t, err = shared.ParseTime(node.Header, hr.options.Global.DateFormat)
			if err != nil {
				return err
			}
			if hr.inInterval(t) {
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

func (hr Hranoprovod) inInterval(t time.Time) bool {
	if hr.options.Reporter.HasBeginning && !isGoodDate(t, hr.options.Reporter.BeginningTime, dateBeginning) {
		return false
	}
	if hr.options.Reporter.HasEnd && !isGoodDate(t, hr.options.Reporter.EndTime, dateEnd) {
		return false
	}
	return true
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
