package reporter

import (
	"bufio"
	"fmt"

	"github.com/aquilax/hranoprovod-cli/v2"
)

// UnsolvedReporter is unresolved reporter
type UnsolvedReporter struct {
	config Config
	db     hranoprovod.DBNodeMap
	output *bufio.Writer
	list   map[string]bool
}

// NewUnsolvedReporter returns reporter for unresolved elements
func NewUnsolvedReporter(config Config, db hranoprovod.DBNodeMap) *UnsolvedReporter {
	return &UnsolvedReporter{
		config,
		db,
		bufio.NewWriter(config.Output),
		make(map[string]bool),
	}
}

// Process handles single node
func (r *UnsolvedReporter) Process(ln *hranoprovod.LogNode) error {
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
