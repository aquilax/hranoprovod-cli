package cmd

import (
	"io"
	"time"

	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/aquilax/hranoprovod-cli/v2/options"
	"github.com/urfave/cli/v2"
)

type SummaryCmd func(logStream, dbStream io.Reader, sc app.SummaryConfig) error

func newSummaryCommand(cu cmdUtils, summary SummaryCmd) *cli.Command {
	return &cli.Command{
		Name:  "summary",
		Usage: "Show summary",
		Subcommands: []*cli.Command{
			newSummaryTodayCommand(cu, summary),
			newSummaryYesterdayCommand(cu, summary),
		},
	}
}

func newSummaryTodayCommand(cu cmdUtils, summary SummaryCmd) *cli.Command {
	return &cli.Command{
		Name:  "today",
		Usage: "Show summary for today",
		Action: func(c *cli.Context) error {
			return cu.withOptions(c, func(o *options.Options) error {
				return cu.withFileReaders([]string{o.GlobalConfig.DbFileName, o.GlobalConfig.LogFileName}, func(streams []io.Reader) error {
					dbStream, logStream := streams[0], streams[1]
					t := time.Now().Local()
					bTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
					o.FilterConfig.BeginningTime = &bTime
					eTime := time.Date(t.Year(), t.Month(), t.Day(), 24, 0, 0, -1, t.Location())
					o.FilterConfig.EndTime = &eTime
					return summary(logStream, dbStream, app.SummaryConfig{
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

func newSummaryYesterdayCommand(cu cmdUtils, summary SummaryCmd) *cli.Command {
	return &cli.Command{
		Name:  "yesterday",
		Usage: "Show summary for yesterday",
		Action: func(c *cli.Context) error {
			return cu.withOptions(c, func(o *options.Options) error {
				return cu.withFileReaders([]string{o.GlobalConfig.DbFileName, o.GlobalConfig.LogFileName}, func(streams []io.Reader) error {
					dbStream, logStream := streams[0], streams[1]
					t := time.Now().Local().AddDate(0, 0, -1)
					bTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
					o.FilterConfig.BeginningTime = &bTime
					eTime := time.Date(t.Year(), t.Month(), t.Day(), 24, 0, 0, -1, t.Location())
					o.FilterConfig.EndTime = &eTime
					return summary(logStream, dbStream, app.RegisterConfig{
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
