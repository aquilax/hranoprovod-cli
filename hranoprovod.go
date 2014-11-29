package main

import (
	"github.com/Hranoprovod/api-client"
	"github.com/Hranoprovod/parser"
	"github.com/Hranoprovod/processor"
	"github.com/Hranoprovod/reporter"
	"github.com/Hranoprovod/resolver"
	"github.com/Hranoprovod/shared"
	"os"
	"time"
)

// Hranoprovod is the main app type
type Hranoprovod struct {
	config *Config
}

// RegisterOptions contains the options for the register command
type RegisterOptions struct {
	DbFileName       string
	LogFileName      string
	ResolverMaxDepth int
	DateFormat       string
	CSV              bool

	Beginning     string
	HasBeginning  bool
	BegginingTime time.Time

	End     string
	HasEnd  bool
	EndTime time.Time

	Unresolved    bool
	SingleElement string
	SingleFood    string
	Totals        bool
	Color         bool
}

// NewHranoprovod creates new application
func NewHranoprovod(config *Config) *Hranoprovod {
	return &Hranoprovod{config}
}

// Validate validates the register options
func (ro *RegisterOptions) Validate() error {
	var err error
	if ro.Beginning != "" {
		ro.BegginingTime, err = time.Parse(ro.DateFormat, ro.Beginning)
		if err != nil {
			return err
		}
		ro.HasBeginning = true
	}
	if ro.End != "" {
		ro.EndTime, err = time.Parse(ro.DateFormat, ro.End)
		if err != nil {
			return err
		}
		ro.HasEnd = true
	}
	return nil
}

// Register generates report
func (hr *Hranoprovod) Register() error {
	parser := parser.NewParser(parser.NewDefaultOptions())
	nl, err := hr.loadDatabase(parser, hr.config.Global.DbFileName)
	if err != nil {
		return err
	}
	resolver.NewResolver(nl, hr.config.Resolver.ResolverMaxDepth).Resolve()
	return hr.processLog(parser, nl)
}

// Search searches the API for the provided query
func (hr *Hranoprovod) Search(q string) error {
	api := client.NewAPIClient(client.GetDefaultAPIClientOptions())
	nl, err := api.Search(q)
	if err != nil {
		return err
	}
	rp := reporter.NewReporter(reporter.NewDefaultOptions(), os.Stdout)
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
		&hr.config.Processor,
		nl,
		reporter.NewReporter(&hr.config.Reporter, os.Stdout),
	)

	go p.ParseFile(hr.config.Global.LogFileName)
	for {
		select {
		case node := <-p.Nodes:
			ln, err := shared.NewLogNodeFromNode(node, hr.config.Global.DateFormat)
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
