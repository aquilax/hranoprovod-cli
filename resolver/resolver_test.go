package resolver

import (
	"testing"

	"github.com/aquilax/hranoprovod-cli/v2/shared"
	"github.com/tj/assert"
)

func TestResolver(t *testing.T) {
	t.Run("Given nodes database and reslover", func(t *testing.T) {
		nl := shared.NewDBNodeList()
		node1 := shared.NewParserNode("node1")
		node1.Elements.Add("element1", 100)
		node1.Elements.Add("element2", 200)
		nl.Push(shared.NewDBNodeFromNode(node1))
		node2 := shared.NewParserNode("node2")
		node2.Elements.Add("node1", 2)
		nl.Push(shared.NewDBNodeFromNode(node2))
		resolver := NewResolver(nl, 1)
		t.Run("Resolve resolves the database", func(t *testing.T) {
			resolver.Resolve()
			t.Run("Elements are resolved", func(t *testing.T) {
				n1 := nl["node1"]
				assert.Equal(t, "element1", n1.Elements[0].Name)
				assert.Equal(t, 100., n1.Elements[0].Val)
				assert.Equal(t, "element2", n1.Elements[1].Name)
				assert.Equal(t, 200., n1.Elements[1].Val)
				n2 := nl["node2"]
				assert.Equal(t, "element1", n2.Elements[0].Name)
				assert.Equal(t, 200., n2.Elements[0].Val)
				assert.Equal(t, "element2", n2.Elements[1].Name)
				assert.Equal(t, 400., n2.Elements[1].Val)
			})
		})
	})
}
