package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/aquilax/hranoprovod-cli/accumulator"
	"github.com/aquilax/hranoprovod-cli/shared"
)

type BalanceSingleReporter struct {
	options *Options
	db      *shared.NodeList
	output  io.Writer
	root    *accumulator.TreeNode
	total   float32
}

func NewBalanceSingleReporter(options *Options, db *shared.NodeList, writer io.Writer) *BalanceSingleReporter {
	return &BalanceSingleReporter{
		options,
		db,
		writer,
		accumulator.NewTreeNode("", 0),
		0,
	}
}

func (r *BalanceSingleReporter) Process(ln *shared.LogNode) error {
	for _, el := range *ln.Elements {
		repl, found := (*r.db)[el.Name]
		if found {
			for _, repl := range *repl.Elements {
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

func (r *BalanceSingleReporter) Flush() error {
	r.printNode(r.root, 0)
	fmt.Printf(strings.Repeat("-", 11) + "|\n")
	fmt.Fprintf(r.output, "%10.2f | %s\n", r.total, r.options.SingleElement)
	return nil
}

func (r *BalanceSingleReporter) printNode(node *accumulator.TreeNode, level int) {
	for _, key := range node.Keys() {
		child := node.Children[key]
		if r.options.CollapseLast && len(child.Children) == 1 && len(child.Children[child.Keys()[0]].Children) == 0 {
			//
			fmt.Fprintf(r.output, "%10.2f | %s%s\n", child.Sum, strings.Repeat("  ", level), child.Name+"/"+child.Children[child.Keys()[0]].Name)
			continue
		} else {
			fmt.Fprintf(r.output, "%10.2f | %s%s\n", child.Sum, strings.Repeat("  ", level), child.Name)
		}
		r.printNode(child, level+1)
	}
}
