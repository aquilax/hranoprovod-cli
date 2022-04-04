package summary

import (
	"io"
	"time"

	"github.com/aquilax/hranoprovod-cli/v2/cmd/hranoprovod-cli/internal/options"
	"github.com/aquilax/hranoprovod-cli/v2/cmd/hranoprovod-cli/internal/utils"
	"github.com/aquilax/hranoprovod-cli/v2/lib/filter"
	"github.com/aquilax/hranoprovod-cli/v2/lib/parser"
	"github.com/aquilax/hranoprovod-cli/v2/lib/reporter"
	"github.com/aquilax/hranoprovod-cli/v2/lib/resolver"
	"github.com/aquilax/hranoprovod-cli/v2/lib/shared"
	"github.com/urfave/cli/v2"
)

type summaryCmd func(logStream, dbStream io.Reader, sc SummaryConfig) error

func Command() *cli.Command {
	return NewSummaryCommand(utils.NewCmdUtils(), Summary)
}

func NewSummaryCommand(cu utils.CmdUtils, summary summaryCmd) *cli.Command {
	return &cli.Command{
		Name:  "summary",
		Usage: "Show summary for date",
		Action: func(c *cli.Context) error {
			return cu.WithOptions(c, func(o *options.Options) error {
				t, err := options.GetTimeFromString(o.GlobalConfig.Now, o.GlobalConfig.DateFormat, c.Args().First())
				if err != nil {
					return err
				}
				return cu.WithFileReaders([]string{o.GlobalConfig.DbFileName, o.GlobalConfig.LogFileName}, func(streams []io.Reader) error {
					dbStream, logStream := streams[0], streams[1]
					bTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
					o.FilterConfig.BeginningTime = &bTime
					eTime := time.Date(t.Year(), t.Month(), t.Day(), 24, 0, 0, -1, t.Location())
					o.FilterConfig.EndTime = &eTime
					return summary(logStream, dbStream, SummaryConfig{
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

type SummaryConfig struct {
	DateFormat     string
	ParserConfig   parser.Config
	ResolverConfig resolver.Config
	ReporterConfig reporter.Config
	FilterConfig   filter.Config
}

// Summary generates summary
func Summary(logStream, dbStream io.Reader, sc SummaryConfig) error {
	return utils.WithResolvedDatabase(dbStream, sc.ParserConfig, sc.ResolverConfig,
		func(nl shared.DBNodeMap) error {
			r := reporter.NewSummaryReporterTemplate(sc.ReporterConfig, nl)
			f := filter.GetIntervalNodeFilter(sc.FilterConfig)
			return utils.WalkNodesInStream(logStream, sc.DateFormat, sc.ParserConfig, f, r)
		})
}
