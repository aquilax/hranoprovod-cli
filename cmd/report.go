package cmd

import (
	"fmt"
	"io"

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
		ArgsUsage: "[element-name]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "desc",
				Usage: "Descending order",
			},
		},
		Before: func(c *cli.Context) error {
			if c.Args().First() == "" {
				return fmt.Errorf("no element name")
			}
			return nil
		},
		Action: func(c *cli.Context) error {
			if o, err := ol(c); err != nil {
				return err
			} else {
				return withFileReader(o.GlobalConfig.DbFileName, func(dbStream io.Reader) error {
					return app.ReportElement(dbStream, c.Args().First(), c.IsSet("desc"), o.ParserConfig, o.ResolverConfig, o.ReporterConfig)
				})
			}
		},
	}
}

func newReportUnresolvedCommand(ol optionLoader) *cli.Command {
	return &cli.Command{
		Name:  "unresolved",
		Usage: "Print list of unresolved elements",
		Action: func(c *cli.Context) error {
			if o, err := ol(c); err != nil {
				return err
			} else {
				return withFileReader(o.GlobalConfig.DbFileName, func(dbStream io.Reader) error {
					return withFileReader(o.GlobalConfig.LogFileName, func(logStream io.Reader) error {
						return app.ReportUnresolved(logStream, dbStream, o.GlobalConfig.DateFormat, o.ParserConfig, o.ResolverConfig, o.ReporterConfig, o.FilterConfig)
					})
				})
			}
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
			return app.ReportQuantity(o.GlobalConfig, c.IsSet("desc"), o.ParserConfig, o.ReporterConfig, o.FilterConfig)
		},
	}
}
