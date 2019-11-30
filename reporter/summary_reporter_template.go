package reporter

import (
	"fmt"
	"io"
	"text/template"
	"time"

	"github.com/aquilax/hranoprovod-cli/accumulator"
	"github.com/aquilax/hranoprovod-cli/shared"
)

const summaryTemplate = `{{printDate .Time}} :
{{- if .Elements }}
{{- range $el := .Elements}}
{{ cNum $el.Val }} : {{ $el.Name }}
{{- end}}
{{- end}}
{{- if .Totals }}
----------------------------------------------------
{{- range $total := .Totals }}
{{ cNum $total.Positive }} : {{ $total.Name }}
{{- end}}
{{- end}}
`

type SummaryReporterTemplate struct {
	options  *Options
	db       shared.DBNodeList
	output   io.Writer
	template *template.Template
}

func NewSummaryReporterTemplate(options *Options, db shared.DBNodeList, writer io.Writer) *SummaryReporterTemplate {
	return &SummaryReporterTemplate{
		options,
		db,
		writer,
		template.Must(template.New("summary").Funcs(template.FuncMap{
			"printDate": func(ts time.Time) string {
				return ts.Format(options.DateFormat)
			},
			"cNum": func(num float64) string {
				return fmt.Sprintf("%10.2f", num)
			},
		}).Parse(summaryTemplate)),
	}
}

func (r *SummaryReporterTemplate) Process(ln *shared.LogNode) error {
	return r.template.Execute(r.output, r.getReportItem(ln, r.db))
}

func (r *SummaryReporterTemplate) Flush() error {
	return nil
}

func (r *SummaryReporterTemplate) getReportItem(ln *shared.LogNode, db shared.DBNodeList) ReportItem {
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
