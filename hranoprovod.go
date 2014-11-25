package main

import (
	"github.com/Hranoprovod/api-client"
	"github.com/Hranoprovod/parser"
	"github.com/Hranoprovod/reporter"
	"os"
)

// Hranoprovod is the main app type
type Hranoprovod struct{}

// NewHranoprovod creates new application
func NewHranoprovod() *Hranoprovod {
	return &Hranoprovod{}
}

// Register generates report
func (hr *Hranoprovod) Register() error {
	return nil
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
