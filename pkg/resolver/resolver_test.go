package resolver

import (
	"fmt"
	"testing"

	"github.com/aquilax/hranoprovod-cli/v3/pkg/element"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/node"
	"github.com/stretchr/testify/assert"
)

func getTestnodeMap() node.DBNodeMap {
	return node.DBNodeMap{
		"node1": &node.DBNode{
			Header: "node1",
			Elements: element.Elements{
				element.Element{Name: "element1", Value: 100},
				element.Element{Name: "element2", Value: 200},
			},
		},
		"node2": &node.DBNode{
			Header: "node2",
			Elements: element.Elements{
				element.Element{Name: "node1", Value: 2},
			},
		},
	}
}

func getSizeNnodeMap(n int) node.DBNodeMap {
	var nl = node.DBNodeMap{}
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("node-%d", i)
		nl[name] = &node.DBNode{
			Header: name,
			Elements: element.Elements{
				element.Element{Name: fmt.Sprintf("node-%d", i+1), Value: float64(i + 1)},
				element.Element{Name: fmt.Sprintf("node-%d", i+2), Value: float64(i + 2)},
				element.Element{Name: fmt.Sprintf("node-%d", i+3), Value: float64(i + 3)},
			},
		}
	}
	return nl
}

func TestResolver(t *testing.T) {
	t.Run("Given nodes database and reslover", func(t *testing.T) {
		nl := getTestnodeMap()
		resolver := NewResolver(nl, Config{10})
		t.Run("Resolve resolves the database", func(t *testing.T) {
			assert.Nil(t, resolver.Resolve())
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

func TestResolver_Resolve(t *testing.T) {
	t.Run("Given nodes database and reslover", func(t *testing.T) {
		nl := getTestnodeMap()
		t.Run("Resolve resolves the database", func(t *testing.T) {
			nl, err := Resolve(Config{10}, nl)
			assert.Equal(t, err, nil)
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

func BenchmarkResolve(b *testing.B) {
	nl := getSizeNnodeMap(100)
	var err error
	for n := 0; n < b.N; n++ {
		_, err = Resolve(Config{10}, nl)
	}
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkResolverResolve(b *testing.B) {
	nl := getSizeNnodeMap(100)
	resolver := NewResolver(nl, Config{10})
	var err error
	for n := 0; n < b.N; n++ {
		err = resolver.Resolve()
	}
	if err != nil {
		b.Fatal(err)
	}
}
