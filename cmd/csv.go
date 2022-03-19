package cmd

import (
	"io"

	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/aquilax/hranoprovod-cli/v2/options"
	"github.com/aquilax/hranoprovod-cli/v2/reporter"
	"github.com/urfave/cli/v2"
)

type (
	csvLogCmd              func(logStream io.Reader, c app.CSVLogConfig) error
	csvDatabaseCmd         func(dbStream io.Reader, cdc app.CSVDatabaseConfig) error
	csvDatabaseResolvedCmd func(dbStream io.Reader, cdrc app.CSVDatabaseResolvedConfig) error
)

func newCSVCommand(cu cmdUtils) *cli.Command {
	return &cli.Command{
		Name:  "csv",
		Usage: "Generates csv exports",
		Subcommands: []*cli.Command{
			newCSVLogCommand(cu, app.CSVLog),
			newCSVDatabaseCommand(cu, app.CSVDatabase),
			newCSVDatabaseResolvedCommand(cu, app.CSVDatabaseResolved),
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
			return cu.withOptions(c, func(o *options.Options) error {
				cfg := app.CSVLogConfig{
					ParserConfig:   o.ParserConfig,
					FilterConfig:   o.FilterConfig,
					ReporterConfig: reporter.NewCSVConfig(reporter.NewCommonConfig(o.ReporterConfig.Color)),
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
			return cu.withOptions(c, func(o *options.Options) error {
				return cu.withFileReaders([]string{o.GlobalConfig.DbFileName}, func(streams []io.Reader) error {
					dbStream := streams[0]
					return csvDatabase(dbStream, app.CSVDatabaseConfig{
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
			return cu.withOptions(c, func(o *options.Options) error {
				return cu.withFileReaders([]string{o.GlobalConfig.DbFileName}, func(streams []io.Reader) error {
					dbStream := streams[0]
					return app.CSVDatabaseResolved(dbStream, app.CSVDatabaseResolvedConfig{
						ParserConfig:   o.ParserConfig,
						ReporterConfig: o.ReporterConfig,
						ResolverConfig: o.ResolverConfig,
					})
				})
			})
		},
	}
}
