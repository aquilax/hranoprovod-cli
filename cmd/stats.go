package cmd

import (
	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/urfave/cli/v2"
)

func newStatsCommand(cu cmdUtils) *cli.Command {
	return &cli.Command{
		Name:  "stats",
		Usage: "Provide stats information",
		Action: func(c *cli.Context) error {
			return cu.withOptions(c, func(o *app.Options) error {
				return app.Stats(o.GlobalConfig, o.ParserConfig, o.ReporterConfig)
			})
		},
	}
}
