package cmd

import (
	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/urfave/cli/v2"
)

func newReportCommand(ol optionLoader) *cli.Command {
	return &cli.Command{
		Name:  "report",
		Usage: "Generates various reports",
		Subcommands: []*cli.Command{
			newReportElementTotalCommand(ol),
			newReportUnresolvedCommand(ol),
			newReportQuantityCommand(ol),
		},
	}
}

func newReportElementTotalCommand(ol optionLoader) *cli.Command {
	return &cli.Command{
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
			o, err := ol(c)
			if err != nil {
				return err
			}
			return app.NewHranoprovod(o).ReportElement(c.Args().First(), c.IsSet("desc"), o.ParserConfig, o.ResolverConfig)
		},
	}
}

func newReportUnresolvedCommand(ol optionLoader) *cli.Command {
	return &cli.Command{
		Name:  "unresolved",
		Usage: "Print list of unresolved elements",
		Action: func(c *cli.Context) error {
			o, err := ol(c)
			if err != nil {
				return err
			}
			return app.NewHranoprovod(o).ReportUnresolved(o.ParserConfig, o.ResolverConfig)
		},
	}
}

func newReportQuantityCommand(ol optionLoader) *cli.Command {
	return &cli.Command{
		Name:  "quantity",
		Usage: "Total quantities per food",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "desc",
				Usage: "Descending order",
			},
		},
		Action: func(c *cli.Context) error {
			o, err := ol(c)
			if err != nil {
				return err
			}
			return app.NewHranoprovod(o).ReportQuantity(c.IsSet("desc"), o.ParserConfig)
		},
	}
}
