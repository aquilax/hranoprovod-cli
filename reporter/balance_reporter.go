package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/aquilax/hranoprovod-cli/accumulator"
	"github.com/aquilax/hranoprovod-cli/shared"
)

type BalanceReporter struct {
	options *Options
	db      *shared.NodeList
	output  io.Writer
	root    *accumulator.TreeNode
}

func NewBalanceReporter(options *Options, db *shared.NodeList, writer io.Writer) *BalanceReporter {
	return &BalanceReporter{
		options,
		db,
		writer,
		accumulator.NewTreeNode("", 0),
	}
}

func (r *BalanceReporter) Process(ln *shared.LogNode) error {
	for _, el := range *ln.Elements {
		r.root.AddDeep(el)
	}
	return nil
}

func (r *BalanceReporter) Flush() error {
	r.printNode(r.root, 0)
	return nil
}

func (r *BalanceReporter) printNode(node *accumulator.TreeNode, level int) {
	for _, key := range node.Keys() {
		child := node.Children[key]
		if len(child.Children) == 0 {
			fmt.Fprintf(r.output, "%10.2f | %s%s\n", child.Sum, strings.Repeat("  ", level), child.Name)
		} else if r.options.CollapseLast && len(child.Children) == 1 && len(child.Children[child.Keys()[0]].Children) == 0 {
			//
			fmt.Fprintf(r.output, "%10.2f | %s%s\n", child.Sum, strings.Repeat("  ", level), child.Name+"/"+child.Children[child.Keys()[0]].Name)
			continue
		} else {
			fmt.Fprintf(r.output, "%10s | %s%s\n", " ", strings.Repeat("  ", level), child.Name)
		}
		r.printNode(child, level+1)
	}
}
