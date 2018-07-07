package reporter

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/aquilax/hranoprovod-cli/accumulator"
	"github.com/aquilax/hranoprovod-cli/shared"
)

// Reporter is the main report structure
type RegReporter struct {
	options *Options
	db      *shared.NodeList
	output  io.Writer
}

// NewReporter creates new reporter
func NewRegReporter(options *Options, db *shared.NodeList, writer io.Writer) *RegReporter {
	return &RegReporter{
		options,
		db,
		writer,
	}
}

func (r *RegReporter) Process(ln *shared.LogNode) error {
	acc := accumulator.NewAccumulator()
	r.PrintDate(ln.Time)
	for _, element := range *ln.Elements {
		if !r.options.TotalsOnly {
			r.PrintElement(element)
		}
		if repl, found := (*r.db)[element.Name]; found {
			for _, repl := range *repl.Elements {
				res := repl.Val * element.Val
				if !r.options.TotalsOnly {
					r.PrintIngredient(repl.Name, res)
				}
				acc.Add(repl.Name, res)
			}
		} else {
			if !r.options.TotalsOnly {
				r.PrintIngredient(element.Name, element.Val)
			}
			acc.Add(element.Name, element.Val)
		}
	}
	if r.options.Totals {
		var ss sort.StringSlice
		if len(*acc) > 0 {
			r.PrintTotalHeader()
			for name := range *acc {
				ss = append(ss, name)
			}
			sort.Sort(ss)
			for _, name := range ss {
				arr := (*acc)[name]
				r.PrintTotalRow(name, arr[accumulator.Positive], arr[accumulator.Negative])
			}
		}
	}
	return nil
}

func (r *RegReporter) Flush() error {
	return nil
}

func (r *RegReporter) cNum(num float32) string {
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

func (r *RegReporter) PrintDate(ts time.Time) {
	fmt.Fprintf(r.output, "%s\n", ts.Format(r.options.DateFormat))
}

func (r *RegReporter) PrintElement(element *shared.Element) {
	fmt.Fprintf(r.output, "\t%-27s :%s\n", element.Name, r.cNum(element.Val))
}

func (r *RegReporter) PrintIngredient(name string, value float32) {
	fmt.Fprintf(r.output, "\t\t%20s %s\n", name, r.cNum(value))
}

func (r *RegReporter) PrintTotalHeader() {
	fmt.Fprintf(r.output, "\t-- %s %s\n", "TOTAL ", strings.Repeat("-", 52))
}

func (r *RegReporter) PrintTotalRow(name string, pos float32, neg float32) {
	fmt.Fprintf(r.output, "\t\t%20s %s %s =%s\n", name, r.cNum(pos), r.cNum(neg), r.cNum(pos+neg))
}

func (r *RegReporter) PrintUnresolvedRow(name string) {
	fmt.Fprintln(r.output, name)
}
