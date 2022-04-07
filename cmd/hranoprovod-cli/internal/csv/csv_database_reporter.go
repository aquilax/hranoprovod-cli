package csv

import (
	"encoding/csv"
	"fmt"

	"github.com/aquilax/hranoprovod-cli/v3/cmd/hranoprovod-cli/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/v3/lib/shared"
)

// CSVReporter outputs report for single food
type CSVDatabaseReporter struct {
	output *csv.Writer
}

// NewCSVReporter creates new CSV reporter
func NewCSVDatabaseReporter(config reporter.Config) CSVDatabaseReporter {
	w := csv.NewWriter(config.Output)
	w.Comma = config.CSVSeparator
	return CSVDatabaseReporter{w}
}

// Process writes single node
func (r CSVDatabaseReporter) Process(n *shared.DBNode) error {
	var err error
	for _, e := range n.Elements {
		if err = r.output.Write([]string{
			n.Header,
			e.Name,
			fmt.Sprintf("%0.2f", e.Value),
		}); err != nil {
			return err
		}
	}
	return nil
}

// Flush flushes the buffer
func (r CSVDatabaseReporter) Flush() error {
	r.output.Flush()
	return r.output.Error()
}