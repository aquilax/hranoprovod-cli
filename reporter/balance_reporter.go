package reporter

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/aquilax/hranoprovod-cli/v2"
)

type balanceReporter struct {
	config Config
	db     hranoprovod.DBNodeMap
	output *bufio.Writer
	root   *hranoprovod.TreeNode
}

func newBalanceReporter(config Config, db hranoprovod.DBNodeMap) *balanceReporter {
	return &balanceReporter{
		config,
		db,
		bufio.NewWriter(config.Output),
		hranoprovod.NewTreeNode("", 0),
	}
}

func (r *balanceReporter) Process(ln *hranoprovod.LogNode) error {
	if len(r.config.SingleElement) > 0 {
		for _, el := range ln.Elements {
			repl, found := r.db[el.Name]
			if found {
				for _, repl := range repl.Elements {
					if repl.Name == r.config.SingleElement {
						r.root.AddDeep(hranoprovod.NewElement(el.Name, repl.Value*el.Value), hranoprovod.DefaultCategorySeparator)
					}
				}
			} else {
				if el.Name == r.config.SingleElement {
					r.root.AddDeep(hranoprovod.NewElement(el.Name, 0), hranoprovod.DefaultCategorySeparator)
				}
			}
		}
		return nil
	}

	for _, el := range ln.Elements {
		el := el
		r.root.AddDeep(el, hranoprovod.DefaultCategorySeparator)
	}
	return nil
}

func (r *balanceReporter) Flush() error {
	if r.config.Collapse {
		r.printNodeCollapsed(r.root, 0)
		return nil
	}
	r.printNode(r.root, 0)
	return r.output.Flush()
}

func (r *balanceReporter) printNode(node *hranoprovod.TreeNode, level int) {
	for _, key := range node.Keys() {
		child := node.Children[key]
		if len(child.Children) == 0 {
			fmt.Fprintf(r.output, "%10.2f | %s%s\n", child.Total, strings.Repeat("  ", level), child.Name)
		} else if r.config.CollapseLast && len(child.Children) == 1 && len(child.Children[child.Keys()[0]].Children) == 0 {
			// combine the last two levels
			fmt.Fprintf(r.output, "%10.2f | %s%s\n", child.Total, strings.Repeat("  ", level), child.Name+"/"+child.Children[child.Keys()[0]].Name)
			continue
		} else {
			fmt.Fprintf(r.output, "%10s | %s%s\n", " ", strings.Repeat("  ", level), child.Name)
		}
		r.printNode(child, level+1)
	}
}

func getJump(node *hranoprovod.TreeNode) []string {
	if len(node.Children) > 1 {
		return []string{}
	}
	if len(node.Children) == 0 {
		return []string{node.Name}
	}
	return append([]string{node.Name}, getJump(node.Children[node.Keys()[0]])...)
}

func (r *balanceReporter) printNodeCollapsed(node *hranoprovod.TreeNode, level int) {
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
