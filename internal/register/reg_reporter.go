package register

import (
	"bufio"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aquilax/hranoprovod-cli/v3/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/accumulator"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/element"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/node"
)

const (
	reset = "\x1B[0m"
	red   = "\x1B[31m"
	green = "\x1B[32m"
)

type regReporter struct {
	config reporter.Config
	db     node.DBNodeMap
	output *bufio.Writer
}

// NewRegReporter creates new response handler
func NewRegReporter(c reporter.Config, db node.DBNodeMap) reporter.Reporter {
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

func newRegReporter(config reporter.Config, db node.DBNodeMap) *regReporter {
	return &regReporter{
		config,
		db,
		bufio.NewWriter(config.Output),
	}
}

func (r *regReporter) Process(ln *node.LogNode) error {
	acc := accumulator.NewAccumulator()
	r.mustPrintDate(ln.Time)
	for _, element := range ln.Elements {
		if !r.config.TotalsOnly {
			r.mustPrintElement(element)
		}
		if repl, found := r.db[element.Name]; found {
			for _, repl := range repl.Elements {
				res := repl.Value * element.Value
				if !r.config.TotalsOnly {
					r.mustPrintIngredient(repl.Name, res)
				}
				acc.Add(repl.Name, res)
			}
		} else {
			if !r.config.TotalsOnly {
				r.mustPrintIngredient(element.Name, element.Value)
			}
			acc.Add(element.Name, element.Value)
		}
	}
	if r.config.Totals {
		var ss sort.StringSlice
		if len(acc) > 0 {
			r.mustPrintTotalHeader()
			for name := range acc {
				ss = append(ss, name)
			}
			sort.Sort(ss)
			for _, name := range ss {
				arr := acc[name]
				r.mustPrintTotalRow(name, arr[accumulator.Positive], arr[accumulator.Negative])
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

func (r *regReporter) mustPrintDate(ts time.Time) {
	if _, err := fmt.Fprintf(r.output, "%s\n", ts.Format(r.config.DateFormat)); err != nil {
		panic(err)
	}
}

func (r *regReporter) mustPrintElement(element element.Element) {
	if _, err := fmt.Fprintf(r.output, "\t%-27s :%s\n", element.Name, r.cNum(element.Value)); err != nil {
		panic(err)
	}
}

func (r *regReporter) mustPrintIngredient(name string, value float64) {
	if _, err := fmt.Fprintf(r.output, "\t\t%20s %s\n", name, r.cNum(value)); err != nil {
		panic(err)
	}
}

func (r *regReporter) mustPrintTotalHeader() {
	if _, err := fmt.Fprintf(r.output, "\t-- %s %s\n", "TOTAL ", strings.Repeat("-", 52)); err != nil {
		panic(err)
	}
}

func (r *regReporter) mustPrintTotalRow(name string, pos float64, neg float64) {
	if _, err := fmt.Fprintf(r.output, "\t\t%20s %s %s =%s\n", name, r.cNum(pos), r.cNum(neg), r.cNum(pos+neg)); err != nil {
		panic(err)
	}
}
