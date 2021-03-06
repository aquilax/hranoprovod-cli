package reporter

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/aquilax/hranoprovod-cli/v2/accumulator"
	"github.com/aquilax/hranoprovod-cli/v2/shared"
)

type regReporter struct {
	options *Options
	db      shared.DBNodeList
	output  io.Writer
}

func newRegReporter(options *Options, db shared.DBNodeList, writer io.Writer) *regReporter {
	return &regReporter{
		options,
		db,
		writer,
	}
}

func (r *regReporter) Process(ln *shared.LogNode) error {
	acc := accumulator.NewAccumulator()
	r.printDate(ln.Time)
	for _, element := range ln.Elements {
		if !r.options.TotalsOnly {
			r.printElement(element)
		}
		if repl, found := r.db[element.Name]; found {
			for _, repl := range repl.Elements {
				res := repl.Val * element.Val
				if !r.options.TotalsOnly {
					r.printIngredient(repl.Name, res)
				}
				acc.Add(repl.Name, res)
			}
		} else {
			if !r.options.TotalsOnly {
				r.printIngredient(element.Name, element.Val)
			}
			acc.Add(element.Name, element.Val)
		}
	}
	if r.options.Totals {
		var ss sort.StringSlice
		if len(acc) > 0 {
			r.printTotalHeader()
			for name := range acc {
				ss = append(ss, name)
			}
			sort.Sort(ss)
			for _, name := range ss {
				arr := acc[name]
				r.printTotalRow(name, arr[accumulator.Positive], arr[accumulator.Negative])
			}
		}
	}
	return nil
}

func (r *regReporter) Flush() error {
	return nil
}

func (r *regReporter) cNum(num float64) string {
	if r.options.Color {
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
	fmt.Fprintf(r.output, "%s\n", ts.Format(r.options.DateFormat))
}

func (r *regReporter) printElement(element shared.Element) {
	fmt.Fprintf(r.output, "\t%-27s :%s\n", element.Name, r.cNum(element.Val))
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
