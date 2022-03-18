package cmd

import (
	"io"

	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/aquilax/hranoprovod-cli/v2/reporter"
	"github.com/urfave/cli/v2"
)

func newCSVCommand(cu cmdUtils) *cli.Command {
	return &cli.Command{
		Name:  "csv",
		Usage: "Generates csv exports",
		Subcommands: []*cli.Command{
			newCSVLogCommand(cu),
			newCSVDatabaseCommand(cu),
			newCSVDatabaseResolvedCommand(cu),
		},
	}
}

func newCSVLogCommand(cu cmdUtils) *cli.Command {
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
			return cu.withOptions(c, func(o *app.Options) error {
				cfg := app.CSVLogConfig{
					ParserConfig:   o.ParserConfig,
					FilterConfig:   o.FilterConfig,
					ReporterConfig: reporter.NewCSVConfig(reporter.NewCommonConfig(o.ReporterConfig.Color)),
				}
				return cu.withFileReaders([]string{o.GlobalConfig.LogFileName}, func(streams []io.Reader) error {
					logStream := streams[0]
					return app.CSVLog(logStream, cfg)
				})
			})
		},
	}
}

func newCSVDatabaseCommand(cu cmdUtils) *cli.Command {
	return &cli.Command{
		Name:  "database",
		Usage: "Exports the database file as CSV",
		Action: func(c *cli.Context) error {
			return cu.withOptions(c, func(o *app.Options) error {
				return cu.withFileReaders([]string{o.GlobalConfig.DbFileName}, func(streams []io.Reader) error {
					dbStream := streams[0]
					return app.CSVDatabase(dbStream, o.ParserConfig, o.ReporterConfig)
				})
			})
		},
	}
}
func newCSVDatabaseResolvedCommand(cu cmdUtils) *cli.Command {
	return &cli.Command{
		Name:  "database-resolved",
		Usage: "Exports the resolved database as CSV",
		Action: func(c *cli.Context) error {
			return cu.withOptions(c, func(o *app.Options) error {
				return cu.withFileReaders([]string{o.GlobalConfig.DbFileName}, func(streams []io.Reader) error {
					dbStream := streams[0]
					return app.CSVDatabaseResolved(dbStream, o.ParserConfig, o.ReporterConfig, o.ResolverConfig)
				})
			})
		},
	}
}
