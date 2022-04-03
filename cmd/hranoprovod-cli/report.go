package main

import (
	"fmt"
	"io"
	"sort"

	"github.com/aquilax/hranoprovod-cli/v2/lib/filter"
	"github.com/aquilax/hranoprovod-cli/v2/lib/parser"
	"github.com/aquilax/hranoprovod-cli/v2/lib/reporter"
	"github.com/aquilax/hranoprovod-cli/v2/lib/resolver"
	"github.com/aquilax/hranoprovod-cli/v2/lib/shared"
	"github.com/urfave/cli/v2"
)

type (
	reportElementCmd    func(dbStream io.Reader, rec ReportElementConfig) error
	reportUnresolvedCmd func(logStream, dbStream io.Reader, ruc ReportUnresolvedConfig) error
	reportQuantityCmd   func(logStream io.Reader, rqc ReportQuantityConfig) error
)

func newReportCommand(cu cmdUtils) *cli.Command {
	return &cli.Command{
		Name:  "report",
		Usage: "Generates various reports",
		Subcommands: []*cli.Command{
			newReportElementTotalCommand(cu, ReportElement),
			newReportUnresolvedCommand(cu, ReportUnresolved),
			newReportQuantityCommand(cu, ReportQuantity),
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
			return cu.withOptions(c, func(o *Options) error {
				return cu.withFileReaders([]string{o.GlobalConfig.DbFileName}, func(streams []io.Reader) error {
					dbStream := streams[0]
					return reportElement(dbStream, ReportElementConfig{
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
			return cu.withOptions(c, func(o *Options) error {
				return cu.withFileReaders([]string{o.GlobalConfig.DbFileName, o.GlobalConfig.LogFileName}, func(streams []io.Reader) error {
					dbStream, logStream := streams[0], streams[1]
					return reportUnresolved(logStream, dbStream, ReportUnresolvedConfig{
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
			return cu.withOptions(c, func(o *Options) error {
				return cu.withFileReaders([]string{o.GlobalConfig.LogFileName}, func(streams []io.Reader) error {
					logStream := streams[0]
					return reportQuantity(logStream, ReportQuantityConfig{
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

type ReportElementConfig struct {
	ElementName    string
	Descending     bool
	ParserConfig   parser.Config
	ResolverConfig resolver.Config
	ReporterConfig reporter.Config
}

// ReportElement generates report for single element
func ReportElement(dbStream io.Reader, rec ReportElementConfig) error {
	nl, err := loadDatabaseFromStream(dbStream, rec.ParserConfig)
	if err != nil {
		return err
	}
	nl, err = resolver.Resolve(rec.ResolverConfig, nl)
	if err != nil {
		return err
	}
	var list []shared.Element
	for name, node := range nl {
		for _, el := range node.Elements {
			if el.Name == rec.ElementName {
				list = append(list, shared.NewElement(name, el.Value))
			}
		}
	}
	if rec.Descending {
		sort.SliceStable(list, func(i, j int) bool {
			return list[i].Value > list[j].Value
		})
	} else {
		sort.SliceStable(list, func(i, j int) bool {
			return list[i].Value < list[j].Value
		})
	}
	return reporter.NewElementReporter(rec.ReporterConfig, list).Flush()
}

type ReportUnresolvedConfig = RegisterConfig

// ReportUnresolved generates report for unresolved elements
func ReportUnresolved(logStream, dbStream io.Reader, ruc ReportUnresolvedConfig) error {
	return withResolvedDatabase(dbStream, ruc.ParserConfig, ruc.ResolverConfig,
		func(nl shared.DBNodeMap) error {
			r := reporter.NewUnsolvedReporter(ruc.ReporterConfig, nl)
			f := filter.GetIntervalNodeFilter(ruc.FilterConfig)
			return walkNodesInStream(logStream, ruc.DateFormat, ruc.ParserConfig, f, r)
		})
}

type ReportQuantityConfig struct {
	DateFormat     string
	Descending     bool
	ParserConfig   parser.Config
	ReporterConfig reporter.Config
	FilterConfig   filter.Config
}

// ReportQuantity Generates a quantity report
func ReportQuantity(logStream io.Reader, rqc ReportQuantityConfig) error {
	r := reporter.NewQuantityReporter(rqc.ReporterConfig, rqc.Descending)
	f := filter.GetIntervalNodeFilter(rqc.FilterConfig)
	return walkNodesInStream(logStream, rqc.DateFormat, rqc.ParserConfig, f, r)
}
