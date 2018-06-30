package shared

import (
	"time"
)

// LogNode contains log node data
type LogNode struct {
	Time     time.Time
	Elements *Elements
}

// NewLogNode creates new log node
func NewLogNode(time time.Time, elements *Elements) *LogNode {
	return &LogNode{time, elements}
}

// NewLogNodeFromNode creates new LogNode from Node
func NewLogNodeFromNode(node *Node, dateFormat string) (*LogNode, error) {
	t, err := time.Parse(dateFormat, node.Header)
	if err != nil {
		return nil, err
	}
	return NewLogNode(t, node.Elements), nil
}
