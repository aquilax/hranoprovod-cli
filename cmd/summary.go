package cmd

import (
	"time"

	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/urfave/cli/v2"
)

func newSummaryCommand(ol optionLoader) *cli.Command {
	return &cli.Command{
		Name:  "summary",
		Usage: "Show summary",
		Subcommands: []*cli.Command{
			newSummaryTodayCommand(ol),
			newSummaryYesterdayCommand(ol),
		},
	}
}

func newSummaryTodayCommand(ol optionLoader) *cli.Command {
	return &cli.Command{
		Name:  "today",
		Usage: "Show summary for today",
		Action: func(c *cli.Context) error {
			o, err := ol(c)
			if err != nil {
				return err
			}
			t := time.Now().Local()
			o.ReporterConfig.HasBeginning = true
			o.ReporterConfig.BeginningTime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
			o.ReporterConfig.HasEnd = true
			o.ReporterConfig.EndTime = time.Date(t.Year(), t.Month(), t.Day(), 24, 0, 0, -1, t.Location())
			return app.NewHranoprovod(o).Summary(o.ParserConfig, o.ResolverConfig, o.ReporterConfig)
		},
	}
}

func newSummaryYesterdayCommand(ol optionLoader) *cli.Command {
	return &cli.Command{
		Name:  "yesterday",
		Usage: "Show summary for yesterday",
		Action: func(c *cli.Context) error {
			o, err := ol(c)
			if err != nil {
				return err
			}
			t := time.Now().Local().AddDate(0, 0, -1)
			o.ReporterConfig.HasBeginning = true
			o.ReporterConfig.BeginningTime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
			o.ReporterConfig.HasEnd = true
			o.ReporterConfig.EndTime = time.Date(t.Year(), t.Month(), t.Day(), 24, 0, 0, -1, t.Location())
			return app.NewHranoprovod(o).Summary(o.ParserConfig, o.ResolverConfig, o.ReporterConfig)
		},
	}
}
