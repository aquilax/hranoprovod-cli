package reporter

import (
	"bufio"
	"fmt"
	"sort"

	"github.com/aquilax/hranoprovod-cli/v2/lib/shared"
)

// elementByFoodReporter outputs report for single element groupped by food
type elementByFoodReporter struct {
	config Config
	db     shared.DBNodeMap
	output *bufio.Writer
	acc    shared.Accumulator
}

func newElementByFoodReporter(config Config, db shared.DBNodeMap) *elementByFoodReporter {
	return &elementByFoodReporter{
		config,
		db,
		bufio.NewWriter(config.Output),
		shared.NewAccumulator(),
	}
}

func (r *elementByFoodReporter) Process(ln *shared.LogNode) error {
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
	keys := make([]string, len(r.acc))
	i := 0
	for name := range r.acc {
		keys[i] = name
		i++
	}
	sort.Strings(keys)
	for _, name := range keys {
		arr := r.acc[name]
		r.printSingleElementByFoodRow(name, arr[shared.Positive], arr[shared.Negative])
	}
	return r.output.Flush()
}

func (r *elementByFoodReporter) printSingleElementByFoodRow(name string, pos float64, neg float64) {
	format := "%10.2f\t%s\n"
	fmt.Fprintf(r.output, format, pos+neg, name)
}
