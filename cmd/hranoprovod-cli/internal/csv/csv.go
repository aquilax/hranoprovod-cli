package csv

import (
	"io"
	"sort"

	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/options"
	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/utils"
	shared "github.com/aquilax/hranoprovod-cli/v3"
	"github.com/aquilax/hranoprovod-cli/v3/filter"
	"github.com/aquilax/hranoprovod-cli/v3/parser"
	"github.com/aquilax/hranoprovod-cli/v3/resolver"
	"github.com/urfave/cli/v2"
)

type (
	CSVLogCmd              func(logStream io.Reader, c CSVLogConfig) error
	CSVDatabaseCmd         func(dbStream io.Reader, cdc CSVDatabaseConfig) error
	CSVDatabaseResolvedCmd func(dbStream io.Reader, cdrc CSVDatabaseResolvedConfig) error
)

func Command() *cli.Command {
	return NewCSVCommand(utils.NewCmdUtils())
}

func NewCSVCommand(cu utils.CmdUtils) *cli.Command {
	return &cli.Command{
		Name:  "csv",
		Usage: "Generates csv exports",
		Subcommands: []*cli.Command{
			NewCSVLogCommand(cu, CSVLog),
			NewCSVDatabaseCommand(cu, CSVDatabase),
			NewCSVDatabaseResolvedCommand(cu, CSVDatabaseResolved),
		},
	}
}

func NewCSVLogCommand(cu utils.CmdUtils, csvLog CSVLogCmd) *cli.Command {
	return &cli.Command{
		Name:  "log",
		Usage: "Exports the log file as CSV",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "begin",
				Aliases: []string{"b"},
				Usage:   "Beginning of period `DATE`",
			},
			&cli.StringFlag{
				Name:    "end",
				Aliases: []string{"e"},
				Usage:   "End of period `DATE`",
			},
		},
		Action: func(c *cli.Context) error {
			return cu.WithOptions(c, func(o *options.Options) error {
				cfg := CSVLogConfig{
					DateFormat:     o.GlobalConfig.DateFormat,
					ParserConfig:   o.ParserConfig,
					FilterConfig:   o.FilterConfig,
					ReporterConfig: NewCSVConfig(reporter.NewCommonConfig(o.ReporterConfig.Output, o.ReporterConfig.Color)),
				}
				return cu.WithFileReaders([]string{o.GlobalConfig.LogFileName}, func(streams []io.Reader) error {
					logStream := streams[0]
					return csvLog(logStream, cfg)
				})
			})
		},
	}
}

func NewCSVDatabaseCommand(cu utils.CmdUtils, csvDatabase CSVDatabaseCmd) *cli.Command {
	return &cli.Command{
		Name:  "database",
		Usage: "Exports the database file as CSV",
		Action: func(c *cli.Context) error {
			return cu.WithOptions(c, func(o *options.Options) error {
				return cu.WithFileReaders([]string{o.GlobalConfig.DbFileName}, func(streams []io.Reader) error {
					dbStream := streams[0]
					return csvDatabase(dbStream, CSVDatabaseConfig{
						ParserConfig:   o.ParserConfig,
						ReporterConfig: o.ReporterConfig,
					})
				})
			})
		},
	}
}
func NewCSVDatabaseResolvedCommand(cu utils.CmdUtils, csvDatabaseResolved CSVDatabaseResolvedCmd) *cli.Command {
	return &cli.Command{
		Name:  "database-resolved",
		Usage: "Exports the resolved database as CSV",
		Action: func(c *cli.Context) error {
			return cu.WithOptions(c, func(o *options.Options) error {
				return cu.WithFileReaders([]string{o.GlobalConfig.DbFileName}, func(streams []io.Reader) error {
					dbStream := streams[0]
					return CSVDatabaseResolved(dbStream, CSVDatabaseResolvedConfig{
						ParserConfig:   o.ParserConfig,
						ReporterConfig: o.ReporterConfig,
						ResolverConfig: o.ResolverConfig,
					})
				})
			})
		},
	}
}

type CSVLogConfig struct {
	DateFormat     string
	ParserConfig   parser.Config
	FilterConfig   filter.Config
	ReporterConfig CSVConfig
}

// CSVLog generates CSV export of the log
func CSVLog(logStream io.Reader, c CSVLogConfig) error {
	r := NewCSVReporter(c.ReporterConfig)
	defer r.Flush()
	f := filter.GetIntervalNodeFilter(c.FilterConfig)
	return utils.WalkNodesInStream(logStream, c.DateFormat, c.ParserConfig, f, r)
}

type CSVDatabaseConfig struct {
	ParserConfig   parser.Config
	ReporterConfig reporter.Config
}

// CSVDatabase generates CSV export of the database
func CSVDatabase(dbStream io.Reader, cdc CSVDatabaseConfig) error {
	r := NewCSVDatabaseReporter(cdc.ReporterConfig)
	defer r.Flush()

	return parser.ParseStreamCallback(dbStream, cdc.ParserConfig, func(n *shared.ParserNode, err error) (stop bool, cbError error) {
		if err := r.Process(shared.NewDBNodeFromNode(n)); err != nil {
			return true, err
		}
		return false, nil
	})
}

type CSVDatabaseResolvedConfig struct {
	ParserConfig   parser.Config
	ReporterConfig reporter.Config
	ResolverConfig resolver.Config
}

// CSVDatabaseResolved generates CSV export of the resolved database
func CSVDatabaseResolved(dbStream io.Reader, cdc CSVDatabaseResolvedConfig) error {
	nl, err := utils.LoadDatabaseFromStream(dbStream, cdc.ParserConfig)
	if err != nil {
		return err
	}
	nl, err = resolver.Resolve(cdc.ResolverConfig, nl)
	if err != nil {
		return err
	}
	keys := make([]string, len(nl))
	i := 0
	for n := range nl {
		keys[i] = n
		i++
	}
	sort.Strings(keys)
	r := NewCSVDatabaseReporter(cdc.ReporterConfig)

	for _, key := range keys {
		if err = r.Process(nl[key]); err != nil {
			return err
		}
	}
	return r.Flush()
}
