package report

import (
	"bufio"
	"fmt"

	"github.com/aquilax/hranoprovod-cli/v3/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/node"
)

// UnsolvedReporter is unresolved reporter
type UnsolvedReporter struct {
	db     node.DBNodeMap
	output *bufio.Writer
	list   map[string]bool
}

// NewUnsolvedReporter returns reporter for unresolved elements
func NewUnsolvedReporter(config reporter.Config, db node.DBNodeMap) *UnsolvedReporter {
	return &UnsolvedReporter{
		db:     db,
		output: bufio.NewWriter(config.Output),
		list:   make(map[string]bool),
	}
}

// Process handles single node
func (r *UnsolvedReporter) Process(ln *node.LogNode) error {
	for _, e := range ln.Elements {
		_, found := r.db[e.Name]
		if !found {
			r.list[e.Name] = true
		}
	}
	return nil
}

// Flush flushes the report
func (r *UnsolvedReporter) Flush() error {
	for name := range r.list {
		if _, err := fmt.Fprintln(r.output, name); err != nil {
			return err
		}
	}
	return r.output.Flush()
}
