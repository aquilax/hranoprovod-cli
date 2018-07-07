package reporter

import (
	"fmt"
	"io"

	"github.com/aquilax/hranoprovod-cli/shared"
)

// Reporter is the main report structure
type UnsolvedReporter struct {
	options *Options
	db      *shared.NodeList
	output  io.Writer
	list    map[string]bool
}

// NewReporter creates new reporter
func NewUnsolvedReporter(options *Options, db *shared.NodeList, writer io.Writer) *UnsolvedReporter {
	return &UnsolvedReporter{
		options,
		db,
		writer,
		make(map[string]bool),
	}
}

func (r *UnsolvedReporter) Process(ln *shared.LogNode) error {
	if (r.options.HasBeginning && !isGoodDate(ln.Time, r.options.BeginningTime, dateBeginning)) || (r.options.HasEnd && !isGoodDate(ln.Time, r.options.EndTime, dateEnd)) {
		return nil
	}
	for _, e := range *ln.Elements {
		_, found := (*r.db)[e.Name]
		if !found {
			r.list[e.Name] = true
		}
	}
	return nil
}

func (r *UnsolvedReporter) Flush() error {
	for name, _ := range r.list {
		r.PrintUnresolvedRow(name)
	}
	return nil
}

func (r *UnsolvedReporter) PrintUnresolvedRow(name string) {
	fmt.Fprintln(r.output, name)
}
