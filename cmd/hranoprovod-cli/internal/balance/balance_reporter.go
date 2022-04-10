package balance

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/reporter"
	shared "github.com/aquilax/hranoprovod-cli/v3"
)

type balanceReporter struct {
	db           shared.DBNodeMap
	output       *bufio.Writer
	root         *shared.TreeNode
	collapseLast bool
}

func newBalanceReporter(config reporter.Config, db shared.DBNodeMap) *balanceReporter {
	return &balanceReporter{
		db,
		bufio.NewWriter(config.Output),
		shared.NewTreeNode("", 0),
		config.CollapseLast,
	}
}

func (r *balanceReporter) Process(ln *shared.LogNode) error {
	for _, el := range ln.Elements {
		el := el
		r.root.AddDeep(el, shared.DefaultCategorySeparator)
	}
	return nil
}

func (r *balanceReporter) Flush() error {
	printNode(r.root, 0, r.output, r.collapseLast)
	return r.output.Flush()
}

func printNode(node *shared.TreeNode, level int, output io.Writer, collapseLast bool) error {
	for _, key := range node.Keys() {
		child := node.Children[key]
		if len(child.Children) == 0 {
			fmt.Fprintf(output, "%10.2f | %s%s\n", child.Total, strings.Repeat("  ", level), child.Name)
		} else if collapseLast && len(child.Children) == 1 && len(child.Children[child.Keys()[0]].Children) == 0 {
			// combine the last two levels
			fmt.Fprintf(output, "%10.2f | %s%s\n", child.Total, strings.Repeat("  ", level), child.Name+"/"+child.Children[child.Keys()[0]].Name)
			continue
		} else {
			fmt.Fprintf(output, "%10.2f | %s%s\n", child.Total, strings.Repeat("  ", level), child.Name)
		}
		printNode(child, level+1, output, collapseLast)
	}
	return nil
}
