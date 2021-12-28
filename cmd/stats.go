package cmd

import (
	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/urfave/cli/v2"
)

func newStatsCommand(ol optionLoader) *cli.Command {
	return &cli.Command{
		Name:  "stats",
		Usage: "Provide stats information",
		Action: func(c *cli.Context) error {
			o, err := ol(c)
			if err != nil {
				return err
			}
			return app.NewHranoprovod(o).Stats(o.GlobalConfig, o.ParserConfig, o.ReporterConfig)
		},
	}
}
