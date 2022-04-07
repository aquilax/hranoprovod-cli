package report

import (
	"bufio"
	"fmt"
	"sort"

	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/reporter"
	shared "github.com/aquilax/hranoprovod-cli/v3"
)

type QuantityReporter struct {
	descending  bool
	accumulator map[string]float64
	output      *bufio.Writer
}

func NewQuantityReporter(config reporter.Config, descending bool) QuantityReporter {
	return QuantityReporter{
		descending,
		make(map[string]float64),
		bufio.NewWriter(config.Output),
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

	if r.descending {
		sort.SliceStable(sortable, func(i, j int) bool {
			return sortable[i].value > sortable[j].value
		})
	} else {
		sort.SliceStable(sortable, func(i, j int) bool {
			return sortable[i].value < sortable[j].value
		})
	}
	for _, el := range sortable {
		el := el
		fmt.Fprintf(r.output, "%0.2f\t%s\n", el.value, el.name)
	}
	return r.output.Flush()
}
