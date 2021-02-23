package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/aquilax/hranoprovod-cli/parser"
	"github.com/aquilax/hranoprovod-cli/reporter"
	"github.com/aquilax/hranoprovod-cli/resolver"
	"github.com/aquilax/hranoprovod-cli/shared"
)

// Hranoprovod is the main app type
type Hranoprovod struct {
	options *Options
}

// NewHranoprovod creates new application
func NewHranoprovod(options *Options) Hranoprovod {
	return Hranoprovod{options}
}

// Register generates report
func (hr Hranoprovod) Register() error {
	parser := parser.NewParser(&hr.options.Parser)
	nl, err := loadDatabase(parser, hr.options.Global.DbFileName)
	if err != nil {
		return err
	}
	resolver.NewResolver(nl, hr.options.Resolver.ResolverMaxDepth).Resolve()
	r := reporter.NewRegReporter(&hr.options.Reporter, nl, os.Stdout)
	return hr.walkNodes(parser, r)
}

// Balance generates balance report
func (hr Hranoprovod) Balance() error {
	parser := parser.NewParser(&hr.options.Parser)
	nl, err := loadDatabase(parser, hr.options.Global.DbFileName)
	if err != nil {
		return err
	}
	resolver.NewResolver(nl, hr.options.Resolver.ResolverMaxDepth).Resolve()
	r := reporter.NewBalanceReporter(&hr.options.Reporter, nl, os.Stdout)
	return hr.walkNodes(parser, r)
}

// Lint lints file
func (hr Hranoprovod) Lint(fileName string) error {
	p := parser.NewParser(&hr.options.Parser)
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
func (hr Hranoprovod) ReportElement(elementName string, ascending bool) error {
	p := parser.NewParser(&hr.options.Parser)
	nl, err := loadDatabase(p, hr.options.Global.DbFileName)
	if err != nil {
		return err
	}
	resolver.NewResolver(nl, hr.options.Resolver.ResolverMaxDepth).Resolve()
	var list []shared.Element
	for name, node := range nl {
		for _, el := range node.Elements {
			if el.Name == elementName {
				list = append(list, shared.NewElement(name, el.Val))
			}
		}
	}
	if ascending {
		sort.SliceStable(list, func(i, j int) bool {
			return list[i].Val > list[j].Val
		})
	} else {
		sort.SliceStable(list, func(i, j int) bool {
			return list[i].Val < list[j].Val
		})
	}
	for _, el := range list {
		fmt.Printf("%0.2f\t%s\n", el.Val, el.Name)
	}
	return nil
}

// ReportUnresolved generates report for unresolved elements
func (hr Hranoprovod) ReportUnresolved() error {
	parser := parser.NewParser(&hr.options.Parser)
	nl, err := loadDatabase(parser, hr.options.Global.DbFileName)
	if err != nil {
		return err
	}
	resolver.NewResolver(nl, hr.options.Resolver.ResolverMaxDepth).Resolve()
	r := reporter.NewUnsolvedReporter(&hr.options.Reporter, nl, os.Stdout)

	return hr.walkNodes(parser, r)
}

// CSV generates CSV export
func (hr Hranoprovod) CSV() error {
	parser := parser.NewParser(&hr.options.Parser)
	r := reporter.NewCSVReporter(&hr.options.Reporter, os.Stdout)
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

	fmt.Printf("  Database file:      %s\n", hr.options.Global.DbFileName)
	fmt.Printf("  Database records:   %d\n", countDb)
	fmt.Println("")
	fmt.Printf("  Log file:           %s\n", hr.options.Global.LogFileName)
	fmt.Printf("  Log records:        %d\n", countLog)
	fmt.Printf("  First record:       %s (%d days ago)\n", firstLogDate.Format(hr.options.Reporter.DateFormat), int(time.Now().Sub(firstLogDate).Hours()/24))
	fmt.Printf("  Last record:        %s (%d days ago)\n", lastLogDate.Format(hr.options.Reporter.DateFormat), int(time.Now().Sub(lastLogDate).Hours()/24))
	return nil
}

// Summary generates summary
func (hr Hranoprovod) Summary() error {
	parser := parser.NewParser(&hr.options.Parser)
	nl, err := loadDatabase(parser, hr.options.Global.DbFileName)
	if err != nil {
		return err
	}
	resolver.NewResolver(nl, hr.options.Resolver.ResolverMaxDepth).Resolve()
	r := reporter.NewSummaryReporterTemplate(&hr.options.Reporter, nl, os.Stdout)
	return hr.walkNodes(parser, r)
}

func loadDatabase(p parser.Parser, fileName string) (shared.DBNodeList, error) {
	nodeList := shared.NewDBNodeList()
	// Database must be optional. If the default file name is used and the file is not found,
	// return empty node list
	if fileName == defaultDbFilename {
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
				ln, err = shared.NewLogNodeFromElements(t, node.Elements)
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
