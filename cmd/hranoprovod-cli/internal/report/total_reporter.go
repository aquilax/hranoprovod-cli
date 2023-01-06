package report

import (
	"bufio"
	"fmt"
	"io"
	"sort"

	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/reporter"
	shared "github.com/aquilax/hranoprovod-cli/v3"
)

type TotalReporter struct {
	output *bufio.Writer
	db     shared.DBNodeMap
	acc    shared.Accumulator
}

func NewTotalReporter(c reporter.Config, db shared.DBNodeMap) *TotalReporter {
	return &TotalReporter{
		output: bufio.NewWriter(c.Output),
		db:     db,
		acc:    shared.NewAccumulator(),
	}
}

func (tr TotalReporter) Process(ln *shared.LogNode) error {
	for _, element := range ln.Elements {
		if repl, found := tr.db[element.Name]; found {
			for _, repl := range repl.Elements {
				res := repl.Value * element.Value
				tr.acc.Add(repl.Name, res)
			}
		} else {
			tr.acc.Add(element.Name, element.Value)
		}
	}
	return nil
}

func (tr TotalReporter) Flush() error {
	if len(tr.acc) > 0 {
		ss := make([]string, len(tr.acc))
		i := 0
		for name := range tr.acc {
			ss[i] = name
			i++
		}
		sort.Strings(ss)
		fmt.Fprintf(tr.output, "%12s  %12s  %12s  %s\n", "positive", "negative", "sum", "element")
		for _, name := range ss {
			arr := tr.acc[name]
			if err := printTotalRow(name, arr[shared.Positive], arr[shared.Negative], tr.output); err != nil {
				return err
			}
		}
	}
	return tr.output.Flush()
}

func printTotalRow(name string, positive, negative float64, output io.Writer) error {
	_, err := fmt.Fprintf(output, "%12.2f  %12.2f  %12.2f  %s\n", positive, negative, positive+negative, name)
	return err
}
