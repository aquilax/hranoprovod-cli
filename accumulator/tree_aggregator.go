package accumulator

import (
	"sort"
	"strings"

	"github.com/aquilax/hranoprovod-cli/v2/shared"
)

// Separator is used to separate the categories
const Separator = "/"

// TreeNode contains data for a single balance tree node
type TreeNode struct {
	Name     string
	Sum      float64
	Children map[string]*TreeNode
}

// NewTreeNode creates new tree node
func NewTreeNode(name string, sum float64) *TreeNode {
	return &TreeNode{
		Name:     name,
		Sum:      sum,
		Children: make(map[string]*TreeNode),
	}
}

// Add appends child node to the current tree node
func (tn *TreeNode) Add(child *TreeNode) *TreeNode {
	if _, ok := tn.Children[child.Name]; !ok {
		tn.Children[child.Name] = child
	} else {
		tn.Children[child.Name].Sum += child.Sum
	}
	return tn.Children[child.Name]
}

// AddDeep adds recursive child nodes to the current node given an element
func (tn *TreeNode) AddDeep(el shared.Element) {
	parent := tn
	names := strings.Split(el.Name, Separator)
	for _, name := range names {
		trn := NewTreeNode(name, el.Val)
		parent = parent.Add(trn)
	}
}

// Keys returns array of children keys
func (tn *TreeNode) Keys() []string {
	keys := make([]string, len(tn.Children))

	i := 0
	for k := range tn.Children {
		keys[i] = k
		i++
	}
	sort.StringSlice.Sort(keys)
	return keys
}
