package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	client "github.com/aquilax/hranoprovod-cli/api-client"
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
	nl, err := hr.loadDatabase(parser, hr.options.Global.DbFileName)
	if err != nil {
		return err
	}
	resolver.NewResolver(nl, hr.options.Resolver.ResolverMaxDepth).Resolve()
	return hr.processLog(parser, nl)
}

// Balance generates balance report
func (hr Hranoprovod) Balance() error {
	parser := parser.NewParser(&hr.options.Parser)
	nl, err := hr.loadDatabase(parser, hr.options.Global.DbFileName)
	if err != nil {
		return err
	}
	resolver.NewResolver(nl, hr.options.Resolver.ResolverMaxDepth).Resolve()
	return hr.processBalance(parser, nl)
}

// Search searches the API for the provided query
func (hr Hranoprovod) Search(q string) error {
	api := client.NewAPIClient(&hr.options.API)
	nl, err := api.Search(q)
	if err != nil {
		return err
	}
	rp := reporter.NewAPIReporter(&hr.options.Reporter, os.Stdout)
	return rp.PrintAPISearchResult(*nl)
}

// Add adds new item to the log
func (hr Hranoprovod) Add(name string, qty string) error {
	println("Adding " + name + " : " + qty)
	return nil
}

// Lint lints file
func (hr Hranoprovod) Lint(fileName string) error {
	p := parser.NewParser(&hr.options.Parser)
	go p.ParseFile(fileName)
	return func() error {
		for {
			select {
			case _ = <-p.Nodes:
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
	nl, err := hr.loadDatabase(p, hr.options.Global.DbFileName)
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

func (hr Hranoprovod) Stats() error {
	f, err := os.Open(hr.options.Global.LogFileName)
	if err != nil {
		return parser.NewErrorIO(err, hr.options.Global.LogFileName)
	}
	defer f.Close()

	count := 0
	parser.ParseStreamCallback(f, '#', func(n *shared.ParserNode, err error) (stop bool) {
		count++
		return false
	})

	fmt.Printf("Food database: %s\n", hr.options.Global.DbFileName)
	fmt.Printf("Log database: %s\n", hr.options.Global.LogFileName)
	fmt.Println("")
	fmt.Printf("Log records: %d\n", count)
	return nil
}

// ReportUnresolved generates report for unresolved elements
func (hr Hranoprovod) ReportUnresolved() error {
	var err error
	p := parser.NewParser(&hr.options.Parser)
	nl, err := hr.loadDatabase(p, hr.options.Global.DbFileName)
	if err != nil {
		return err
	}
	resolver.NewResolver(nl, hr.options.Resolver.ResolverMaxDepth).Resolve()
	r := reporter.NewUnsolvedReporter(&hr.options.Reporter, nl, os.Stdout)
	var node *shared.ParserNode
	var ln *shared.LogNode

	go p.ParseFile(hr.options.Global.LogFileName)
	for {
		select {
		case node = <-p.Nodes:
			ln, err = shared.NewLogNodeFromNode(node, hr.options.Global.DateFormat)
			if err != nil {
				return err
			}
			if hr.inInterval(ln.Time) {
				r.Process(ln)
			}
		case breakingError := <-p.Errors:
			return breakingError
		case <-p.Done:
			r.Flush()
			return nil
		}
	}
}

func (hr Hranoprovod) loadDatabase(p parser.Parser, fileName string) (shared.DBNodeList, error) {
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

func (hr Hranoprovod) processLog(p parser.Parser, nl shared.DBNodeList) error {
	r := reporter.NewRegReporter(&hr.options.Reporter, nl, os.Stdout)
	var node *shared.ParserNode
	var ln *shared.LogNode
	var err error

	go p.ParseFile(hr.options.Global.LogFileName)
	for {
		select {
		case node = <-p.Nodes:
			ln, err = shared.NewLogNodeFromNode(node, hr.options.Global.DateFormat)
			if err != nil {
				return err
			}
			if hr.inInterval(ln.Time) {
				r.Process(ln)
			}
		case breakingError := <-p.Errors:
			return breakingError
		case <-p.Done:
			r.Flush()
			return nil
		}
	}
}

func (hr Hranoprovod) processBalance(p parser.Parser, nl shared.DBNodeList) error {
	r := reporter.NewBalanceReporter(&hr.options.Reporter, nl, os.Stdout)
	var node *shared.ParserNode
	var ln *shared.LogNode
	var err error

	go p.ParseFile(hr.options.Global.LogFileName)
	for {
		select {
		case node = <-p.Nodes:
			ln, err = shared.NewLogNodeFromNode(node, hr.options.Global.DateFormat)
			if err != nil {
				return err
			}
			if hr.inInterval(ln.Time) {
				r.Process(ln)
			}
		case breakingError := <-p.Errors:
			return breakingError
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

// CSV generates CSV export
func (hr Hranoprovod) CSV() error {
	p := parser.NewParser(&hr.options.Parser)
	r := reporter.NewCSVReporter(&hr.options.Reporter, os.Stdout)

	var node *shared.ParserNode
	var ln *shared.LogNode
	var err error

	go p.ParseFile(hr.options.Global.LogFileName)
	for {
		select {
		case node = <-p.Nodes:
			ln, err = shared.NewLogNodeFromNode(node, hr.options.Global.DateFormat)
			if err != nil {
				return err
			}
			if hr.inInterval(ln.Time) {
				if err = r.Process(ln); err != nil {
					return err
				}
			}
		case breakingError := <-p.Errors:
			return breakingError
		case <-p.Done:
			r.Flush()
			return nil
		}
	}
}

// CompareType identifies the type of date comparison
type CompareType bool

const (
	dateBeginning CompareType = true
	dateEnd       CompareType = false
)

func isGoodDate(time, compareTime time.Time, compareType CompareType) bool {
	if time.Equal(compareTime) {
		return true
	}
	if compareType == dateBeginning {
		return time.After(compareTime)
	}
	return time.Before(compareTime)
}
