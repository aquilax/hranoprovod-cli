package utils

import (
	"io"
	"time"

	"github.com/aquilax/hranoprovod-cli/v3/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/filter"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/node"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/parser"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/resolver"
)

type ResolvedCallback = func(nl node.DBNodeMap) error
type ReporterCallback func(rpc reporter.Config, nl node.DBNodeMap) reporter.Reporter

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
		func(nl node.DBNodeMap) error {
			r := rpCb(rpc, nl)
			f := filter.GetIntervalNodeFilter(fc)
			if err := WalkNodesInStream(logStream, dateFormat, pc, f, r); err != nil {
				return err
			}
			return r.Flush()
		})
}

func LoadDatabaseFromStream(dbStream io.Reader, pc parser.Config) (node.DBNodeMap, error) {
	nodeMap := node.NewDBNodeMap()
	return nodeMap, parser.ParseStreamCallback(dbStream, pc, func(pn *node.ParserNode, err error) (stop bool, cbError error) {
		if err != nil {
			return true, err
		} else {
			nodeMap.Push(node.NewDBNodeFromNode(pn))
			return false, nil
		}
	})
}

func WalkNodesInStream(logStream io.Reader, dateFormat string, pc parser.Config, filter *filter.LogNodeFilter, r reporter.Reporter) error {
	var ln *node.LogNode
	var t time.Time
	var ok bool

	cb := func(pn *node.ParserNode, err error) (stop bool, cbError error) {
		if err != nil {
			return true, err
		}
		if t, err = time.Parse(dateFormat, pn.Header); err != nil {
			return true, err
		}
		ok = true
		if filter != nil {
			if ok, err = (*filter)(t, pn); err != nil {
				return true, err
			}
		}
		if ok {
			if ln, err = node.NewLogNodeFromElements(t, pn.Elements, pn.Metadata); err != nil {
				return true, err
			}
			if err = r.Process(ln); err != nil {
				return true, err
			}
		}
		return false, nil
	}
	return parser.ParseStreamCallback(logStream, pc, cb)
}
