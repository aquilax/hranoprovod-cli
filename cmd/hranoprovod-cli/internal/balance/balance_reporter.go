package balance

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/aquilax/hranoprovod-cli/v3/cmd/hranoprovod-cli/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/v3/lib/shared"
)

type balanceReporter struct {
	config reporter.Config
	db     shared.DBNodeMap
	output *bufio.Writer
	root   *shared.TreeNode
}

// NewBalanceReporter returns balance reporter
func NewBalanceReporter(config reporter.Config, db shared.DBNodeMap) reporter.Reporter {
	if len(config.SingleElement) > 0 {
		return newBalanceSingleReporter(config, db)
	}
	return newBalanceReporter(config, db)
}

func newBalanceReporter(config reporter.Config, db shared.DBNodeMap) *balanceReporter {
	return &balanceReporter{
		config,
		db,
		bufio.NewWriter(config.Output),
		shared.NewTreeNode("", 0),
	}
}

func (r *balanceReporter) Process(ln *shared.LogNode) error {
	if len(r.config.SingleElement) > 0 {
		for _, el := range ln.Elements {
			repl, found := r.db[el.Name]
			if found {
				for _, repl := range repl.Elements {
					if repl.Name == r.config.SingleElement {
						r.root.AddDeep(shared.NewElement(el.Name, repl.Value*el.Value), shared.DefaultCategorySeparator)
					}
				}
			} else {
				if el.Name == r.config.SingleElement {
					r.root.AddDeep(shared.NewElement(el.Name, 0), shared.DefaultCategorySeparator)
				}
			}
		}
		return nil
	}

	for _, el := range ln.Elements {
		el := el
		r.root.AddDeep(el, shared.DefaultCategorySeparator)
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

func (r *balanceReporter) printNode(node *shared.TreeNode, level int) {
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

func getJump(node *shared.TreeNode) []string {
	if len(node.Children) > 1 {
		return []string{}
	}
	if len(node.Children) == 0 {
		return []string{node.Name}
	}
	return append([]string{node.Name}, getJump(node.Children[node.Keys()[0]])...)
}

func (r *balanceReporter) printNodeCollapsed(node *shared.TreeNode, level int) {
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
