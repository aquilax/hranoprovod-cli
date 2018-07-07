package reporter

import (
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/aquilax/hranoprovod-cli/shared"
)

// Reporter is the main report structure
type SingleFoodReporter struct {
	options *Options
	db      *shared.NodeList
	output  io.Writer
}

// NewReporter creates new reporter
func NewSingleFoodReporter(options *Options, db *shared.NodeList, writer io.Writer) *SingleFoodReporter {
	return &SingleFoodReporter{
		options,
		db,
		writer,
	}
}

func (r *SingleFoodReporter) Process(ln *shared.LogNode) error {
	for _, e := range *ln.Elements {
		matched, err := regexp.MatchString(r.options.SingleFood, e.Name)
		if err != nil {
			return err
		}
		if matched {
			r.PrintSingleFoodRow(ln.Time, e.Name, e.Val)
		}
	}
	return nil
}

func (r *SingleFoodReporter) Flush() error {
	return nil
}

func (r *SingleFoodReporter) PrintSingleFoodRow(ts time.Time, name string, val float32) {
	fmt.Fprintf(r.output, "%s\t%s\t%0.2f\n", ts.Format(r.options.DateFormat), name, val)
}
