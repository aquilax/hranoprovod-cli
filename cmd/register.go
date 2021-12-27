package cmd

import (
	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/urfave/cli/v2"
)

func newRegisterCommand(ol optionLoader) *cli.Command {
	return &cli.Command{
		Name:    "register",
		Aliases: []string{"reg"},
		Usage:   "Shows the log register report",
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
			&cli.StringFlag{
				Name:    "single-food",
				Aliases: []string{"f"},
				Usage:   "Show only single food",
			},
			&cli.StringFlag{
				Name:    "single-element",
				Aliases: []string{"s"},
				Usage:   "Show only single element",
			},
			&cli.BoolFlag{
				Name:    "group-food",
				Aliases: []string{"g"},
				Usage:   "Single element grouped by food",
			},
			&cli.BoolFlag{
				Name:  "csv",
				Usage: "Export as CSV",
			},
			&cli.BoolFlag{
				Name:  "no-color",
				Usage: "Disable color output",
			},
			&cli.BoolFlag{
				Name:  "no-totals",
				Usage: "Disable totals",
			},
			&cli.BoolFlag{
				Name:  "totals-only",
				Usage: "Show only totals",
			},
			&cli.BoolFlag{
				Name:  "shorten",
				Usage: "Shorten longer strings",
			},
			&cli.BoolFlag{
				Name:  "use-old-reg-reporter",
				Usage: "Use the old reg reporter",
			},
			&cli.StringFlag{
				Name:  "internal-template-name",
				Usage: "Name of the internal demplate to use: [default, left-aligned]",
				Value: "default",
			},
			&cli.BoolFlag{
				Name:  "unresolved",
				Usage: "Deprecated: Show unresolved elements only (moved to 'report unresolved')",
			},
		},
		Action: func(c *cli.Context) error {
			o, err := ol(c)
			if err != nil {
				return err
			}
			return app.NewHranoprovod(o).Register(o.ParserConfig, o.ResolverConfig)
		},
	}
}
