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

	elList := NewElements()

	for _, el := range *node.Elements {
		if ndx, exists := elList.Index(el.Name); exists {
			(*elList)[ndx].Val += el.Val
		} else {
			elList.Add(el.Name, el.Val)
		}
	}

	return NewLogNode(t, elList), nil
}
