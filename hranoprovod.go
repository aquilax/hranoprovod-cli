package main

import (
	"github.com/Hranoprovod/api-client"
	"github.com/Hranoprovod/parser"
	"github.com/Hranoprovod/reporter"
	"github.com/Hranoprovod/resolver"
	"github.com/Hranoprovod/shared"
	"os"
)

// Hranoprovod is the main app type
type Hranoprovod struct{}

// NewHranoprovod creates new application
func NewHranoprovod() *Hranoprovod {
	return &Hranoprovod{}
}

// Register generates report
func (hr *Hranoprovod) Register(dbFileName string) error {
	parser := parser.NewParser(parser.NewDefaultOptions())
	nodeList, err := hr.loadDatabase(parser, dbFileName)
	if err != nil {
		return err
	}

	// TODO: Magic number
	resolverMaxDepth := 10
	resolver.NewResolver(nodeList, resolverMaxDepth).Resolve()

	return hr.processLog(parser, nodeList)
}

// Search searches the API for the provided query
func (hr *Hranoprovod) Search(q string) error {
	api := client.NewAPIClient(client.GetDefaultAPIClientOptions())
	nl, err := api.Search(q)
	if err != nil {
		return err
	}
	rp := reporter.NewReporter(os.Stdout)
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

	// processor := NewProcessor(
	// 	options,
	// 	nodeList,
	// 	NewReporter(options, os.Stdout),
	// )

	// go parser.parseFile(options.logFileName)
	// for {
	// 	select {
	// 	case node := <-parser.nodes:
	// 		processor.process(node)
	// 	case breakingError := <-parser.errors:
	// 		return breakingError
	// 	case <-parser.done:
	// 		return nil
	// 	}
	// }
	return nil
}
