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
	return r.template.Execute(r.output, getReportItem(ln, r.db, r.options))
}

func (r *regReporterTemplate) Flush() error {
	return nil
}

func newTotalFromAccumulator(acc accumulator.Accumulator) *[]total {
	var result = make([]total, len(acc))
	var ss = make(sort.StringSlice, len(acc))
	i := 0
	for name := range acc {
		ss[i] = name
		i++
	}
	sort.Sort(ss)
	for i, name := range ss {
		result[i] = total{name, acc[name][accumulator.Positive], acc[name][accumulator.Negative], acc[name][accumulator.Positive] + acc[name][accumulator.Negative]}
	}
	return &result
}
