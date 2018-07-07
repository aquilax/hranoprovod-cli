package main

import (
	"os"

	"github.com/aquilax/hranoprovod-cli/api-client"
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
func NewHranoprovod(options *Options) *Hranoprovod {
	return &Hranoprovod{options}
}

// Register generates report
func (hr *Hranoprovod) Register() error {
	parser := parser.NewParser(&hr.options.Parser)
	nl, err := hr.loadDatabase(parser, hr.options.Global.DbFileName)
	if err != nil {
		return err
	}
	resolver.NewResolver(nl, hr.options.Resolver.ResolverMaxDepth).Resolve()
	return hr.processLog(parser, nl)
}

// Balance generates balance report
func (hr *Hranoprovod) Balance() error {
	parser := parser.NewParser(&hr.options.Parser)
	nl, err := hr.loadDatabase(parser, hr.options.Global.DbFileName)
	if err != nil {
		return err
	}
	resolver.NewResolver(nl, hr.options.Resolver.ResolverMaxDepth).Resolve()
	return hr.processBalance(parser, nl)
}

// Search searches the API for the provided query
func (hr *Hranoprovod) Search(q string) error {
	api := client.NewAPIClient(&hr.options.API)
	nl, err := api.Search(q)
	if err != nil {
		return err
	}
	rp := reporter.NewAPIReporter(&hr.options.Reporter, os.Stdout)
	return rp.PrintAPISearchResult(*nl)
}

// Add adds new item to the log
func (hr *Hranoprovod) Add(name string, qty string) error {
	println("Adding " + name + " : " + qty)
	return nil
}

// Lint lints file
func (hr *Hranoprovod) Lint(fileName string) error {
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

func (hr *Hranoprovod) loadDatabase(p *parser.Parser, fileName string) (*shared.NodeList, error) {
	nodeList := shared.NewNodeList()
	go p.ParseFile(fileName)
	return func() (*shared.NodeList, error) {
		for {
			select {
			case node := <-p.Nodes:
				nodeList.Push(node)
			case error := <-p.Errors:
				return nil, error
			case <-p.Done:
				return nodeList, nil
			}
		}
	}()
}

func (hr *Hranoprovod) processLog(p *parser.Parser, nl *shared.NodeList) error {
	r := reporter.NewReporter(reporter.Reg, &hr.options.Reporter, nl, os.Stdout)

	go p.ParseFile(hr.options.Global.LogFileName)
	for {
		select {
		case node := <-p.Nodes:
			ln, err := shared.NewLogNodeFromNode(node, hr.options.Global.DateFormat)
			if err != nil {
				return err
			}
			r.Process(ln)
		case breakingError := <-p.Errors:
			return breakingError
		case <-p.Done:
			return nil
		}
	}
}

func (hr *Hranoprovod) processBalance(p *parser.Parser, nl *shared.NodeList) error {
	_ = reporter.NewReporter(reporter.Reg, &hr.options.Reporter, nl, os.Stdout)
	return nil
}
