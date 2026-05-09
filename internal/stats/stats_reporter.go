package stats

import (
	"bufio"
	"time"

	"github.com/aquilax/hranoprovod-cli/v3/internal/errwriter"
	"github.com/aquilax/hranoprovod-cli/v3/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/node"
)

type StatsData struct {
	DbFileName      string
	LogFileName     string
	DbRecordsCount  int
	LogRecordsCount int
	Now             time.Time
	LogFirstRecord  time.Time
	LogLastRecord   time.Time
}

type StatsReporter struct {
	stats      *StatsData
	output     *bufio.Writer
	dateFormat string
}

func NewStatsReporter(c reporter.Config, stats *StatsData) *StatsReporter {
	return &StatsReporter{
		stats:      stats,
		output:     bufio.NewWriter(c.Output),
		dateFormat: c.DateFormat,
	}
}

func (sr StatsReporter) Process(ln *node.LogNode) error {
	return nil
}
func (sr StatsReporter) Flush() error {
	writer := errwriter.New(sr.output)

	writer.Fprintf("  Database file:      %s\n", sr.stats.DbFileName)
	writer.Fprintf("  Database records:   %d\n", sr.stats.DbRecordsCount)
	writer.Fprintf("\n")
	writer.Fprintf("  Log file:           %s\n", sr.stats.LogFileName)
	writer.Fprintf("  Log records:        %d\n", sr.stats.LogRecordsCount)
	writer.Fprintf("  Today:              %s\n", sr.stats.Now.Format(sr.dateFormat))
	writer.Fprintf("  First record:       %s (%d days ago)\n", sr.stats.LogFirstRecord.Format(sr.dateFormat), int(sr.stats.Now.Sub(sr.stats.LogFirstRecord).Hours()/24))
	writer.Fprintf("  Last record:        %s (%d days ago)\n", sr.stats.LogLastRecord.Format(sr.dateFormat), int(sr.stats.Now.Sub(sr.stats.LogLastRecord).Hours()/24))
	if writer.Err != nil {
		return writer.Err
	}
	return sr.output.Flush()
}
