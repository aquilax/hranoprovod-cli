package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/aquilax/hranoprovod-cli/accumulator"
	"github.com/aquilax/hranoprovod-cli/shared"
)

type balanceSingleReporter struct {
	options *Options
	db      shared.DBNodeList
	output  io.Writer
	root    *accumulator.TreeNode
	total   float64
}

func newBalanceSingleReporter(options *Options, db shared.DBNodeList, writer io.Writer) *balanceSingleReporter {
	return &balanceSingleReporter{
		options,
		db,
		writer,
		accumulator.NewTreeNode("", 0),
		0,
	}
}

func (r *balanceSingleReporter) Process(ln *shared.LogNode) error {
	for _, el := range ln.Elements {
		repl, found := r.db[el.Name]
		if found {
			for _, repl := range repl.Elements {
				if repl.Name == r.options.SingleElement {
					r.root.AddDeep(shared.NewElement(el.Name, repl.Val*el.Val))
					// Add to grand total
					r.total += repl.Val * el.Val
				}
			}
		} else {
			if el.Name == r.options.SingleElement {
				r.root.AddDeep(shared.NewElement(el.Name, 0))
			}
		}
	}
	return nil
}

func (r *balanceSingleReporter) Flush() error {
	if r.options.Collapse {
		r.printNodeCollapsed(r.root, 0)
	} else {
		r.printNode(r.root, 0)
	}

	fmt.Fprintf(r.output, "%s|\n", strings.Repeat("-", 11))
	fmt.Fprintf(r.output, "%10.2f | %s\n", r.total, r.options.SingleElement)
	return nil
}

func (r *balanceSingleReporter) printNode(node *accumulator.TreeNode, level int) {
	for _, key := range node.Keys() {
		child := node.Children[key]
		if r.options.CollapseLast && len(child.Children) == 1 && len(child.Children[child.Keys()[0]].Children) == 0 {
			fmt.Fprintf(r.output, "%10.2f | %s%s\n", child.Sum, strings.Repeat("  ", level), child.Name+"/"+child.Children[child.Keys()[0]].Name)
			continue
		} else {
			fmt.Fprintf(r.output, "%10.2f | %s%s\n", child.Sum, strings.Repeat("  ", level), child.Name)
		}
		r.printNode(child, level+1)
	}
}

func (r *balanceSingleReporter) printNodeCollapsed(node *accumulator.TreeNode, level int) {
	for _, key := range node.Keys() {
		child := node.Children[key]

		jump := getJump(child)
		if len(jump) > 0 {
			fmt.Fprintf(r.output, "%10.2f | %s%s\n", child.Sum, strings.Repeat("  ", level), strings.Join(jump, "/"))
			continue
		}
		fmt.Fprintf(r.output, "%10.2f | %s%s\n", child.Sum, strings.Repeat("  ", level), child.Name)
		r.printNodeCollapsed(child, level+1)
	}
}
