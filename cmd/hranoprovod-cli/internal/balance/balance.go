package balance

import (
	"io"

	"github.com/aquilax/hranoprovod-cli/v2/cmd/hranoprovod-cli/internal/options"
	"github.com/aquilax/hranoprovod-cli/v2/cmd/hranoprovod-cli/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/v2/cmd/hranoprovod-cli/internal/utils"
	"github.com/aquilax/hranoprovod-cli/v2/lib/filter"
	"github.com/aquilax/hranoprovod-cli/v2/lib/parser"
	"github.com/aquilax/hranoprovod-cli/v2/lib/resolver"
	"github.com/aquilax/hranoprovod-cli/v2/lib/shared"
	"github.com/urfave/cli/v2"
)

type balanceCmd func(logStream, dbStream io.Reader, bc BalanceConfig) error

func Command() *cli.Command {
	return NewBalanceCommand(utils.NewCmdUtils(), Balance)
}

func NewBalanceCommand(cu utils.CmdUtils, balance balanceCmd) *cli.Command {
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
			return cu.WithOptions(c, func(o *options.Options) error {
				return cu.WithFileReaders([]string{o.GlobalConfig.DbFileName, o.GlobalConfig.LogFileName}, func(streams []io.Reader) error {
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

type BalanceConfig struct {
	DateFormat     string
	ParserConfig   parser.Config
	ResolverConfig resolver.Config
	ReporterConfig reporter.Config
	FilterConfig   filter.Config
}

// Balance generates balance report
func Balance(logStream, dbStream io.Reader, bc BalanceConfig) error {
	return utils.WithResolvedDatabase(dbStream, bc.ParserConfig, bc.ResolverConfig,
		func(nl shared.DBNodeMap) error {
			r := NewBalanceReporter(bc.ReporterConfig, nl)
			f := filter.GetIntervalNodeFilter(bc.FilterConfig)
			return utils.WalkNodesInStream(logStream, bc.DateFormat, bc.ParserConfig, f, r)
		})
}
