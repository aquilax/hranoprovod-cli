package hranoprovod

import (
	"sort"
	"strings"
)

// DefaultCategorySeparator the default string used to separate categories
const DefaultCategorySeparator = "/"

// TreeNode contains data for a single balance tree node
type TreeNode struct {
	Name     string
	Total    float64
	Children map[string]*TreeNode
}

// NewTreeNode creates new tree node
func NewTreeNode(name string, sum float64) *TreeNode {
	return &TreeNode{
		Name:     name,
		Total:    sum,
		Children: make(map[string]*TreeNode),
	}
}

// Add appends child node to the current tree node
func (tn *TreeNode) Add(child *TreeNode) *TreeNode {
	if _, ok := tn.Children[child.Name]; !ok {
		tn.Children[child.Name] = child
	} else {
		tn.Children[child.Name].Total += child.Total
	}
	return tn.Children[child.Name]
}

// AddDeep adds recursive child nodes to the current node given an element
func (tn *TreeNode) AddDeep(el Element, separator string) {
	parent := tn
	names := strings.Split(el.Name, separator)
	for _, name := range names {
		trn := NewTreeNode(name, el.Value)
		parent = parent.Add(trn)
	}
}

// Keys returns sorted array of children keys
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

// FirstChild returns the first child of a TreeNode
func (tn *TreeNode) FirstChild() *TreeNode {
	if len(tn.Children) == 0 {
		return nil
	}
	return tn.Children[tn.Keys()[0]]
}
