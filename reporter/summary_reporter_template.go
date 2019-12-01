package reporter

import (
	"fmt"
	"io"
	"text/template"
	"time"

	"github.com/aquilax/hranoprovod-cli/shared"
)

const summaryTemplate = `{{printDate .Time}} :
{{- if .Totals }}
{{- range $total := .Totals }}
{{ cNum $total.Positive }} : {{ $total.Name }}
{{- end}}
{{- end}}
------------
{{- if .Elements }}
{{- range $el := .Elements}}
{{ cNum $el.Val }} : {{ $el.Name }}
{{- end}}
{{- end}}
`

// SummaryReporterTemplate is a summary reporter
type SummaryReporterTemplate struct {
	options  *Options
	db       shared.DBNodeList
	output   io.Writer
	template *template.Template
}

// NewSummaryReporterTemplate creates new summary reporter
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

// Process process shared node
func (r *SummaryReporterTemplate) Process(ln *shared.LogNode) error {
	return r.template.Execute(r.output, getReportItem(ln, r.db, r.options))
}

// Flush does nothing
func (r *SummaryReporterTemplate) Flush() error {
	return nil
}
