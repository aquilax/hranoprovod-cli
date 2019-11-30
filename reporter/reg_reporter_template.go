package reporter

import (
	"fmt"
	"io"
	"sort"
	"text/template"
	"time"

	"github.com/aquilax/hranoprovod-cli/accumulator"
	"github.com/aquilax/hranoprovod-cli/shared"
)

const dayTemplate = `{{printDate .Time}}
{{- if .Elements }}
{{- range $el := .Elements}}
{{ printf "\t%-27s :%s" (shorten $el.Name 27) (cNum $el.Val) }}
{{- range $ing := $el.Ingredients}}
{{ printf "\t\t%20s %s" (shorten $ing.Name 20) (cNum $ing.Val) }}
{{- end}}
{{- end}}
{{- end}}
{{- if .Totals }}
	-- TOTAL  ----------------------------------------------------
{{- range $total := .Totals }}
{{ printf "\t\t%20s %s %s =%s" (shorten $total.Name 20) (cNum $total.Positive) (cNum $total.Negative) (cNum $total.Sum) }}
{{- end}}
{{- end}}
`

type ReportElement struct {
	shared.Element
	Ingredients shared.Elements
}

type Total struct {
	Name     string
	Positive float64
	Negative float64
	Sum      float64
}

type ReportItem struct {
	Time     time.Time
	Elements *[]ReportElement
	Totals   *[]Total
}

type regReporterTemplate struct {
	options  *Options
	db       shared.DBNodeList
	output   io.Writer
	template *template.Template
}

func newRegReporterTemplate(options *Options, db shared.DBNodeList, writer io.Writer) *regReporterTemplate {
	var shorter = shorten
	if !options.ShortenStrings {
		shorter = func(s string, n int) string { return s }
	}
	return &regReporterTemplate{
		options,
		db,
		writer,
		template.Must(template.New("dayTemplate").Funcs(template.FuncMap{
			"printDate": func(ts time.Time) string {
				return ts.Format(options.DateFormat)
			},
			"cNum": func(num float64) string {
				if options.Color {
					if num > 0 {
						return red + fmt.Sprintf("%10.2f", num) + reset
					}
					if num < 0 {
						return green + fmt.Sprintf("%10.2f", num) + reset
					}
				}
				return fmt.Sprintf("%10.2f", num)
			},
			"shorten": shorter,
			"add": func(num1, num2 float64) float64 {
				return num1 + num2
			},
		}).Parse(dayTemplate)),
	}
}

func (r *regReporterTemplate) Process(ln *shared.LogNode) error {
	return r.template.Execute(r.output, r.getReportItem(ln, r.db))
}

func (r *regReporterTemplate) Flush() error {
	return nil
}

func newTotalFromAccumulator(acc accumulator.Accumulator) *[]Total {
	var result = make([]Total, len(acc))
	var ss = make(sort.StringSlice, len(acc))
	i := 0
	for name := range acc {
		ss[i] = name
		i++
	}
	sort.Sort(ss)
	for i, name := range ss {
		result[i] = Total{name, acc[name][accumulator.Positive], acc[name][accumulator.Negative], acc[name][accumulator.Positive] + acc[name][accumulator.Negative]}
	}
	return &result
}

func (r *regReporterTemplate) getReportItem(ln *shared.LogNode, db shared.DBNodeList) ReportItem {
	var acc accumulator.Accumulator
	if r.options.Totals {
		acc = accumulator.NewAccumulator()
	}
	re := make([]ReportElement, len(ln.Elements))
	for i := range ln.Elements {
		re[i].Name = ln.Elements[i].Name
		re[i].Val = ln.Elements[i].Val
		if repl, found := db[re[i].Name]; found {
			for _, repl := range repl.Elements {
				res := repl.Val * re[i].Val
				re[i].Ingredients.Add(repl.Name, res)
				if r.options.Totals {
					acc.Add(repl.Name, res)
				}
			}
		} else {
			re[i].Ingredients.Add(re[i].Name, re[i].Val)
			if r.options.Totals {
				acc.Add(ln.Elements[i].Name, ln.Elements[i].Val)
			}
		}
	}
	var totals *[]Total
	if r.options.TotalsOnly {
		re = nil
	}
	if r.options.Totals {
		totals = newTotalFromAccumulator(acc)
	}
	return ReportItem{ln.Time, &re, totals}
}
