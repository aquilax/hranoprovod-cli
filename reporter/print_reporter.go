package reporter

import (
	"fmt"
	"io"

	"github.com/aquilax/hranoprovod-cli/v2/shared"
)

// PrintReporter outputs log report
type PrintReporter struct {
	config Config
	output io.Writer
}

func NewPrintReporter(config Config) *PrintReporter {
	return &PrintReporter{
		config,
		config.Output,
	}
}

func (pr PrintReporter) Process(ln *shared.LogNode) error {
	var err error
	if _, err = fmt.Fprintf(pr.output, "%s:\n", ln.Time.Format(pr.config.DateFormat)); err != nil {
		return err
	}
	if ln.Metadata != nil {
		for _, md := range *ln.Metadata {
			if md.Name != "" {
				if _, err = fmt.Fprintf(pr.output, "  # %s: %s\n", md.Name, md.Value); err != nil {
					return err
				}
			} else {
				if _, err = fmt.Fprintf(pr.output, "  # %s\n", md.Value); err != nil {
					return err
				}
			}
		}
	}
	for _, el := range ln.Elements {
		if _, err = fmt.Fprintf(pr.output, "  - %s: %0.2f\n", el.Name, el.Value); err != nil {
			return err
		}
	}
	_, err = fmt.Fprintln(pr.output, "")
	return err
}

func (pr PrintReporter) Flush() error {
	return nil
}
