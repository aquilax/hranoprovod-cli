package resolver

import (
	"fmt"

	shared "github.com/aquilax/hranoprovod-cli/v3"
)

const DefaultMaxDepth = 10

type Config struct {
	MaxDepth int
}

func NewDefaultConfig() Config {
	return Config{DefaultMaxDepth}
}

// Resolver contains the resolver data
// Deprecated: Deprecated in favor of using Resolve function directly
type Resolver struct {
	db     shared.DBNodeMap
	config Config
}

// NewResolver creates new resolver
// Deprecated: Deprecated in favor of using Resolve function directly
func NewResolver(db shared.DBNodeMap, c Config) Resolver {
	return Resolver{db, c}
}

// Resolve resolves the current database
// Deprecated: Deprecated in favor of using Resolve function directly
func (r Resolver) Resolve() error {
	var err error
	for name := range r.db {
		if err = r.resolveNode(name, 0); err != nil {
			return err
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
		if err := r.resolveNode(e.Name, level+1); err != nil {
			return err
		}
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

func resolveNode(maxDepth int, db shared.DBNodeMap, name string, level int) error {
	if level >= maxDepth {
		return fmt.Errorf("maximum resolution depth reached")
	}

	node, exists := db[name]
	if !exists {
		return nil
	}

	nel := shared.NewElements()

	for _, e := range node.Elements {
		if err := resolveNode(maxDepth, db, e.Name, level+1); err != nil {
			return err
		}
		if foundNode, exists := db[e.Name]; exists {
			nel.SumMerge(foundNode.Elements, e.Value)
		} else {
			var tm shared.Elements
			tm.Add(e.Name, e.Value)
			nel.SumMerge(tm, 1)
		}
	}
	nel.Sort()
	db[name].Elements = nel
	return nil
}

func Resolve(c Config, db shared.DBNodeMap) (shared.DBNodeMap, error) {
	for name := range db {
		if err := resolveNode(c.MaxDepth, db, name, 0); err != nil {
			return db, err
		}
	}
	return db, nil
}
