package app

import (
	"fmt"
	"sort"
	"time"

	"github.com/aquilax/hranoprovod-cli/v2/parser"
	"github.com/aquilax/hranoprovod-cli/v2/reporter"
	"github.com/aquilax/hranoprovod-cli/v2/resolver"
	"github.com/aquilax/hranoprovod-cli/v2/shared"
)

// Register generates report
func Register(gc GlobalConfig, pc parser.Config, rc resolver.Config, rpc reporter.Config, fc FilterConfig) error {
	parser := parser.NewParser(pc)
	nl, err := loadDatabase(parser, gc.DbFileName)
	if err != nil {
		return err
	}
	err = resolver.NewResolver(nl, rc).Resolve()
	if err != nil {
		return err
	}
	r := reporter.NewRegReporter(rpc, nl, rpc.Output)
	return walkNodes(gc.LogFileName, gc.DateFormat, parser, getIntervalNodeFilter(fc), r)
}

// Balance generates balance report
func Balance(gc GlobalConfig, pc parser.Config, rc resolver.Config, rpc reporter.Config, fc FilterConfig) error {
	parser := parser.NewParser(pc)
	nl, err := loadDatabase(parser, gc.DbFileName)
	if err != nil {
		return err
	}
	err = resolver.NewResolver(nl, rc).Resolve()
	if err != nil {
		return err
	}
	r := reporter.NewBalanceReporter(rpc, nl, rpc.Output)
	return walkNodes(gc.LogFileName, gc.DateFormat, parser, getIntervalNodeFilter(fc), r)
}

// ReportUnresolved generates report for unresolved elements
func ReportUnresolved(gc GlobalConfig, pc parser.Config, rc resolver.Config, rpc reporter.Config, fc FilterConfig) error {
	parser := parser.NewParser(pc)
	nl, err := loadDatabase(parser, gc.DbFileName)
	if err != nil {
		return err
	}
	err = resolver.NewResolver(nl, rc).Resolve()
	if err != nil {
		return err
	}
	r := reporter.NewUnsolvedReporter(rpc, nl, rpc.Output)
	return walkNodes(gc.LogFileName, gc.DateFormat, parser, getIntervalNodeFilter(fc), r)
}

// Summary generates summary
func Summary(gc GlobalConfig, pc parser.Config, rc resolver.Config, rpc reporter.Config, fc FilterConfig) error {
	parser := parser.NewParser(pc)
	nl, err := loadDatabase(parser, gc.DbFileName)
	if err != nil {
		return err
	}
	err = resolver.NewResolver(nl, rc).Resolve()
	if err != nil {
		return err
	}
	r := reporter.NewSummaryReporterTemplate(rpc, nl, rpc.Output)
	return walkNodes(gc.LogFileName, gc.DateFormat, parser, getIntervalNodeFilter(fc), r)
}

// Print reads and prints back out the log file
func Print(gc GlobalConfig, pc parser.Config, rpc reporter.Config, fc FilterConfig) error {
	parser := parser.NewParser(pc)
	r := reporter.NewPrintReporter(rpc, rpc.Output)
	return walkNodes(gc.LogFileName, gc.DateFormat, parser, getIntervalNodeFilter(fc), r)
}

// ReportQuantity Generates a quantity report
func ReportQuantity(gc GlobalConfig, ascending bool, pc parser.Config, rpc reporter.Config, fc FilterConfig) error {
	parser := parser.NewParser(pc)
	r := reporter.NewQuantityReporter(ascending, rpc.Output)
	return walkNodes(gc.LogFileName, gc.DateFormat, parser, getIntervalNodeFilter(fc), r)
}

// CSV generates CSV export
func CSV(gc GlobalConfig, pc parser.Config, rpc reporter.Config, fc FilterConfig) error {
	parser := parser.NewParser(pc)
	r := reporter.NewCSVReporter(rpc, rpc.Output)
	return walkNodes(gc.LogFileName, gc.DateFormat, parser, getIntervalNodeFilter(fc), r)
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
				fmt.Println(err)
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
	err = resolver.NewResolver(nl, rc).Resolve()
	if err != nil {
		return err
	}
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

// Stats generates statistics report
func Stats(gc GlobalConfig, pc parser.Config, rpc reporter.Config) error {
	var err error
	var firstLogDate time.Time
	var lastLogDate time.Time

	countLog := 0
	err = parser.ParseFileCallback(gc.LogFileName, pc, func(n *shared.ParserNode, _ error) (stop bool) {
		lastLogDate, err = time.Parse(gc.DateFormat, n.Header)
		if err == nil {
			if firstLogDate.IsZero() {
				firstLogDate = lastLogDate
			}
		}
		countLog++
		return false
	})
	if err != nil {
		return err
	}

	countDb := 0
	err = parser.ParseFileCallback(gc.DbFileName, pc, func(n *shared.ParserNode, _ error) (stop bool) {
		countDb++
		return false
	})
	if err != nil {
		return err
	}

	return reporter.NewStatsReporter(rpc, &reporter.Stats{
		DbFileName:      gc.DbFileName,
		LogFileName:     gc.LogFileName,
		DbRecordsCount:  countDb,
		LogRecordsCount: countLog,
		LogFirstRecord:  firstLogDate,
		LogLastRecord:   lastLogDate,
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

type LogNodeFilter = func(t time.Time, node *shared.ParserNode) (bool, error)

func getIntervalNodeFilter(fc FilterConfig) *LogNodeFilter {
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

func walkNodes(logFileName string, dateFormat string, p parser.Parser, filter *LogNodeFilter, r reporter.Reporter) error {
	var node *shared.ParserNode
	var ln *shared.LogNode
	var err error
	var t time.Time
	var ok bool

	go p.ParseFile(logFileName)
	for {
		select {
		case node = <-p.Nodes:
			t, err = time.Parse(dateFormat, node.Header)
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
