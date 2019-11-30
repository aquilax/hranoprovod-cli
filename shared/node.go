package shared

import (
	"time"
)

// ParserNode contains general node data
type ParserNode struct {
	Header   string
	Elements Elements
}

// NewParserNode creates new geneal node
func NewParserNode(header string) *ParserNode {
	return &ParserNode{
		header,
		NewElements(),
	}
}

// DBNode contains general node data
type DBNode struct {
	Header   string
	Elements Elements
}

// NewDBNodeFromNode creates DB Node from Parser node
func NewDBNodeFromNode(n *ParserNode) *DBNode {
	dbn := DBNode(*n)
	return &dbn
}

// DBNodeList contains list of general nodes
type DBNodeList map[string]*DBNode

// NewDBNodeList creates new list of general nodes
func NewDBNodeList() DBNodeList {
	return DBNodeList{}
}

// Push adds node to the node list
func (db *DBNodeList) Push(node *DBNode) {
	(*db)[(*node).Header] = node
}

// LogNode contains log node data
type LogNode struct {
	Time     time.Time
	Elements Elements
}

// NewLogNode creates new log node
func NewLogNode(time time.Time, elements Elements) *LogNode {
	return &LogNode{time, elements}
}

// GetHeaderTimeFromNode tries to parse node's time from the header and returns it as time.Time
func ParseTime(header string, dateFormat string) (time.Time, error) {
	return time.Parse(dateFormat, header)
}

// NewLogNodeFromElements creates new LogNode from ParserNode elements and time
func NewLogNodeFromElements(time time.Time, elements Elements) (*LogNode, error) {
	elList := NewElements()

	for _, el := range elements {
		if ndx, exists := elList.Index(el.Name); exists {
			elList[ndx].Val += el.Val
		} else {
			elList.Add(el.Name, el.Val)
		}
	}

	return NewLogNode(time, elList), nil
}
