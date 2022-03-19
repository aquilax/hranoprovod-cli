package cmd

import (
	"fmt"
	"io"

	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/aquilax/hranoprovod-cli/v2/options"
	"github.com/urfave/cli/v2"
)

type (
	reportElementCmd    func(dbStream io.Reader, rec app.ReportElementConfig) error
	reportUnresolvedCmd func(logStream, dbStream io.Reader, ruc app.ReportUnresolvedConfig) error
	reportQuantityCmd   func(logStream io.Reader, rqc app.ReportQuantityConfig) error
)

func newReportCommand(cu cmdUtils) *cli.Command {
	return &cli.Command{
		Name:  "report",
		Usage: "Generates various reports",
		Subcommands: []*cli.Command{
			newReportElementTotalCommand(cu, app.ReportElement),
			newReportUnresolvedCommand(cu, app.ReportUnresolved),
			newReportQuantityCommand(cu, app.ReportQuantity),
		},
	}
}

func newReportElementTotalCommand(cu cmdUtils, reportElement reportElementCmd) *cli.Command {
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
			return cu.withOptions(c, func(o *options.Options) error {
				return cu.withFileReaders([]string{o.GlobalConfig.DbFileName}, func(streams []io.Reader) error {
					dbStream := streams[0]
					return reportElement(dbStream, app.ReportElementConfig{
						ElementName:    c.Args().First(),
						Descending:     c.IsSet("desc"),
						ParserConfig:   o.ParserConfig,
						ResolverConfig: o.ResolverConfig,
						ReporterConfig: o.ReporterConfig,
					})
				})
			})
		},
	}
}

func newReportUnresolvedCommand(cu cmdUtils, reportUnresolved reportUnresolvedCmd) *cli.Command {
	return &cli.Command{
		Name:  "unresolved",
		Usage: "Print list of unresolved elements",
		Action: func(c *cli.Context) error {
			return cu.withOptions(c, func(o *options.Options) error {
				return cu.withFileReaders([]string{o.GlobalConfig.DbFileName, o.GlobalConfig.LogFileName}, func(streams []io.Reader) error {
					dbStream, logStream := streams[0], streams[1]
					return reportUnresolved(logStream, dbStream, app.ReportUnresolvedConfig{
						DateFormat:     o.GlobalConfig.DateFormat,
						ParserConfig:   o.ParserConfig,
						ResolverConfig: o.ResolverConfig,
						ReporterConfig: o.ReporterConfig,
						FilterConfig:   o.FilterConfig,
					})
				})
			})
		},
	}
}

func newReportQuantityCommand(cu cmdUtils, reportQuantity reportQuantityCmd) *cli.Command {
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
			return cu.withOptions(c, func(o *options.Options) error {
				return cu.withFileReaders([]string{o.GlobalConfig.LogFileName}, func(streams []io.Reader) error {
					logStream := streams[0]
					return reportQuantity(logStream, app.ReportQuantityConfig{
						DateFormat:     o.GlobalConfig.DateFormat,
						Descending:     c.IsSet("desc"),
						ParserConfig:   o.ParserConfig,
						ReporterConfig: o.ReporterConfig,
						FilterConfig:   o.FilterConfig,
					})
				})
			})
		},
	}
}
