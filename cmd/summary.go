package cmd

import (
	"io"
	"time"

	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/aquilax/hranoprovod-cli/v2/options"
	"github.com/urfave/cli/v2"
)

type summaryCmd func(logStream, dbStream io.Reader, sc app.SummaryConfig) error

func newSummaryCommand(cu cmdUtils, summary summaryCmd) *cli.Command {
	return &cli.Command{
		Name:  "summary",
		Usage: "Show summary for date",
		Action: func(c *cli.Context) error {
			return cu.withOptions(c, func(o *options.Options) error {
				t, err := options.GetTimeFromString(o.GlobalConfig.Now, o.GlobalConfig.DateFormat, c.Args().First())
				if err != nil {
					return err
				}
				return cu.withFileReaders([]string{o.GlobalConfig.DbFileName, o.GlobalConfig.LogFileName}, func(streams []io.Reader) error {
					dbStream, logStream := streams[0], streams[1]
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
