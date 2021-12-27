package cmd

import (
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
			o, err := ol(c)
			if err != nil {
				return err
			}
			return app.NewHranoprovod(o).Balance(o.ParserConfig)
		},
	}
}
