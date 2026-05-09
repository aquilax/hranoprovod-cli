package balance

import (
	"bufio"
	"io"
	"strings"

	"github.com/aquilax/hranoprovod-cli/v3/internal/errwriter"
	"github.com/aquilax/hranoprovod-cli/v3/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/node"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/tree"
)

type balanceReporter struct {
	db           node.DBNodeMap
	output       *bufio.Writer
	root         *tree.TreeNode
	collapseLast bool
}

func newBalanceReporter(config reporter.Config, db node.DBNodeMap) *balanceReporter {
	return &balanceReporter{
		db,
		bufio.NewWriter(config.Output),
		tree.NewTreeNode("", 0),
		config.CollapseLast,
	}
}

func (r *balanceReporter) Process(ln *node.LogNode) error {
	for _, el := range ln.Elements {
		el := el
		r.root.AddDeep(el, tree.DefaultCategorySeparator)
	}
	return nil
}

func (r *balanceReporter) Flush() error {
	if err := printNode(r.root, 0, r.output, r.collapseLast); err != nil {
		return err
	}
	return r.output.Flush()
}

func printNode(node *tree.TreeNode, level int, output io.Writer, collapseLast bool) error {
	writer := errwriter.New(output)
	for _, key := range node.Keys() {
		child := node.Children[key]
		if len(child.Children) == 0 {
			writer.Fprintf("%10.2f | %s%s\n", child.Total, strings.Repeat("  ", level), child.Name)
		} else if collapseLast && len(child.Children) == 1 && len(child.FirstChild().Children) == 0 {
			// combine the last two levels
			writer.Fprintf("%10.2f | %s%s\n", child.Total, strings.Repeat("  ", level), child.Name+"/"+child.FirstChild().Name)
			continue
		} else {
			writer.Fprintf("%10.2f | %s%s\n", child.Total, strings.Repeat("  ", level), child.Name)
		}
		if err := printNode(child, level+1, output, collapseLast); err != nil {
			return err
		}
	}
	return writer.Err
}
