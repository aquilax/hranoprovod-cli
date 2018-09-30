package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/aquilax/hranoprovod-cli/accumulator"
	"github.com/aquilax/hranoprovod-cli/shared"
)

type balanceReporter struct {
	options *Options
	db      *shared.NodeList
	output  io.Writer
	root    *accumulator.TreeNode
}

func newBalanceReporter(options *Options, db *shared.NodeList, writer io.Writer) *balanceReporter {
	return &balanceReporter{
		options,
		db,
		writer,
		accumulator.NewTreeNode("", 0),
	}
}

func (r *balanceReporter) Process(ln *shared.LogNode) error {
	if len(r.options.SingleElement) > 0 {
		for _, el := range *ln.Elements {
			repl, found := (*r.db)[el.Name]
			if found {
				for _, repl := range *repl.Elements {
					if repl.Name == r.options.SingleElement {
						r.root.AddDeep(shared.NewElement(el.Name, repl.Val*el.Val))
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

	for _, el := range *ln.Elements {
		r.root.AddDeep(el)
	}
	return nil
}

func (r *balanceReporter) Flush() error {
	if r.options.Collapse {
		r.printNodeCollapsed(r.root, 0)
		return nil
	}
	r.printNode(r.root, 0)
	return nil
}

func (r *balanceReporter) printNode(node *accumulator.TreeNode, level int) {
	for _, key := range node.Keys() {
		child := node.Children[key]
		if len(child.Children) == 0 {
			fmt.Fprintf(r.output, "%10.2f | %s%s\n", child.Sum, strings.Repeat("  ", level), child.Name)
		} else if r.options.CollapseLast && len(child.Children) == 1 && len(child.Children[child.Keys()[0]].Children) == 0 {
			// combine the last two levels
			fmt.Fprintf(r.output, "%10.2f | %s%s\n", child.Sum, strings.Repeat("  ", level), child.Name+"/"+child.Children[child.Keys()[0]].Name)
			continue
		} else {
			fmt.Fprintf(r.output, "%10s | %s%s\n", " ", strings.Repeat("  ", level), child.Name)
		}
		r.printNode(child, level+1)
	}
}

func getJump(node *accumulator.TreeNode) []string {
	if len(node.Children) > 1 {
		return []string{}
	}
	if len(node.Children) == 0 {
		return []string{node.Name}
	}
	return append([]string{node.Name}, getJump(node.Children[node.Keys()[0]])...)
}

func (r *balanceReporter) printNodeCollapsed(node *accumulator.TreeNode, level int) {
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
