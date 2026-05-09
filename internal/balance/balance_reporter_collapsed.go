package balance

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/aquilax/hranoprovod-cli/v3/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/node"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/tree"
)

type balanceReporterCollapsed struct {
	db     node.DBNodeMap
	output *bufio.Writer
	root   *tree.TreeNode
}

func newBalanceReporterCollapsed(config reporter.Config, db node.DBNodeMap) *balanceReporterCollapsed {
	return &balanceReporterCollapsed{
		db,
		bufio.NewWriter(config.Output),
		tree.NewTreeNode("", 0),
	}
}

func (r *balanceReporterCollapsed) Process(ln *node.LogNode) error {
	for _, el := range ln.Elements {
		el := el
		r.root.AddDeep(el, tree.DefaultCategorySeparator)
	}
	return nil
}

func (r *balanceReporterCollapsed) Flush() error {
	if err := printNodeCollapsed(r.root, 0, r.output); err != nil {
		return err
	}
	return r.output.Flush()
}

func getJump(node *tree.TreeNode) []string {
	if len(node.Children) == 0 {
		return []string{node.Name}
	}
	if len(node.Children) == 1 {
		return append([]string{node.Name}, getJump(node.FirstChild())...)
	}
	return []string{}
}

func printNodeCollapsed(node *tree.TreeNode, level int, output io.Writer) error {
	var err error
	for _, key := range node.Keys() {
		child := node.Children[key]

		jump := getJump(child)
		if len(jump) > 0 {
			if _, err = fmt.Fprintf(output, "%10.2f | %s%s\n", child.Total, strings.Repeat("  ", level), strings.Join(jump, "/")); err != nil {
				return err
			}
			continue
		}
		if _, err = fmt.Fprintf(output, "%10.2f | %s%s\n", child.Total, strings.Repeat("  ", level), child.Name); err != nil {
			return err
		}
		if err = printNodeCollapsed(child, level+1, output); err != nil {
			return err
		}
	}
	return nil
}
