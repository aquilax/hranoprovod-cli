package reporter

import (
	"io"
	"text/template"

	"github.com/aquilax/hranoprovod-cli/v2/shared"
)

const summaryTemplate = `{{formatDate .Time}} :
{{- if .Totals }}
{{- range $total := .Totals }}
{{ formatValue $total.Positive }} : {{ $total.Name }}
{{- end}}
{{- end}}
------------
{{- if .Elements }}
{{- range $el := .Elements}}
{{ formatValue $el.Val }} : {{ $el.Name }}
{{- end}}
{{- end}}
`

// SummaryReporterTemplate is a summary reporter
type SummaryReporterTemplate struct {
	options  Config
	db       shared.DBNodeList
	output   io.Writer
	template *template.Template
}

// NewSummaryReporterTemplate creates new summary reporter
func NewSummaryReporterTemplate(options Config, db shared.DBNodeList, writer io.Writer) *SummaryReporterTemplate {
	return &SummaryReporterTemplate{
		options,
		db,
		writer,
		template.Must(template.New("summary").Funcs(getTemplateFunctions(options)).Parse(summaryTemplate)),
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
