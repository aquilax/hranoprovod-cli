package report

import (
	"bufio"
	"fmt"

	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/reporter"
	shared "github.com/aquilax/hranoprovod-cli/v3"
)

// UnsolvedReporter is unresolved reporter
type UnsolvedReporter struct {
	db     shared.DBNodeMap
	output *bufio.Writer
	list   map[string]bool
}

// NewUnsolvedReporter returns reporter for unresolved elements
func NewUnsolvedReporter(config reporter.Config, db shared.DBNodeMap) *UnsolvedReporter {
	return &UnsolvedReporter{
		db:     db,
		output: bufio.NewWriter(config.Output),
		list:   make(map[string]bool),
	}
}

// Process handles single node
func (r *UnsolvedReporter) Process(ln *shared.LogNode) error {
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
		fmt.Fprintln(r.output, name)
	}
	return r.output.Flush()
}
