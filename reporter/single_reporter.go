package reporter

import (
	"fmt"
	"io"
	"time"

	"github.com/aquilax/hranoprovod-cli/accumulator"
	"github.com/aquilax/hranoprovod-cli/shared"
)

// Reporter is the main report structure
type SingleReporter struct {
	options *Options
	db      *shared.NodeList
	output  io.Writer
}

// NewReporter creates new reporter
func NewSingleReporter(options *Options, db *shared.NodeList, writer io.Writer) *SingleReporter {
	return &SingleReporter{
		options,
		db,
		writer,
	}
}

func (r *SingleReporter) Process(ln *shared.LogNode) error {
	acc := accumulator.NewAccumulator()
	singleElement := r.options.SingleElement
	for _, e := range *ln.Elements {
		repl, found := (*r.db)[e.Name]
		if found {
			for _, repl := range *repl.Elements {
				if repl.Name == singleElement {
					acc.Add(repl.Name, repl.Val*e.Val)
				}
			}
		} else {
			if e.Name == singleElement {
				acc.Add(e.Name, e.Val)
			}
		}
	}
	if len(*acc) > 0 {
		arr := (*acc)[singleElement]
		r.PrintSingleElementRow(ln.Time, r.options.SingleElement, arr[accumulator.Positive], arr[accumulator.Negative])
	}
	return nil

}

func (r *SingleReporter) Flush() error {
	return nil
}

func (r *SingleReporter) PrintSingleElementRow(ts time.Time, name string, pos float32, neg float32) {
	format := "%s %20s %10.2f %10.2f =%10.2f\n"
	if r.options.CSV {
		format = "%s;\"%s\";%0.2f;%0.2f;%0.2f\n"
	}
	fmt.Fprintf(r.output, format, ts.Format(r.options.DateFormat), name, pos, -1*neg, pos+neg)
}
