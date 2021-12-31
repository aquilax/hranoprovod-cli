package cmd

import (
	"io"

	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/urfave/cli/v2"
)

func newBalanceCommand(ol optionLoader) *cli.Command {
	return &cli.Command{
		Name:    "balance",
		Aliases: []string{"bal"},
		Usage:   "Shows food balance as tree",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "begin",
				Aliases: []string{"b"},
				Usage:   "Beginning of period",
			},
			&cli.StringFlag{
				Name:    "end",
				Aliases: []string{"e"},
				Usage:   "End of period",
			},
			&cli.BoolFlag{
				Name:  "collapse-last",
				Usage: "Collapses last dimension",
			},
			&cli.BoolFlag{
				Name:    "collapse",
				Aliases: []string{"c"},
				Usage:   "Collapses sole branches",
			},
			&cli.StringFlag{
				Name:    "single-element, s",
				Aliases: []string{"s"},
				Usage:   "Show only single element",
			},
		},
		Action: func(c *cli.Context) error {
			if o, err := ol(c); err != nil {
				return err
			} else {
				return withFileReader(o.GlobalConfig.DbFileName, func(dbStream io.Reader) error {
					return withFileReader(o.GlobalConfig.LogFileName, func(logStream io.Reader) error {
						return app.Balance(logStream, dbStream, o.GlobalConfig.DateFormat, o.ParserConfig, o.ResolverConfig, o.ReporterConfig, o.FilterConfig)
					})
				})
			}
		},
	}
}
