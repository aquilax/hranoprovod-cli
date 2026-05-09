package balance

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/aquilax/hranoprovod-cli/v3/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/element"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/node"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/tree"
)

type balanceSingleReporter struct {
	db            node.DBNodeMap
	output        *bufio.Writer
	root          *tree.TreeNode
	total         float64
	singleElement string
	collapse      bool
	collapseLast  bool
}

func newBalanceSingleReporter(config reporter.Config, db node.DBNodeMap) *balanceSingleReporter {
	return &balanceSingleReporter{
		db,
		bufio.NewWriter(config.Output),
		tree.NewTreeNode("", 0),
		0,
		config.SingleElement,
		config.Collapse,
		config.CollapseLast,
	}
}

func (r *balanceSingleReporter) Process(ln *node.LogNode) error {
	for _, el := range ln.Elements {
		repl, found := r.db[el.Name]
		if found {
			for _, repl := range repl.Elements {
				if repl.Name == r.singleElement {
					r.root.AddDeep(element.NewElement(el.Name, repl.Value*el.Value), tree.DefaultCategorySeparator)
					// Add to grand total
					r.total += repl.Value * el.Value
				}
			}
		} else {
			if el.Name == r.singleElement {
				r.root.AddDeep(element.NewElement(el.Name, 0), tree.DefaultCategorySeparator)
			}
		}
	}
	return nil
}

func (r *balanceSingleReporter) Flush() error {
	if r.collapse {
		if err := printNodeCollapsed(r.root, 0, r.output); err != nil {
			return err
		}
	} else {
		if err := printNode(r.root, 0, r.output, r.collapseLast); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(r.output, "%s|\n", strings.Repeat("-", 11)); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(r.output, "%10.2f | %s\n", r.total, r.singleElement); err != nil {
		return err
	}
	return r.output.Flush()
}
