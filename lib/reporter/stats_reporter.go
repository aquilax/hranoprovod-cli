package reporter

import (
	"bufio"
	"fmt"
	"time"

	"github.com/aquilax/hranoprovod-cli/v2/lib/shared"
)

type Stats struct {
	DbFileName      string
	LogFileName     string
	DbRecordsCount  int
	LogRecordsCount int
	Now             time.Time
	LogFirstRecord  time.Time
	LogLastRecord   time.Time
}

type StatsReporter struct {
	config Config
	stats  *Stats
}

func NewStatsReporter(c Config, stats *Stats) *StatsReporter {
	return &StatsReporter{c, stats}
}

func (sr StatsReporter) Process(ln *shared.LogNode) error {
	return nil
}
func (sr StatsReporter) Flush() error {
	w := bufio.NewWriter(sr.config.Output)
	fmt.Fprintf(w, "  Database file:      %s\n", sr.stats.DbFileName)
	fmt.Fprintf(w, "  Database records:   %d\n", sr.stats.DbRecordsCount)
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "  Log file:           %s\n", sr.stats.LogFileName)
	fmt.Fprintf(w, "  Log records:        %d\n", sr.stats.LogRecordsCount)
	fmt.Fprintf(w, "  Today:              %s\n", sr.stats.Now.Format(sr.config.DateFormat))
	fmt.Fprintf(w, "  First record:       %s (%d days ago)\n", sr.stats.LogFirstRecord.Format(sr.config.DateFormat), int(sr.stats.Now.Sub(sr.stats.LogFirstRecord).Hours()/24))
	fmt.Fprintf(w, "  Last record:        %s (%d days ago)\n", sr.stats.LogLastRecord.Format(sr.config.DateFormat), int(sr.stats.Now.Sub(sr.stats.LogLastRecord).Hours()/24))
	return w.Flush()
}