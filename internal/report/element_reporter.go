package report

import (
	"bufio"
	"fmt"

	"github.com/aquilax/hranoprovod-cli/v3/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/element"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/node"
)

type ElementReporter struct {
	list   []element.Element
	output *bufio.Writer
}

func NewElementReporter(c reporter.Config, list []element.Element) *ElementReporter {
	return &ElementReporter{list, bufio.NewWriter(c.Output)}
}

func (er ElementReporter) Process(ln *node.LogNode) error {
	return nil
}
func (er ElementReporter) Flush() error {
	var err error
	for _, el := range er.list {
		if _, err = fmt.Fprintf(er.output, "%0.2f\t%s\n", el.Value, el.Name); err != nil {
			return err
		}
	}
	return er.output.Flush()
}
