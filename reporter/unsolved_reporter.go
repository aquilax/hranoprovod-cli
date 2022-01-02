package reporter

import (
	"fmt"
	"io"

	"github.com/aquilax/hranoprovod-cli/v2/shared"
)

// UnsolvedReporter is unresolved reporter
type UnsolvedReporter struct {
	config Config
	db     shared.DBNodeList
	output io.Writer
	list   map[string]bool
}

// NewUnsolvedReporter returns reporter for unresolved elements
func NewUnsolvedReporter(config Config, db shared.DBNodeList) *UnsolvedReporter {
	return &UnsolvedReporter{
		config,
		db,
		config.Output,
		make(map[string]bool),
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
	return nil
}
