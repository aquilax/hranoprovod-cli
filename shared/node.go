package shared

// Node contains general node data
type Node struct {
	Header   string
	Elements *Elements
}

// NewNode creates new geneal node
func NewNode(header string) *Node {
	return &Node{
		header,
		NewElements(),
	}
}

// NodeList contains list of general nodes
type NodeList map[string]*Node

// NewNodeList creates new list of general nodes
func NewNodeList() *NodeList {
	return &NodeList{}
}

// Push adds node to the node list
func (db *NodeList) Push(node *Node) {
	(*db)[(*node).Header] = node
}
