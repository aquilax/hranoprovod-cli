package register

import (
	"bufio"
	"fmt"
	"regexp"

	"github.com/aquilax/hranoprovod-cli/v3/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/node"
)

type singleFoodReporter struct {
	config reporter.Config
	db     node.DBNodeMap
	output *bufio.Writer
}

func newSingleFoodReporter(config reporter.Config, db node.DBNodeMap) *singleFoodReporter {
	return &singleFoodReporter{
		config,
		db,
		bufio.NewWriter(config.Output),
	}
}

func (r *singleFoodReporter) Process(ln *node.LogNode) error {
	for _, e := range ln.Elements {
		matched, err := regexp.MatchString(r.config.SingleFood, e.Name)
		if err != nil {
			return err
		}
		if matched {
			if _, err := fmt.Fprintf(r.output, "%s\t%s\t%0.2f\n", ln.Time.Format(r.config.DateFormat), e.Name, e.Value); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *singleFoodReporter) Flush() error {
	return r.output.Flush()
}
