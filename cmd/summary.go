package cmd

import (
	"io"
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
			if o, err := ol(c); err != nil {
				return err
			} else {
				return withFileReader(o.GlobalConfig.DbFileName, func(dbStream io.Reader) error {
					return withFileReader(o.GlobalConfig.LogFileName, func(logStream io.Reader) error {
						t := time.Now().Local()
						btime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
						o.FilterConfig.BeginningTime = &btime
						etime := time.Date(t.Year(), t.Month(), t.Day(), 24, 0, 0, -1, t.Location())
						o.FilterConfig.EndTime = &etime
						return app.Summary(logStream, dbStream, o.GlobalConfig.DateFormat, o.ParserConfig, o.ResolverConfig, o.ReporterConfig, o.FilterConfig)
					})
				})
			}
		},
	}
}

func newSummaryYesterdayCommand(ol optionLoader) *cli.Command {
	return &cli.Command{
		Name:  "yesterday",
		Usage: "Show summary for yesterday",
		Action: func(c *cli.Context) error {
			if o, err := ol(c); err != nil {
				return err
			} else {
				return withFileReader(o.GlobalConfig.DbFileName, func(dbStream io.Reader) error {
					return withFileReader(o.GlobalConfig.LogFileName, func(logStream io.Reader) error {
						t := time.Now().Local().AddDate(0, 0, -1)
						btime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
						o.FilterConfig.BeginningTime = &btime
						etime := time.Date(t.Year(), t.Month(), t.Day(), 24, 0, 0, -1, t.Location())
						o.FilterConfig.EndTime = &etime
						return app.Summary(logStream, dbStream, o.GlobalConfig.DateFormat, o.ParserConfig, o.ResolverConfig, o.ReporterConfig, o.FilterConfig)
					})
				})
			}
		},
	}
}
