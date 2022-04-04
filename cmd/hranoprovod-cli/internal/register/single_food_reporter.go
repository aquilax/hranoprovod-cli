package register

import (
	"bufio"
	"fmt"
	"regexp"

	"github.com/aquilax/hranoprovod-cli/v2/cmd/hranoprovod-cli/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/v2/lib/shared"
)

type singleFoodReporter struct {
	config reporter.Config
	db     shared.DBNodeMap
	output *bufio.Writer
}

func newSingleFoodReporter(config reporter.Config, db shared.DBNodeMap) *singleFoodReporter {
	return &singleFoodReporter{
		config,
		db,
		bufio.NewWriter(config.Output),
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
	return r.output.Flush()
}
