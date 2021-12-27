package resolver

import (
	"testing"

	"github.com/aquilax/hranoprovod-cli/v2/shared"
	"github.com/stretchr/testify/assert"
)

func TestResolver(t *testing.T) {
	t.Run("Given nodes database and reslover", func(t *testing.T) {
		nl := shared.DBNodeList{
			"node1": &shared.DBNode{
				Header: "node1",
				Elements: shared.Elements{
					shared.Element{Name: "element1", Value: 100},
					shared.Element{Name: "element2", Value: 200},
				},
			},
			"node2": &shared.DBNode{
				Header: "node2",
				Elements: shared.Elements{
					shared.Element{Name: "node1", Value: 2},
				},
			},
		}
		resolver := NewResolver(nl, Config{10})
		t.Run("Resolve resolves the database", func(t *testing.T) {
			resolver.Resolve()
			t.Run("Elements are resolved", func(t *testing.T) {
				n1 := nl["node1"]
				assert.Equal(t, "element1", n1.Elements[0].Name)
				assert.Equal(t, 100., n1.Elements[0].Value)
				assert.Equal(t, "element2", n1.Elements[1].Name)
				assert.Equal(t, 200., n1.Elements[1].Value)
				n2 := nl["node2"]
				assert.Equal(t, "element1", n2.Elements[0].Name)
				assert.Equal(t, 200., n2.Elements[0].Value)
				assert.Equal(t, "element2", n2.Elements[1].Name)
				assert.Equal(t, 400., n2.Elements[1].Value)
			})
		})
	})
}
