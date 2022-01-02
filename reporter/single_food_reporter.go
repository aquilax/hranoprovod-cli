package reporter

import (
	"fmt"
	"io"
	"regexp"

	"github.com/aquilax/hranoprovod-cli/v2/shared"
)

type singleFoodReporter struct {
	config Config
	db     shared.DBNodeList
	output io.Writer
}

func newSingleFoodReporter(config Config, db shared.DBNodeList) *singleFoodReporter {
	return &singleFoodReporter{
		config,
		db,
		config.Output,
	}
}

func (r *singleFoodReporter) Process(ln *shared.LogNode) error {
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
	return nil
}
