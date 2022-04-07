package stats

import (
	"time"

	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/options"
	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/utils"
	"github.com/aquilax/hranoprovod-cli/lib/parser/v3"
	shared "github.com/aquilax/hranoprovod-cli/v3"
	"github.com/urfave/cli/v2"
)

type statsCmd func(logFileName, dbFileName string, sc StatsConfig) error

func Command() *cli.Command {
	return NewStatsCommand(utils.NewCmdUtils(), Stats)
}

func NewStatsCommand(cu utils.CmdUtils, stats statsCmd) *cli.Command {
	return &cli.Command{
		Name:  "stats",
		Usage: "Provide stats information",
		Action: func(c *cli.Context) error {
			return cu.WithOptions(c, func(o *options.Options) error {
				return stats(o.GlobalConfig.LogFileName, o.GlobalConfig.LogFileName, StatsConfig{
					Now:            o.GlobalConfig.Now,
					ParserConfig:   o.ParserConfig,
					ReporterConfig: o.ReporterConfig,
				})
			})
		},
	}
}

type StatsConfig struct {
	Now            time.Time
	ParserConfig   parser.Config
	ReporterConfig reporter.Config
}

// Stats generates statistics report
func Stats(logFileName, dbFileName string, sc StatsConfig) error {
	var err error
	var firstLogDate time.Time
	var lastLogDate time.Time

	countLog := 0
	if err = parser.ParseFileCallback(logFileName, sc.ParserConfig, func(n *shared.ParserNode, _ error) (stop bool, cbError error) {
		lastLogDate, err = time.Parse(sc.ReporterConfig.DateFormat, n.Header)
		if err == nil {
			if firstLogDate.IsZero() {
				firstLogDate = lastLogDate
			}
		}
		countLog++
		return false, nil
	}); err != nil {
		return err
	}

	countDb := 0
	if err = parser.ParseFileCallback(dbFileName, sc.ParserConfig, func(n *shared.ParserNode, _ error) (stop bool, cbError error) {
		countDb++
		return false, nil
	}); err != nil {
		return err
	}

	return NewStatsReporter(sc.ReporterConfig, &StatsData{
		DbFileName:      dbFileName,
		LogFileName:     logFileName,
		DbRecordsCount:  countDb,
		LogRecordsCount: countLog,
		Now:             sc.Now,
		LogFirstRecord:  firstLogDate,
		LogLastRecord:   lastLogDate,
	}).Flush()
}
