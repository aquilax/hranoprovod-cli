package app

import (
	"io"

	"github.com/aquilax/hranoprovod-cli/v2/filter"
	"github.com/aquilax/hranoprovod-cli/v2/parser"
	"github.com/aquilax/hranoprovod-cli/v2/reporter"
)

type (
	LintCmd  func(stream io.Reader, silent bool, pc parser.Config, rpc reporter.Config) error
	PrintCmd func(logStream io.Reader, dateFormat string, pc parser.Config, rpc reporter.Config, fc filter.Config) error
)

type CSVLogConfig struct {
	ParserConfig   parser.Config
	FilterConfig   filter.Config
	ReporterConfig reporter.CSVConfig
}
