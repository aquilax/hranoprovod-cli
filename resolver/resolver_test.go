package resolver

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestnodeMap() hranoprovod.DBNodeMap {
	return hranoprovod.DBNodeMap{
		"node1": &hranoprovod.DBNode{
			Header: "node1",
			Elements: hranoprovod.Elements{
				hranoprovod.Element{Name: "element1", Value: 100},
				hranoprovod.Element{Name: "element2", Value: 200},
			},
		},
		"node2": &hranoprovod.DBNode{
			Header: "node2",
			Elements: hranoprovod.Elements{
				hranoprovod.Element{Name: "node1", Value: 2},
			},
		},
	}
}

func getSizeNnodeMap(n int) hranoprovod.DBNodeMap {
	var nl = hranoprovod.DBNodeMap{}
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("node-%d", i)
		nl[name] = &hranoprovod.DBNode{
			Header: name,
			Elements: hranoprovod.Elements{
				hranoprovod.Element{Name: fmt.Sprintf("node-%d", i+1), Value: float64(i + 1)},
				hranoprovod.Element{Name: fmt.Sprintf("node-%d", i+2), Value: float64(i + 2)},
				hranoprovod.Element{Name: fmt.Sprintf("node-%d", i+3), Value: float64(i + 3)},
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
	for n := 0; n < b.N; n++ {
		Resolve(Config{10}, nl)
	}
}

func BenchmarkResolverResolve(b *testing.B) {
	nl := getSizeNnodeMap(100)
	resolver := NewResolver(nl, Config{10})
	for n := 0; n < b.N; n++ {
		resolver.Resolve()
	}
}
