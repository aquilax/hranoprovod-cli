package reporter

import (
	"fmt"
	"io"
	"time"

	"github.com/aquilax/hranoprovod-cli/v2/accumulator"
	"github.com/aquilax/hranoprovod-cli/v2/shared"
)

// singleReporter outputs report for single food
type singleReporter struct {
	options Options
	db      shared.DBNodeList
	output  io.Writer
}

func newSingleReporter(options Options, db shared.DBNodeList, writer io.Writer) *singleReporter {
	return &singleReporter{
		options,
		db,
		writer,
	}
}

func (r *singleReporter) Process(ln *shared.LogNode) error {
	acc := accumulator.NewAccumulator()
	singleElement := r.options.SingleElement
	for _, e := range ln.Elements {
		repl, found := r.db[e.Name]
		if found {
			for _, repl := range repl.Elements {
				if repl.Name == singleElement {
					acc.Add(repl.Name, repl.Value*e.Value)
				}
			}
		} else {
			if e.Name == singleElement {
				acc.Add(e.Name, e.Value)
			}
		}
	}
	if len(acc) > 0 {
		arr := (acc)[singleElement]
		r.printSingleElementRow(ln.Time, r.options.SingleElement, arr[accumulator.Positive], arr[accumulator.Negative])
	}
	return nil

}

func (r *singleReporter) Flush() error {
	return nil
}

func (r *singleReporter) printSingleElementRow(ts time.Time, name string, pos float64, neg float64) {
	format := "%s %20s %10.2f %10.2f =%10.2f\n"
	if r.options.CSV {
		format = "%s;\"%s\";%0.2f;%0.2f;%0.2f\n"
	}
	fmt.Fprintf(r.output, format, ts.Format(r.options.DateFormat), name, pos, -1*neg, pos+neg)
}
