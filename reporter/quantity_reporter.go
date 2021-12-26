package reporter

import (
	"fmt"
	"io"
	"sort"

	"github.com/aquilax/hranoprovod-cli/v2/shared"
)

type QuantityReporter struct {
	ascending   bool
	accumulator map[string]float64
	output      io.Writer
}

func NewQuantityReporter(ascending bool, writer io.Writer) QuantityReporter {
	return QuantityReporter{
		ascending,
		make(map[string]float64),
		writer,
	}
}

// Process writes single node
func (r QuantityReporter) Process(ln *shared.LogNode) error {
	for _, e := range ln.Elements {
		if _, ok := r.accumulator[e.Name]; !ok {
			r.accumulator[e.Name] = 0
		}
		r.accumulator[e.Name] += e.Value
	}
	return nil
}

type SortTuple struct {
	name  string
	value float64
}

// Flush does nothing
func (r QuantityReporter) Flush() error {
	sortable := make([]SortTuple, 0, len(r.accumulator))
	for k, v := range r.accumulator {
		sortable = append(sortable, SortTuple{k, v})
	}

	if r.ascending {
		sort.SliceStable(sortable, func(i, j int) bool {
			return sortable[i].value > sortable[j].value
		})
	} else {
		sort.SliceStable(sortable, func(i, j int) bool {
			return sortable[i].value < sortable[j].value
		})
	}
	for _, el := range sortable {
		fmt.Printf("%0.2f\t%s\n", el.value, el.name)
	}
	return nil
}
