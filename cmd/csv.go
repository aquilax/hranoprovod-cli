package cmd

import (
	"io"

	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/aquilax/hranoprovod-cli/v2/reporter"
	"github.com/urfave/cli/v2"
)

func newCSVCommand(ol optionLoader) *cli.Command {
	return &cli.Command{
		Name:  "csv",
		Usage: "Generates csv exports",
		Subcommands: []*cli.Command{
			newCSVLogCommand(ol),
			newCSVDatabaseCommand(ol),
			newCSVDatabaseResolvedCommand(ol),
		},
	}
}

func newCSVLogCommand(ol optionLoader) *cli.Command {
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
			if o, err := ol(c); err != nil {
				return err
			} else {
				// TODO: better config loading strategy is needed
				cfg := app.CSVLogConfig{
					ParserConfig:   o.ParserConfig,
					FilterConfig:   o.FilterConfig,
					ReporterConfig: reporter.NewCSVConfig(reporter.CommonConfig{Color: o.ReporterConfig.Color}),
				}
				return withFileReader(o.GlobalConfig.LogFileName, func(logStream io.Reader) error {
					return app.CSVLog(logStream, cfg)
				})
			}
		},
	}
}

func newCSVDatabaseCommand(ol optionLoader) *cli.Command {
	return &cli.Command{
		Name:  "database",
		Usage: "Exports the database file as CSV",
		Action: func(c *cli.Context) error {
			if o, err := ol(c); err != nil {
				return err
			} else {
				return withFileReader(o.GlobalConfig.DbFileName, func(dbStream io.Reader) error {
					return app.CSVDatabase(dbStream, o.ParserConfig, o.ReporterConfig)
				})
			}
		},
	}
}
func newCSVDatabaseResolvedCommand(ol optionLoader) *cli.Command {
	return &cli.Command{
		Name:  "database-resolved",
		Usage: "Exports the resolved database as CSV",
		Action: func(c *cli.Context) error {
			if o, err := ol(c); err != nil {
				return err
			} else {
				return withFileReader(o.GlobalConfig.DbFileName, func(dbStream io.Reader) error {
					return app.CSVDatabaseResolved(dbStream, o.ParserConfig, o.ReporterConfig, o.ResolverConfig)
				})
			}
		},
	}
}
