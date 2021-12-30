package cmd

import (
	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/urfave/cli/v2"
)

func newLintCommand(ol optionLoader) *cli.Command {
	return &cli.Command{
		Name:  "lint",
		Usage: "Lints file",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "silent",
				Aliases: []string{"s"},
				Usage:   "stay silent if no errors are found",
			},
		},
		Action: func(c *cli.Context) error {
			o, err := ol(c)
			if err != nil {
				return err
			}
			o.ParserConfig.StopOnError = false
			return app.Lint(c.Args().First(), c.IsSet("silent"), o.ParserConfig, o.ReporterConfig)
		},
	}
}
