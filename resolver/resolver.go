package resolver

import (
	"fmt"

	"github.com/aquilax/hranoprovod-cli/v2/shared"
)

type Config struct {
	MaxDepth int
}

// Resolver contains the resolver data
type Resolver struct {
	db     shared.DBNodeList
	config Config
}

// NewResolver creates new resolver
func NewResolver(db shared.DBNodeList, c Config) Resolver {
	return Resolver{db, c}
}

// Resolve resolves the current database
func (r Resolver) Resolve() error {
	var err error
	for name := range r.db {
		if err = r.resolveNode(name, 0); err != nil {
			return nil
		}
	}
	return nil
}

func (r Resolver) resolveNode(name string, level int) error {
	if level >= r.config.MaxDepth {
		return fmt.Errorf("maximum resolution depth reached")
	}

	node, exists := r.db[name]
	if !exists {
		return nil
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
	return nil
}
