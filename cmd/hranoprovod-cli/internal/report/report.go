package report

import (
	"fmt"
	"io"
	"sort"

	"github.com/aquilax/hranoprovod-cli/v2/cmd/hranoprovod-cli/internal/options"
	"github.com/aquilax/hranoprovod-cli/v2/cmd/hranoprovod-cli/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/v2/cmd/hranoprovod-cli/internal/utils"
	"github.com/aquilax/hranoprovod-cli/v2/lib/filter"
	"github.com/aquilax/hranoprovod-cli/v2/lib/parser"
	"github.com/aquilax/hranoprovod-cli/v2/lib/resolver"
	"github.com/aquilax/hranoprovod-cli/v2/lib/shared"
	"github.com/urfave/cli/v2"
)

type (
	reportElementCmd    func(dbStream io.Reader, rec ReportElementConfig) error
	reportUnresolvedCmd func(logStream, dbStream io.Reader, ruc ReportUnresolvedConfig) error
	reportQuantityCmd   func(logStream io.Reader, rqc ReportQuantityConfig) error
)

func Command() *cli.Command {
	return NewReportCommand(utils.NewCmdUtils())
}

func NewReportCommand(cu utils.CmdUtils) *cli.Command {
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

func newReportElementTotalCommand(cu utils.CmdUtils, reportElement reportElementCmd) *cli.Command {
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
			return cu.WithOptions(c, func(o *options.Options) error {
				return cu.WithFileReaders([]string{o.GlobalConfig.DbFileName}, func(streams []io.Reader) error {
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

func newReportUnresolvedCommand(cu utils.CmdUtils, reportUnresolved reportUnresolvedCmd) *cli.Command {
	return &cli.Command{
		Name:  "unresolved",
		Usage: "Print list of unresolved elements",
		Action: func(c *cli.Context) error {
			return cu.WithOptions(c, func(o *options.Options) error {
				return cu.WithFileReaders([]string{o.GlobalConfig.DbFileName, o.GlobalConfig.LogFileName}, func(streams []io.Reader) error {
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

func newReportQuantityCommand(cu utils.CmdUtils, reportQuantity reportQuantityCmd) *cli.Command {
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
			return cu.WithOptions(c, func(o *options.Options) error {
				return cu.WithFileReaders([]string{o.GlobalConfig.LogFileName}, func(streams []io.Reader) error {
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
	nl, err := utils.LoadDatabaseFromStream(dbStream, rec.ParserConfig)
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
	return NewElementReporter(rec.ReporterConfig, list).Flush()
}

type ReportUnresolvedConfig struct {
	DateFormat     string
	ParserConfig   parser.Config
	ResolverConfig resolver.Config
	ReporterConfig reporter.Config
	FilterConfig   filter.Config
}

// ReportUnresolved generates report for unresolved elements
func ReportUnresolved(logStream, dbStream io.Reader, ruc ReportUnresolvedConfig) error {
	return utils.WithResolvedDatabase(dbStream, ruc.ParserConfig, ruc.ResolverConfig,
		func(nl shared.DBNodeMap) error {
			r := NewUnsolvedReporter(ruc.ReporterConfig, nl)
			f := filter.GetIntervalNodeFilter(ruc.FilterConfig)
			return utils.WalkNodesInStream(logStream, ruc.DateFormat, ruc.ParserConfig, f, r)
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
	r := NewQuantityReporter(rqc.ReporterConfig, rqc.Descending)
	f := filter.GetIntervalNodeFilter(rqc.FilterConfig)
	return utils.WalkNodesInStream(logStream, rqc.DateFormat, rqc.ParserConfig, f, r)
}
