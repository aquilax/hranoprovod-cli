package shared

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDBnodeMap(t *testing.T) {
	t.Run("Given nodeMap", func(t *testing.T) {
		nl := NewDBNodeMap()
		t.Run("Creates new DBnodeMap", func(t *testing.T) {
			assert.NotNil(t, nl)
		})
		t.Run("Adding new node", func(t *testing.T) {
			node := NewDBNodeFromNode(NewParserNode("test"))
			nl.Push(node)
			t.Run("Increases the number of nodes in the list", func(t *testing.T) {
				assert.Equal(t, 1, len(nl))
			})
		})
	})
}

func TestNewLogNode(t *testing.T) {
	t.Run("Given NewLogNode", func(t *testing.T) {
		now := time.Now()
		elements := NewElements()
		elements.Add("test", 1.22)
		logNode := NewLogNode(now, elements, nil)
		t.Run("Creates new log node with the proper fields", func(t *testing.T) {
			assert.True(t, logNode.Time.Equal(now))
			assert.Equal(t, "test", (logNode.Elements)[0].Name)
			assert.Equal(t, 1.22, (logNode.Elements)[0].Value)
		})
	})
	t.Run("Given Parser Node", func(t *testing.T) {
		t.Run("Creates new node on valid date", func(t *testing.T) {
			node := NewParserNode("2006/01/02")
			logNode, err := NewLogNodeFromElements(time.Now(), node.Elements, nil)
			assert.NotNil(t, logNode)
			assert.Nil(t, err)
		})
	})
}
