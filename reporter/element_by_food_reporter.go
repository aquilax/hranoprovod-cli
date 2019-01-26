package reporter

import (
	"fmt"
	"io"

	"github.com/aquilax/hranoprovod-cli/accumulator"
	"github.com/aquilax/hranoprovod-cli/shared"
)

// elementByFoodReporter outputs report for single element groupped by food
type elementByFoodReporter struct {
	options *Options
	db      *shared.NodeList
	output  io.Writer
	acc     *accumulator.Accumulator
}

func newElementByFoodReporter(options *Options, db *shared.NodeList, writer io.Writer) *elementByFoodReporter {
	return &elementByFoodReporter{
		options,
		db,
		writer,
		accumulator.NewAccumulator(),
	}
}

func (r *elementByFoodReporter) Process(ln *shared.LogNode) error {
	singleElement := r.options.SingleElement
	for _, e := range *ln.Elements {
		node, found := (*r.db)[e.Name]
		if found {
			for _, repl := range *node.Elements {
				if repl.Name == singleElement {
					r.acc.Add(node.Header, repl.Val*e.Val)
				}
			}
		}
	}
	return nil
}

func (r *elementByFoodReporter) Flush() error {
	for name, arr := range *r.acc {
		r.printSingleElementByFoodRow(name, arr[accumulator.Positive], arr[accumulator.Negative])
	}
	return nil
}

func (r *elementByFoodReporter) printSingleElementByFoodRow(name string, pos float32, neg float32) {
	format := "%10.2f\t%s\n"
	fmt.Fprintf(r.output, format, pos+neg, name)
}
