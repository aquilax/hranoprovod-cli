package resolver

import (
	"github.com/aquilax/hranoprovod-cli/v2/shared"
)

// Resolver contains the resolver data
type Resolver struct {
	db       shared.DBNodeList
	maxDepth int
}

// NewResolver creates new resolver
func NewResolver(db shared.DBNodeList, maxDepth int) Resolver {
	return Resolver{db, maxDepth}
}

// Resolve resolves the current database
func (r Resolver) Resolve() {
	for name := range r.db {
		r.resolveNode(name, 0)
	}
}

func (r Resolver) resolveNode(name string, level int) {
	if level >= r.maxDepth {
		return
	}

	node, exists := r.db[name]
	if !exists {
		return
	}

	nel := shared.NewElements()

	for _, e := range node.Elements {
		r.resolveNode(e.Name, level+1)
		foundNode, exists := r.db[e.Name]
		if exists {
			nel.SumMerge(foundNode.Elements, e.Value)
		} else {
			var tm shared.Elements
			tm.Add(e.Name, e.Value)
			nel.SumMerge(tm, 1)
		}
	}
	nel.Sort()
	r.db[name].Elements = nel
}
