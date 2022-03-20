package hranoprovod

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreeNode_AddDeep(t *testing.T) {
	tests := []struct {
		name   string
		caller func() *TreeNode
		want   *TreeNode
	}{
		{
			"spreads correctly the values to categories",
			func() *TreeNode {
				tn := NewTreeNode("root", 0)
				tn.AddDeep(NewElement("one/two", 10), DefaultCategorySeparator)
				return tn
			},
			&TreeNode{
				Name:  "root",
				Total: 0,
				Children: map[string]*TreeNode{
					"one": {
						Name:  "one",
						Total: 10,
						Children: map[string]*TreeNode{
							"two": {
								Name:     "two",
								Total:    10,
								Children: map[string]*TreeNode{},
							},
						},
					},
				}},
		},
		{
			"accumulates correctly the values",
			func() *TreeNode {
				tn := NewTreeNode("root", 0)
				tn.AddDeep(NewElement("one", 10), DefaultCategorySeparator)
				tn.AddDeep(NewElement("one/two", 20), DefaultCategorySeparator)
				return tn
			},
			&TreeNode{
				Name:  "root",
				Total: 0,
				Children: map[string]*TreeNode{
					"one": {
						Name:  "one",
						Total: 30,
						Children: map[string]*TreeNode{
							"two": {
								Name:     "two",
								Total:    20,
								Children: map[string]*TreeNode{},
							},
						},
					},
				}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tn := tt.caller()
			if !reflect.DeepEqual(tn, tt.want) {
				t.Errorf("want %+v got %+v", tt.want, tt)
			}
			assert.Equal(t, tt.want, tn)
		})
	}
}

func TestTreeNode_Keys(t *testing.T) {
	tests := []struct {
		name string
		tn   *TreeNode
		want []string
	}{
		{
			"returns sorted list of keys",
			&TreeNode{
				Name:  "root",
				Total: 0,
				Children: map[string]*TreeNode{
					"999":  {Name: "999", Total: 0, Children: map[string]*TreeNode{}},
					"zzzz": {Name: "zzzz", Total: 0, Children: map[string]*TreeNode{}},
					"a999": {Name: "a999", Total: 0, Children: map[string]*TreeNode{}},
				},
			},
			[]string{"999", "a999", "zzzz"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tn.Keys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TreeNode.Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}
