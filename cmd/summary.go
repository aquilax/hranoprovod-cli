package cmd

import (
	"time"

	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/urfave/cli/v2"
)

func newSummaryCommand() *cli.Command {
	return &cli.Command{
		Name:  "summary",
		Usage: "Show summary",
		Subcommands: []*cli.Command{
			{
				Name:  "today",
				Usage: "Show summary for today",
				Action: func(c *cli.Context) error {
					o := app.NewOptions()
					if err := o.Load(c); err != nil {
						return err
					}
					t := time.Now().Local()
					o.Reporter.HasBeginning = true
					o.Reporter.BeginningTime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
					o.Reporter.HasEnd = true
					o.Reporter.EndTime = time.Date(t.Year(), t.Month(), t.Day(), 24, 0, 0, -1, t.Location())
					return app.NewHranoprovod(o).Summary()
				},
			},
			{
				Name:  "yesterday",
				Usage: "Show summary for yesterday",
				Action: func(c *cli.Context) error {
					o := app.NewOptions()
					if err := o.Load(c); err != nil {
						return err
					}
					t := time.Now().Local().AddDate(0, 0, -1)
					o.Reporter.HasBeginning = true
					o.Reporter.BeginningTime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
					o.Reporter.HasEnd = true
					o.Reporter.EndTime = time.Date(t.Year(), t.Month(), t.Day(), 24, 0, 0, -1, t.Location())
					return app.NewHranoprovod(o).Summary()
				},
			},
		},
	}
}
