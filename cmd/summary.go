package cmd

import (
	"io"
	"time"

	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/urfave/cli/v2"
)

func newSummaryCommand(cu cmdUtils) *cli.Command {
	return &cli.Command{
		Name:  "summary",
		Usage: "Show summary",
		Subcommands: []*cli.Command{
			newSummaryTodayCommand(cu),
			newSummaryYesterdayCommand(cu),
		},
	}
}

func newSummaryTodayCommand(cu cmdUtils) *cli.Command {
	return &cli.Command{
		Name:  "today",
		Usage: "Show summary for today",
		Action: func(c *cli.Context) error {
			return cu.withOptions(c, func(o *app.Options) error {
				return cu.withFileReaders([]string{o.GlobalConfig.DbFileName, o.GlobalConfig.LogFileName}, func(streams []io.Reader) error {
					dbStream, logStream := streams[0], streams[1]
					t := time.Now().Local()
					bTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
					o.FilterConfig.BeginningTime = &bTime
					eTime := time.Date(t.Year(), t.Month(), t.Day(), 24, 0, 0, -1, t.Location())
					o.FilterConfig.EndTime = &eTime
					return app.Summary(logStream, dbStream, o.GlobalConfig.DateFormat, o.ParserConfig, o.ResolverConfig, o.ReporterConfig, o.FilterConfig)
				})
			})
		},
	}
}

func newSummaryYesterdayCommand(cu cmdUtils) *cli.Command {
	return &cli.Command{
		Name:  "yesterday",
		Usage: "Show summary for yesterday",
		Action: func(c *cli.Context) error {
			return cu.withOptions(c, func(o *app.Options) error {
				return cu.withFileReaders([]string{o.GlobalConfig.DbFileName, o.GlobalConfig.LogFileName}, func(streams []io.Reader) error {
					dbStream, logStream := streams[0], streams[1]
					t := time.Now().Local().AddDate(0, 0, -1)
					bTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
					o.FilterConfig.BeginningTime = &bTime
					eTime := time.Date(t.Year(), t.Month(), t.Day(), 24, 0, 0, -1, t.Location())
					o.FilterConfig.EndTime = &eTime
					return app.Summary(logStream, dbStream, o.GlobalConfig.DateFormat, o.ParserConfig, o.ResolverConfig, o.ReporterConfig, o.FilterConfig)
				})
			})
		},
	}
}
