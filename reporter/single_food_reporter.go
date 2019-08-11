package reporter

import (
	"fmt"
	"io"
	"regexp"

	"github.com/aquilax/hranoprovod-cli/shared"
)

type singleFoodReporter struct {
	options *Options
	db      shared.DBNodeList
	output  io.Writer
}

func newSingleFoodReporter(options *Options, db shared.DBNodeList, writer io.Writer) *singleFoodReporter {
	return &singleFoodReporter{
		options,
		db,
		writer,
	}
}

func (r *singleFoodReporter) Process(ln *shared.LogNode) error {
	for _, e := range ln.Elements {
		matched, err := regexp.MatchString(r.options.SingleFood, e.Name)
		if err != nil {
			return err
		}
		if matched {
			fmt.Fprintf(r.output, "%s\t%s\t%0.2f\n", ln.Time.Format(r.options.DateFormat), e.Name, e.Val)
		}
	}
	return nil
}

func (r *singleFoodReporter) Flush() error {
	return nil
}
