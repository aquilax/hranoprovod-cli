package register

import (
	"bufio"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aquilax/hranoprovod-cli/v2/cmd/hranoprovod-cli/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/v2/lib/shared"
)

const (
	reset = "\x1B[0m"
	red   = "\x1B[31m"
	green = "\x1B[32m"
)

type regReporter struct {
	config reporter.Config
	db     shared.DBNodeMap
	output *bufio.Writer
}

// NewRegReporter creates new response handler
func NewRegReporter(c reporter.Config, db shared.DBNodeMap) reporter.Reporter {
	if len(c.SingleElement) > 0 {
		if c.ElementGroupByFood {
			return newElementByFoodReporter(c, db)
		}
		return newSingleReporter(c, db)
	}
	if len(c.SingleFood) > 0 {
		return newSingleFoodReporter(c, db)
	}
	if c.UseOldRegReporter {
		return newRegReporter(c, db)
	}
	return newRegReporterTemplate(c, db)
}

func newRegReporter(config reporter.Config, db shared.DBNodeMap) *regReporter {
	return &regReporter{
		config,
		db,
		bufio.NewWriter(config.Output),
	}
}

func (r *regReporter) Process(ln *shared.LogNode) error {
	acc := shared.NewAccumulator()
	r.printDate(ln.Time)
	for _, element := range ln.Elements {
		if !r.config.TotalsOnly {
			r.printElement(element)
		}
		if repl, found := r.db[element.Name]; found {
			for _, repl := range repl.Elements {
				res := repl.Value * element.Value
				if !r.config.TotalsOnly {
					r.printIngredient(repl.Name, res)
				}
				acc.Add(repl.Name, res)
			}
		} else {
			if !r.config.TotalsOnly {
				r.printIngredient(element.Name, element.Value)
			}
			acc.Add(element.Name, element.Value)
		}
	}
	if r.config.Totals {
		var ss sort.StringSlice
		if len(acc) > 0 {
			r.printTotalHeader()
			for name := range acc {
				ss = append(ss, name)
			}
			sort.Sort(ss)
			for _, name := range ss {
				arr := acc[name]
				r.printTotalRow(name, arr[shared.Positive], arr[shared.Negative])
			}
		}
	}
	return nil
}

func (r *regReporter) Flush() error {
	return r.output.Flush()
}

func (r *regReporter) cNum(num float64) string {
	if r.config.Color {
		if num > 0 {
			return red + fmt.Sprintf("%10.2f", num) + reset
		}
		if num < 0 {
			return green + fmt.Sprintf("%10.2f", num) + reset
		}
	}
	return fmt.Sprintf("%10.2f", num)
}

func (r *regReporter) printDate(ts time.Time) {
	fmt.Fprintf(r.output, "%s\n", ts.Format(r.config.DateFormat))
}

func (r *regReporter) printElement(element shared.Element) {
	fmt.Fprintf(r.output, "\t%-27s :%s\n", element.Name, r.cNum(element.Value))
}

func (r *regReporter) printIngredient(name string, value float64) {
	fmt.Fprintf(r.output, "\t\t%20s %s\n", name, r.cNum(value))
}

func (r *regReporter) printTotalHeader() {
	fmt.Fprintf(r.output, "\t-- %s %s\n", "TOTAL ", strings.Repeat("-", 52))
}

func (r *regReporter) printTotalRow(name string, pos float64, neg float64) {
	fmt.Fprintf(r.output, "\t\t%20s %s %s =%s\n", name, r.cNum(pos), r.cNum(neg), r.cNum(pos+neg))
}
