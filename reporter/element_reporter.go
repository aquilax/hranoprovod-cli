package reporter

import (
	"bufio"
	"fmt"

	"github.com/aquilax/hranoprovod-cli/v2"
)

type ElementReporter struct {
	config Config
	list   []hranoprovod.Element
	output *bufio.Writer
}

func NewElementReporter(c Config, list []hranoprovod.Element) *ElementReporter {
	return &ElementReporter{c, list, bufio.NewWriter(c.Output)}
}

func (er ElementReporter) Process(ln *hranoprovod.LogNode) error {
	return nil
}
func (er ElementReporter) Flush() error {
	var err error
	for _, el := range er.list {
		if _, err = fmt.Fprintf(er.config.Output, "%0.2f\t%s\n", el.Value, el.Name); err != nil {
			return err
		}
	}
	return er.output.Flush()
}
