package main

import (
	"os"
)

type Hranoprovod struct {}

func NewHranoprovod() *Hranoprovod {
	return &Hranoprovod{}
}

func (hr *Hranoprovod) register() error {
	return nil
}

func (hr *Hranoprovod) search(q string) error {
	api := NewAPIClient(GetDefaultAPIClientOptions())
	nl, err := api.Search(q)
	if err != nil {
		return err
	}
	reporter := NewReporter(os.Stdout)
	return reporter.PrintAPISearchResult(*nl)
}