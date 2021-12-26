package reporter

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/aquilax/hranoprovod-cli/v2/shared"
)

// CSVReporter outputs report for single food
type CSVReporter struct {
	options *Options
	output  *csv.Writer
}

// NewCSVReporter creates new CSV reporter
func NewCSVReporter(options *Options, writer io.Writer) CSVReporter {
	w := csv.NewWriter(writer)
	w.Comma = ';'
	return CSVReporter{
		options,
		w,
	}
}

// Process writes single node
func (r CSVReporter) Process(ln *shared.LogNode) error {
	var err error
	for _, e := range ln.Elements {
		if err = r.output.Write([]string{
			ln.Time.Format("2006-01-02"),
			e.Name,
			fmt.Sprintf("%0.2f", e.Value),
		}); err != nil {
			return err
		}
	}
	return nil
}

// Flush does nothing
func (r CSVReporter) Flush() error {
	r.output.Flush()
	return nil
}
