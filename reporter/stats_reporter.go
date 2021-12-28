package reporter

import (
	"bufio"

	"github.com/aquilax/hranoprovod-cli/v2/shared"
)

type StatsReporter struct {
	config Config
	stats  []string
}

func NewStatsReporter(c Config, stats []string) *StatsReporter {
	return &StatsReporter{c, stats}
}

func (sr StatsReporter) Process(ln *shared.LogNode) error {
	return nil
}
func (sr StatsReporter) Flush() error {
	w := bufio.NewWriter(sr.config.Output)
	for i := range sr.stats {
		w.WriteString(sr.stats[i])
	}
	return w.Flush()
}
