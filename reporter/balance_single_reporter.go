package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/aquilax/hranoprovod-cli/v2/accumulator"
	"github.com/aquilax/hranoprovod-cli/v2/shared"
)

type balanceSingleReporter struct {
	config Config
	db     shared.DBNodeList
	output io.Writer
	root   *accumulator.TreeNode
	total  float64
}

func newBalanceSingleReporter(config Config, db shared.DBNodeList) *balanceSingleReporter {
	return &balanceSingleReporter{
		config,
		db,
		config.Output,
		accumulator.NewTreeNode("", 0),
		0,
	}
}

func (r *balanceSingleReporter) Process(ln *shared.LogNode) error {
	for _, el := range ln.Elements {
		repl, found := r.db[el.Name]
		if found {
			for _, repl := range repl.Elements {
				if repl.Name == r.config.SingleElement {
					r.root.AddDeep(shared.NewElement(el.Name, repl.Value*el.Value), accumulator.DefaultCategorySeparator)
					// Add to grand total
					r.total += repl.Value * el.Value
				}
			}
		} else {
			if el.Name == r.config.SingleElement {
				r.root.AddDeep(shared.NewElement(el.Name, 0), accumulator.DefaultCategorySeparator)
			}
		}
	}
	return nil
}

func (r *balanceSingleReporter) Flush() error {
	if r.config.Collapse {
		r.printNodeCollapsed(r.root, 0)
	} else {
		r.printNode(r.root, 0)
	}

	fmt.Fprintf(r.output, "%s|\n", strings.Repeat("-", 11))
	fmt.Fprintf(r.output, "%10.2f | %s\n", r.total, r.config.SingleElement)
	return nil
}

func (r *balanceSingleReporter) printNode(node *accumulator.TreeNode, level int) {
	for _, key := range node.Keys() {
		child := node.Children[key]
		if r.config.CollapseLast && len(child.Children) == 1 && len(child.Children[child.Keys()[0]].Children) == 0 {
			fmt.Fprintf(r.output, "%10.2f | %s%s\n", child.Total, strings.Repeat("  ", level), child.Name+"/"+child.Children[child.Keys()[0]].Name)
			continue
		} else {
			fmt.Fprintf(r.output, "%10.2f | %s%s\n", child.Total, strings.Repeat("  ", level), child.Name)
		}
		r.printNode(child, level+1)
	}
}

func (r *balanceSingleReporter) printNodeCollapsed(node *accumulator.TreeNode, level int) {
	for _, key := range node.Keys() {
		child := node.Children[key]

		jump := getJump(child)
		if len(jump) > 0 {
			fmt.Fprintf(r.output, "%10.2f | %s%s\n", child.Total, strings.Repeat("  ", level), strings.Join(jump, "/"))
			continue
		}
		fmt.Fprintf(r.output, "%10.2f | %s%s\n", child.Total, strings.Repeat("  ", level), child.Name)
		r.printNodeCollapsed(child, level+1)
	}
}
