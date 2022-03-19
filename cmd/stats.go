package cmd

import (
	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/aquilax/hranoprovod-cli/v2/options"
	"github.com/urfave/cli/v2"
)

type statsCmd func(logFileName, dbFileName string, sc app.StatsConfig) error

func newStatsCommand(cu cmdUtils, stats statsCmd) *cli.Command {
	return &cli.Command{
		Name:  "stats",
		Usage: "Provide stats information",
		Action: func(c *cli.Context) error {
			return cu.withOptions(c, func(o *options.Options) error {
				return stats(o.GlobalConfig.LogFileName, o.GlobalConfig.LogFileName, app.StatsConfig{
					ParserConfig:   o.ParserConfig,
					ReporterConfig: o.ReporterConfig,
				})
			})
		},
	}
}
