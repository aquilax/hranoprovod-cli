package cmd

import (
	"io"

	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/urfave/cli/v2"
)

func newCsvCommand(ol optionLoader) *cli.Command {
	return &cli.Command{
		Name:  "csv",
		Usage: "Generates csv exports",
		Subcommands: []*cli.Command{
			{
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
						return withFileReader(o.GlobalConfig.LogFileName, func(logStream io.Reader) error {
							return app.CSVLog(logStream, o.GlobalConfig.DateFormat, o.ParserConfig, o.ReporterConfig, o.FilterConfig)
						})
					}
				},
			},
			{
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
			},
			{
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
			},
		},
	}
}
