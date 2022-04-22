package csv

import (
	"encoding/csv"
	"fmt"

	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/reporter"
	shared "github.com/aquilax/hranoprovod-cli/v3"
)

const DefaultOutputTimeFormat = "2006-01-02"
const DefaultCSVSeparator = ','

type CSVConfig struct {
	reporter.CommonConfig
	CSVSeparator     rune
	OutputTimeFormat string
}

func NewCSVConfig(c reporter.CommonConfig) CSVConfig {
	return CSVConfig{
		CommonConfig:     c,
		CSVSeparator:     DefaultCSVSeparator,
		OutputTimeFormat: DefaultOutputTimeFormat,
	}

}

// CSVReporter outputs report for single food
type CSVReporter struct {
	output           *csv.Writer
	outputTimeFormat string
}

// NewCSVReporter creates new CSV reporter
func NewCSVReporter(config CSVConfig) CSVReporter {
	w := csv.NewWriter(config.Output)
	w.Comma = config.CSVSeparator
	return CSVReporter{
		w,
		config.OutputTimeFormat,
	}
}

// Process writes single node
func (r CSVReporter) Process(ln *shared.LogNode) error {
	var err error
	for _, e := range ln.Elements {
		if err = r.output.Write([]string{
			ln.Time.Format(r.outputTimeFormat),
			e.Name,
			fmt.Sprintf("%0.3f", e.Value),
		}); err != nil {
			return err
		}
	}
	return nil
}

// Flush flushes the buffer
func (r CSVReporter) Flush() error {
	r.output.Flush()
	return r.output.Error()
}
