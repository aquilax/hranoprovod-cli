package cmd

import (
	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/urfave/cli/v2"
)

func newLintCommand(ol optionLoader) *cli.Command {
	return &cli.Command{
		Name:  "lint",
		Usage: "Lints file",
		Action: func(c *cli.Context) error {
			o, err := ol(c)
			if err != nil {
				return err
			}
			o.ParserConfig.StopOnError = false
			return app.Lint(c.Args().First(), o.ParserConfig)
		},
	}
}
