package reporter

import (
	"fmt"

	"github.com/aquilax/hranoprovod-cli/v2/shared"
)

type ElementReporter struct {
	config Config
	list   []shared.Element
}

func NewElementReporter(c Config, list []shared.Element) *ElementReporter {
	return &ElementReporter{c, list}
}

func (er ElementReporter) Process(ln *shared.LogNode) error {
	return nil
}
func (er ElementReporter) Flush() error {
	var err error
	for _, el := range er.list {
		if _, err = fmt.Fprintf(er.config.Output, "%0.2f\t%s\n", el.Value, el.Name); err != nil {
			return err
		}
	}
	return err
}
