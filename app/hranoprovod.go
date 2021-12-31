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

type resolvedCallback = func(gc GlobalConfig, p parser.Parser, nl shared.DBNodeList, rpc reporter.Config, nf *LogNodeFilter) error

func withResolvedDatabase(gc GlobalConfig, p parser.Parser, rc resolver.Config, rpc reporter.Config, fc FilterConfig, cb resolvedCallback) error {
	if nl, err := loadDatabase(p, gc.DbFileName); err == nil {
		if nl, err = resolver.Resolve(rc, nl); err == nil {
			return cb(gc, p, nl, rpc, getIntervalNodeFilter(fc))
		} else {
			return err
		}
	} else {
		return err
	}
}

// Register generates report
func Register(gc GlobalConfig, pc parser.Config, rc resolver.Config, rpc reporter.Config, fc FilterConfig) error {
	return withResolvedDatabase(gc, parser.NewParser(pc), rc, rpc, fc,
		func(gc GlobalConfig, p parser.Parser, nl shared.DBNodeList, rpc reporter.Config, nf *LogNodeFilter) error {
			r := reporter.NewRegReporter(rpc, nl, rpc.Output)
			return walkNodes(gc.LogFileName, gc.DateFormat, p, nf, r)
		})
}

// Balance generates balance report
func Balance(gc GlobalConfig, pc parser.Config, rc resolver.Config, rpc reporter.Config, fc FilterConfig) error {
	return withResolvedDatabase(gc, parser.NewParser(pc), rc, rpc, fc,
		func(gc GlobalConfig, p parser.Parser, nl shared.DBNodeList, rpc reporter.Config, nf *LogNodeFilter) error {
			r := reporter.NewBalanceReporter(rpc, nl, rpc.Output)
			return walkNodes(gc.LogFileName, gc.DateFormat, p, nf, r)
		})
}

// ReportUnresolved generates report for unresolved elements
func ReportUnresolved(gc GlobalConfig, pc parser.Config, rc resolver.Config, rpc reporter.Config, fc FilterConfig) error {
	return withResolvedDatabase(gc, parser.NewParser(pc), rc, rpc, fc,
		func(gc GlobalConfig, p parser.Parser, nl shared.DBNodeList, rpc reporter.Config, nf *LogNodeFilter) error {
			r := reporter.NewUnsolvedReporter(rpc, nl, rpc.Output)
			return walkNodes(gc.LogFileName, gc.DateFormat, p, nf, r)
		})
}

// Summary generates summary
func Summary(gc GlobalConfig, pc parser.Config, rc resolver.Config, rpc reporter.Config, fc FilterConfig) error {
	return withResolvedDatabase(gc, parser.NewParser(pc), rc, rpc, fc,
		func(gc GlobalConfig, p parser.Parser, nl shared.DBNodeList, rpc reporter.Config, nf *LogNodeFilter) error {
			r := reporter.NewSummaryReporterTemplate(rpc, nl, rpc.Output)
			return walkNodes(gc.LogFileName, gc.DateFormat, p, nf, r)
		})
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

// CSVLog generates CSV export of the log
func CSVLog(gc GlobalConfig, pc parser.Config, rpc reporter.Config, fc FilterConfig) error {
	parser := parser.NewParser(pc)
	r := reporter.NewCSVReporter(rpc, rpc.Output)
	return walkNodes(gc.LogFileName, gc.DateFormat, parser, getIntervalNodeFilter(fc), r)
}

// CSVDatabase generates CSV export of the database
func CSVDatabase(fileName string, pc parser.Config, rpc reporter.Config) error {
	p := parser.NewParser(pc)
	r := reporter.NewCSVDatabaseReporter(rpc)
	go p.ParseFile(fileName)
	return func() error {
		for {
			select {
			case node := <-p.Nodes:
				r.Process(shared.NewDBNodeFromNode(node))
			case error := <-p.Errors:
				return error
			case <-p.Done:
				return r.Flush()
			}
		}
	}()
}

// CSVDatabaseResolved generates CSV export of the resolved database
func CSVDatabaseResolved(fileName string, pc parser.Config, rpc reporter.Config, rc resolver.Config) error {
	p := parser.NewParser(pc)
	nl, err := mustLoadDatabase(p, fileName)
	if err != nil {
		return err
	}
	nl, err = resolver.Resolve(rc, nl)
	if err != nil {
		return err
	}
	r := reporter.NewCSVDatabaseReporter(rpc)
	for _, n := range nl {
		if err = r.Process(n); err != nil {
			return err
		}
	}
	return r.Flush()
}

// Lint lints file
func Lint(fileName string, silent bool, pc parser.Config, rpc reporter.Config) error {
	parser := parser.NewParser(pc)
	go parser.ParseFile(fileName)
	err := func() error {
		for {
			select {
			case <-parser.Nodes:
			case err := <-parser.Errors:
				fmt.Fprintln(rpc.Output, err)
			case <-parser.Done:
				return nil
			}
		}
	}()
	if err != nil {
		return err
	}
	if !silent {
		fmt.Fprintln(rpc.Output, "No errors found")
	}
	return nil
}

// ReportElement generates report for single element
func ReportElement(dbFileName string, elementName string, ascending bool, pc parser.Config, rc resolver.Config, rpc reporter.Config) error {
	parser := parser.NewParser(pc)
	nl, err := mustLoadDatabase(parser, dbFileName)
	if err != nil {
		return err
	}
	nl, err = resolver.Resolve(rc, nl)
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
	// Database must be optional. If the default file name is used and the file is not found,
	// return empty node list
	if fileName == DefaultDbFilename {
		if exists, _ := fileExists(fileName); !exists {
			return shared.NewDBNodeList(), nil
		}
	}
	return mustLoadDatabase(p, fileName)
}

func mustLoadDatabase(p parser.Parser, fileName string) (shared.DBNodeList, error) {
	nodeList := shared.NewDBNodeList()
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
