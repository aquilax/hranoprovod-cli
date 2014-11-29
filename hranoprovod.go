package main

import (
	"github.com/Hranoprovod/api-client"
	"github.com/Hranoprovod/parser"
	"github.com/Hranoprovod/processor"
	"github.com/Hranoprovod/reporter"
	"github.com/Hranoprovod/resolver"
	"github.com/Hranoprovod/shared"
	"os"
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
	parser := parser.NewParser(parser.NewDefaultOptions())
	nl, err := hr.loadDatabase(parser, hr.options.Global.DbFileName)
	if err != nil {
		return err
	}
	resolver.NewResolver(nl, hr.options.Resolver.ResolverMaxDepth).Resolve()
	return hr.processLog(parser, nl)
}

// Search searches the API for the provided query
func (hr *Hranoprovod) Search(q string) error {
	api := client.NewAPIClient(&hr.options.API)
	nl, err := api.Search(q)
	if err != nil {
		return err
	}
	rp := reporter.NewReporter(&hr.options.Reporter, os.Stdout)
	return rp.PrintAPISearchResult(*nl)
}

// Add adds new item to the log
func (hr *Hranoprovod) Add(name string, qty string) error {
	println("Adding " + name + " : " + qty)
	return nil
}

// Lint lints file
func (hr *Hranoprovod) Lint(fileName string) error {
	p := parser.NewParser(parser.NewDefaultOptions())
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
	pr := processor.NewProcessor(
		&hr.options.Processor,
		nl,
		reporter.NewReporter(&hr.options.Reporter, os.Stdout),
	)

	go p.ParseFile(hr.options.Global.LogFileName)
	for {
		select {
		case node := <-p.Nodes:
			ln, err := shared.NewLogNodeFromNode(node, hr.options.Global.DateFormat)
			if err != nil {
				return err
			}
			pr.Process(ln)
		case breakingError := <-p.Errors:
			return breakingError
		case <-p.Done:
			return nil
		}
	}
}
