package cmd

import (
	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/urfave/cli/v2"
)

func newReportCommand() *cli.Command {
	return &cli.Command{
		Name:  "report",
		Usage: "Generates various reports",
		Subcommands: []*cli.Command{
			{
				Name:      "element-total",
				Usage:     "Generates total sum for element grouped by food",
				ArgsUsage: "[element name]",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "desc",
						Usage: "Descending order",
					},
				},
				Action: func(c *cli.Context) error {
					o := app.NewOptions()
					if err := o.Load(c); err != nil {
						return err
					}
					return app.NewHranoprovod(o).ReportElement(c.Args().First(), c.IsSet("desc"))
				},
			},
			{
				Name:  "unresolved",
				Usage: "Print list of unresolved elements",
				Action: func(c *cli.Context) error {
					o := app.NewOptions()
					if err := o.Load(c); err != nil {
						return err
					}
					return app.NewHranoprovod(o).ReportUnresolved()
				},
			},
			{
				Name:  "quantity",
				Usage: "Total quantities per food",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "desc",
						Usage: "Descending order",
					},
				},
				Action: func(c *cli.Context) error {
					o := app.NewOptions()
					if err := o.Load(c); err != nil {
						return err
					}
					return app.NewHranoprovod(o).ReportQuantity(c.IsSet("desc"))
				},
			},
		},
	}
}
