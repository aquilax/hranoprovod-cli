package shared

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

// NodeList contains list of general nodes
type NodeList map[string]*ParserNode

// NewNodeList creates new list of general nodes
func NewNodeList() NodeList {
	nl := &NodeList{}
	return *nl
}

// Push adds node to the node list
func (db *NodeList) Push(node *ParserNode) {
	(*db)[(*node).Header] = node
}

// DBNode contains general node data
type DBNode struct {
	Header   string
	Elements Elements
}

func NewDBNodeFromNode(n *ParserNode) *DBNode {
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
