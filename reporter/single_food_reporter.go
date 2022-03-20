package reporter

import (
	"bufio"
	"fmt"
	"regexp"

	"github.com/aquilax/hranoprovod-cli/v2"
)

type singleFoodReporter struct {
	config Config
	db     hranoprovod.DBNodeMap
	output *bufio.Writer
}

func newSingleFoodReporter(config Config, db hranoprovod.DBNodeMap) *singleFoodReporter {
	return &singleFoodReporter{
		config,
		db,
		bufio.NewWriter(config.Output),
	}
}

func (r *singleFoodReporter) Process(ln *hranoprovod.LogNode) error {
	for _, e := range ln.Elements {
		matched, err := regexp.MatchString(r.config.SingleFood, e.Name)
		if err != nil {
			return err
		}
		if matched {
			fmt.Fprintf(r.output, "%s\t%s\t%0.2f\n", ln.Time.Format(r.config.DateFormat), e.Name, e.Value)
		}
	}
	return nil
}

func (r *singleFoodReporter) Flush() error {
	return r.output.Flush()
}
