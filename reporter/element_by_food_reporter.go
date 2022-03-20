package reporter

import (
	"bufio"
	"fmt"

	"github.com/aquilax/hranoprovod-cli/v2"
	"github.com/aquilax/hranoprovod-cli/v2/accumulator"
)

// elementByFoodReporter outputs report for single element groupped by food
type elementByFoodReporter struct {
	config Config
	db     hranoprovod.DBNodeMap
	output *bufio.Writer
	acc    accumulator.Accumulator
}

func newElementByFoodReporter(config Config, db hranoprovod.DBNodeMap) *elementByFoodReporter {
	return &elementByFoodReporter{
		config,
		db,
		bufio.NewWriter(config.Output),
		accumulator.NewAccumulator(),
	}
}

func (r *elementByFoodReporter) Process(ln *hranoprovod.LogNode) error {
	singleElement := r.config.SingleElement
	for _, e := range ln.Elements {
		node, found := r.db[e.Name]
		if found {
			for _, repl := range node.Elements {
				if repl.Name == singleElement {
					r.acc.Add(node.Header, repl.Value*e.Value)
				}
			}
		}
	}
	return nil
}

func (r *elementByFoodReporter) Flush() error {
	for name, arr := range r.acc {
		r.printSingleElementByFoodRow(name, arr[accumulator.Positive], arr[accumulator.Negative])
	}
	return r.output.Flush()
}

func (r *elementByFoodReporter) printSingleElementByFoodRow(name string, pos float64, neg float64) {
	format := "%10.2f\t%s\n"
	fmt.Fprintf(r.output, format, pos+neg, name)
}
