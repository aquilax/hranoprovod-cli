package main

import (
	"io"

	"github.com/aquilax/hranoprovod-cli/v2/lib/filter"
	"github.com/aquilax/hranoprovod-cli/v2/lib/reporter"
	"github.com/aquilax/hranoprovod-cli/v2/lib/shared"
	"github.com/urfave/cli/v2"
)

type balanceCmd func(logStream, dbStream io.Reader, bc BalanceConfig) error

func newBalanceCommand(cu cmdUtils, balance balanceCmd) *cli.Command {
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
			return cu.withOptions(c, func(o *Options) error {
				return cu.withFileReaders([]string{o.GlobalConfig.DbFileName, o.GlobalConfig.LogFileName}, func(streams []io.Reader) error {
					dbStream, logStream := streams[0], streams[1]
					return balance(logStream, dbStream, BalanceConfig{
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

type BalanceConfig = RegisterConfig

// Balance generates balance report
func Balance(logStream, dbStream io.Reader, bc BalanceConfig) error {
	return withResolvedDatabase(dbStream, bc.ParserConfig, bc.ResolverConfig,
		func(nl shared.DBNodeMap) error {
			r := reporter.NewBalanceReporter(bc.ReporterConfig, nl)
			f := filter.GetIntervalNodeFilter(bc.FilterConfig)
			return walkNodesInStream(logStream, bc.DateFormat, bc.ParserConfig, f, r)
		})
}
