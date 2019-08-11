package shared

// Node contains general node data
type Node struct {
	Header   string
	Elements Elements
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
func NewNodeList() NodeList {
	nl := &NodeList{}
	return *nl
}

// Push adds node to the node list
func (db *NodeList) Push(node *Node) {
	(*db)[(*node).Header] = node
}

// DBNode contains general node data
type DBNode struct {
	Header   string
	Elements Elements
}

func NewDBNodeFromNode(n *Node) *DBNode {
	dbn := DBNode(*n)
	return &dbn
}

// DBNodeList contains list of general nodes
type DBNodeList map[string]*DBNode

// NewDBNodeList creates new list of general nodes
func NewDBNodeList() DBNodeList {
	dbnl := &DBNodeList{}
	return *dbnl
}

// Push adds node to the node list
func (db *DBNodeList) Push(node *DBNode) {
	(*db)[(*node).Header] = node
}
