package stats

import (
	"bufio"
	"fmt"
	"time"

	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/reporter"
	shared "github.com/aquilax/hranoprovod-cli/v3"
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

func (sr StatsReporter) Process(ln *shared.LogNode) error {
	return nil
}
func (sr StatsReporter) Flush() error {
	fmt.Fprintf(sr.output, "  Database file:      %s\n", sr.stats.DbFileName)
	fmt.Fprintf(sr.output, "  Database records:   %d\n", sr.stats.DbRecordsCount)
	fmt.Fprintln(sr.output, "")
	fmt.Fprintf(sr.output, "  Log file:           %s\n", sr.stats.LogFileName)
	fmt.Fprintf(sr.output, "  Log records:        %d\n", sr.stats.LogRecordsCount)
	fmt.Fprintf(sr.output, "  Today:              %s\n", sr.stats.Now.Format(sr.dateFormat))
	fmt.Fprintf(sr.output, "  First record:       %s (%d days ago)\n", sr.stats.LogFirstRecord.Format(sr.dateFormat), int(sr.stats.Now.Sub(sr.stats.LogFirstRecord).Hours()/24))
	fmt.Fprintf(sr.output, "  Last record:        %s (%d days ago)\n", sr.stats.LogLastRecord.Format(sr.dateFormat), int(sr.stats.Now.Sub(sr.stats.LogLastRecord).Hours()/24))
	return sr.output.Flush()
}
