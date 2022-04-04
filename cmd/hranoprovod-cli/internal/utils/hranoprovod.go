package utils

import (
	"io"
	"time"

	"github.com/aquilax/hranoprovod-cli/v2/cmd/hranoprovod-cli/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/v2/lib/filter"
	"github.com/aquilax/hranoprovod-cli/v2/lib/parser"
	"github.com/aquilax/hranoprovod-cli/v2/lib/resolver"
	"github.com/aquilax/hranoprovod-cli/v2/lib/shared"
)

type ResolvedCallback = func(nl shared.DBNodeMap) error
type ReporterCallback func(rpc reporter.Config, nl shared.DBNodeMap) reporter.Reporter

func WithResolvedDatabase(dbStream io.Reader, pc parser.Config, rc resolver.Config, cb ResolvedCallback) error {
	if nl, err := LoadDatabaseFromStream(dbStream, pc); err == nil {
		if nl, err = resolver.Resolve(rc, nl); err == nil {
			return cb(nl)
		} else {
			return err
		}
	} else {
		return err
	}
}

func WalkWithReporter(logStream, dbStream io.Reader, dateFormat string, pc parser.Config, rc resolver.Config, rpc reporter.Config, fc filter.Config, rpCb ReporterCallback) error {
	return WithResolvedDatabase(dbStream, pc, rc,
		func(nl shared.DBNodeMap) error {
			r := rpCb(rpc, nl)
			f := filter.GetIntervalNodeFilter(fc)
			return WalkNodesInStream(logStream, dateFormat, pc, f, r)
		})
}

func LoadDatabaseFromStream(dbStream io.Reader, pc parser.Config) (shared.DBNodeMap, error) {
	nodeMap := shared.NewDBNodeMap()
	return nodeMap, parser.ParseStreamCallback(dbStream, pc, func(node *shared.ParserNode, err error) (stop bool, cbError error) {
		if err != nil {
			return true, err
		} else {
			nodeMap.Push(shared.NewDBNodeFromNode(node))
			return false, nil
		}
	})
}

func WalkNodesInStream(logStream io.Reader, dateFormat string, pc parser.Config, filter *filter.LogNodeFilter, r reporter.Reporter) error {
	var ln *shared.LogNode
	var t time.Time
	var ok bool

	cb := func(node *shared.ParserNode, err error) (stop bool, cbError error) {
		if err != nil {
			return true, err
		}
		if t, err = time.Parse(dateFormat, node.Header); err != nil {
			return true, err
		}
		ok = true
		if filter != nil {
			if ok, err = (*filter)(t, node); err != nil {
				return true, err
			}
		}
		if ok {
			if ln, err = shared.NewLogNodeFromElements(t, node.Elements, node.Metadata); err != nil {
				return true, err
			}
			if err = r.Process(ln); err != nil {
				return true, err
			}
		}
		return false, nil
	}
	err := parser.ParseStreamCallback(logStream, pc, cb)
	if err != nil {
		return err
	}
	return r.Flush()
}
