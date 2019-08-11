package reporter

import (
	"fmt"
	"io"

	"github.com/aquilax/hranoprovod-cli/shared"
)

type unsolvedReporter struct {
	options *Options
	db      shared.DBNodeList
	output  io.Writer
	list    map[string]bool
}

func newUnsolvedReporter(options *Options, db shared.DBNodeList, writer io.Writer) *unsolvedReporter {
	return &unsolvedReporter{
		options,
		db,
		writer,
		make(map[string]bool),
	}
}

func (r *unsolvedReporter) Process(ln *shared.LogNode) error {
	for _, e := range *ln.Elements {
		_, found := r.db[e.Name]
		if !found {
			r.list[e.Name] = true
		}
	}
	return nil
}

func (r *unsolvedReporter) Flush() error {
	for name := range r.list {
		fmt.Fprintln(r.output, name)
	}
	return nil
}
