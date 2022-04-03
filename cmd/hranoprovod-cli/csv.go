package main

import (
	"io"
	"sort"

	"github.com/aquilax/hranoprovod-cli/v2/lib/filter"
	"github.com/aquilax/hranoprovod-cli/v2/lib/parser"
	"github.com/aquilax/hranoprovod-cli/v2/lib/reporter"
	"github.com/aquilax/hranoprovod-cli/v2/lib/resolver"
	"github.com/aquilax/hranoprovod-cli/v2/lib/shared"
	"github.com/urfave/cli/v2"
)

type (
	csvLogCmd              func(logStream io.Reader, c CSVLogConfig) error
	csvDatabaseCmd         func(dbStream io.Reader, cdc CSVDatabaseConfig) error
	csvDatabaseResolvedCmd func(dbStream io.Reader, cdrc CSVDatabaseResolvedConfig) error
)

func newCSVCommand(cu cmdUtils) *cli.Command {
	return &cli.Command{
		Name:  "csv",
		Usage: "Generates csv exports",
		Subcommands: []*cli.Command{
			newCSVLogCommand(cu, CSVLog),
			newCSVDatabaseCommand(cu, CSVDatabase),
			newCSVDatabaseResolvedCommand(cu, CSVDatabaseResolved),
		},
	}
}

func newCSVLogCommand(cu cmdUtils, csvLog csvLogCmd) *cli.Command {
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
			return cu.withOptions(c, func(o *Options) error {
				cfg := CSVLogConfig{
					DateFormat:     o.GlobalConfig.DateFormat,
					ParserConfig:   o.ParserConfig,
					FilterConfig:   o.FilterConfig,
					ReporterConfig: reporter.NewCSVConfig(reporter.NewCommonConfig(o.ReporterConfig.Output, o.ReporterConfig.Color)),
				}
				return cu.withFileReaders([]string{o.GlobalConfig.LogFileName}, func(streams []io.Reader) error {
					logStream := streams[0]
					return csvLog(logStream, cfg)
				})
			})
		},
	}
}

func newCSVDatabaseCommand(cu cmdUtils, csvDatabase csvDatabaseCmd) *cli.Command {
	return &cli.Command{
		Name:  "database",
		Usage: "Exports the database file as CSV",
		Action: func(c *cli.Context) error {
			return cu.withOptions(c, func(o *Options) error {
				return cu.withFileReaders([]string{o.GlobalConfig.DbFileName}, func(streams []io.Reader) error {
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
func newCSVDatabaseResolvedCommand(cu cmdUtils, csvDatabaseResolved csvDatabaseResolvedCmd) *cli.Command {
	return &cli.Command{
		Name:  "database-resolved",
		Usage: "Exports the resolved database as CSV",
		Action: func(c *cli.Context) error {
			return cu.withOptions(c, func(o *Options) error {
				return cu.withFileReaders([]string{o.GlobalConfig.DbFileName}, func(streams []io.Reader) error {
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
	ReporterConfig reporter.CSVConfig
}

// CSVLog generates CSV export of the log
func CSVLog(logStream io.Reader, c CSVLogConfig) error {
	r := reporter.NewCSVReporter(c.ReporterConfig)
	f := filter.GetIntervalNodeFilter(c.FilterConfig)
	return walkNodesInStream(logStream, c.DateFormat, c.ParserConfig, f, r)
}

type CSVDatabaseConfig struct {
	ParserConfig   parser.Config
	ReporterConfig reporter.Config
}

// CSVDatabase generates CSV export of the database
func CSVDatabase(dbStream io.Reader, cdc CSVDatabaseConfig) error {
	p := parser.NewParser(cdc.ParserConfig)
	r := reporter.NewCSVDatabaseReporter(cdc.ReporterConfig)
	go p.ParseStream(dbStream)
	return func() error {
		for {
			select {
			case node := <-p.Nodes:
				r.Process(shared.NewDBNodeFromNode(node))
			case error := <-p.Errors:
				return error
			case <-p.Done:
				return r.Flush()
			}
		}
	}()
}

type CSVDatabaseResolvedConfig struct {
	ParserConfig   parser.Config
	ReporterConfig reporter.Config
	ResolverConfig resolver.Config
}

// CSVDatabaseResolved generates CSV export of the resolved database
func CSVDatabaseResolved(dbStream io.Reader, cdc CSVDatabaseResolvedConfig) error {
	nl, err := loadDatabaseFromStream(dbStream, cdc.ParserConfig)
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
	r := reporter.NewCSVDatabaseReporter(cdc.ReporterConfig)
	for _, key := range keys {
		if err = r.Process(nl[key]); err != nil {
			return err
		}
	}
	return r.Flush()
}
