package cmd

import (
	"fmt"
	"io"

	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/urfave/cli/v2"
)

func newReportCommand(cu cmdUtils) *cli.Command {
	return &cli.Command{
		Name:  "report",
		Usage: "Generates various reports",
		Subcommands: []*cli.Command{
			newReportElementTotalCommand(cu),
			newReportUnresolvedCommand(cu),
			newReportQuantityCommand(cu),
		},
	}
}

func newReportElementTotalCommand(cu cmdUtils) *cli.Command {
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
			return cu.withOptions(c, func(o *app.Options) error {
				return cu.withFileReaders([]string{o.GlobalConfig.DbFileName}, func(streams []io.Reader) error {
					dbStream := streams[0]
					return app.ReportElement(dbStream, c.Args().First(), c.IsSet("desc"), o.ParserConfig, o.ResolverConfig, o.ReporterConfig)
				})
			})
		},
	}
}

func newReportUnresolvedCommand(cu cmdUtils) *cli.Command {
	return &cli.Command{
		Name:  "unresolved",
		Usage: "Print list of unresolved elements",
		Action: func(c *cli.Context) error {
			return cu.withOptions(c, func(o *app.Options) error {
				return cu.withFileReaders([]string{o.GlobalConfig.DbFileName, o.GlobalConfig.LogFileName}, func(streams []io.Reader) error {
					dbStream, logStream := streams[0], streams[1]
					return app.ReportUnresolved(logStream, dbStream, o.GlobalConfig.DateFormat, o.ParserConfig, o.ResolverConfig, o.ReporterConfig, o.FilterConfig)
				})
			})
		},
	}
}

func newReportQuantityCommand(cu cmdUtils) *cli.Command {
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
			return cu.withOptions(c, func(o *app.Options) error {
				return cu.withFileReaders([]string{o.GlobalConfig.LogFileName}, func(streams []io.Reader) error {
					logStream := streams[0]
					return app.ReportQuantity(logStream, o.GlobalConfig.DateFormat, c.IsSet("desc"), o.ParserConfig, o.ReporterConfig, o.FilterConfig)
				})
			})
		},
	}
}
