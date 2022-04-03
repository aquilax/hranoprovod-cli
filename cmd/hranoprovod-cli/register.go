package main

import (
	"io"

	"github.com/aquilax/hranoprovod-cli/v2/lib/filter"
	"github.com/aquilax/hranoprovod-cli/v2/lib/parser"
	"github.com/aquilax/hranoprovod-cli/v2/lib/reporter"
	"github.com/aquilax/hranoprovod-cli/v2/lib/resolver"
	"github.com/aquilax/hranoprovod-cli/v2/lib/shared"
	"github.com/urfave/cli/v2"
)

type registerCmd func(logStream, dbStream io.Reader, rc RegisterConfig) error

func newRegisterCommand(cu cmdUtils, register registerCmd) *cli.Command {
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
			return cu.withOptions(c, func(o *Options) error {
				return cu.withFileReaders([]string{o.GlobalConfig.DbFileName, o.GlobalConfig.LogFileName}, func(streams []io.Reader) error {
					dbStream, logStream := streams[0], streams[1]
					return register(logStream, dbStream, RegisterConfig{
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

type RegisterConfig struct {
	DateFormat     string
	ParserConfig   parser.Config
	ResolverConfig resolver.Config
	ReporterConfig reporter.Config
	FilterConfig   filter.Config
}

// Register generates report
func Register(logStream, dbStream io.Reader, rc RegisterConfig) error {
	rpCb := func(rpc reporter.Config, nl shared.DBNodeMap) reporter.Reporter {
		return reporter.NewRegReporter(rpc, nl)
	}
	return walkWithReporter(logStream, dbStream, rc.DateFormat, rc.ParserConfig, rc.ResolverConfig, rc.ReporterConfig, rc.FilterConfig, rpCb)
}