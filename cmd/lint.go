package cmd

import (
	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/urfave/cli/v2"
)

func newLintCommand() *cli.Command {
	return &cli.Command{
		Name:  "lint",
		Usage: "Lints file",
		Action: func(c *cli.Context) error {
			o := app.NewOptions()
			if err := o.Load(c); err != nil {
				return err
			}
			return app.NewHranoprovod(o).Lint(c.Args().First())
		},
	}
}
