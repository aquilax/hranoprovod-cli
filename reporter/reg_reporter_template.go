package reporter

import (
	"io"
	"sort"
	"text/template"

	"github.com/aquilax/hranoprovod-cli/v2/accumulator"
	"github.com/aquilax/hranoprovod-cli/v2/shared"
)

const defaultTemplate = `{{formatDate .Time}}
{{- if .Elements }}
{{- range $el := .Elements}}
{{ printf "\t%-27s :%s" (shorten $el.Name 27) (formatValue $el.Val) }}
{{- range $ing := $el.Ingredients}}
{{ printf "\t\t%20s %s" (shorten $ing.Name 20) (formatValue $ing.Val) }}
{{- end}}
{{- end}}
{{- end}}
{{- if .Totals }}
	-- TOTAL  ----------------------------------------------------
{{- range $total := .Totals }}
{{ printf "\t\t%20s %s %s =%s" (shorten $total.Name 20) (formatValue $total.Positive) (formatValue $total.Negative) (formatValue $total.Sum) }}
{{- end}}
{{- end}}
`
const leftAlignedTemplate = `{{formatDate .Time}}
{{- if .Elements }}
{{- range $el := .Elements}}
{{ printf "  %s  %s" (formatValue $el.Val) $el.Name}}
{{- range $ing := $el.Ingredients}}
{{ printf "  %s    %s" (formatValue $ing.Val) $ing.Name }}
{{- end}}
{{- end}}
{{- end}}
{{- if .Totals }}
------------------------------------------------------- TOTAL --
{{- range $total := .Totals }}
{{ printf "  %s %s = %s  %s" (formatValue $total.Positive) (formatValue $total.Negative) (formatValue $total.Sum) $total.Name }}
{{- end}}
{{- end}}
`

type regReporterTemplate struct {
	options  *Options
	db       shared.DBNodeList
	output   io.Writer
	template *template.Template
}

func getInternalTemplate(internalTemplateName string) string {
	if internalTemplateName == "left-aligned" {
		return leftAlignedTemplate
	}
	return defaultTemplate
}

func newRegReporterTemplate(options *Options, db shared.DBNodeList, writer io.Writer) *regReporterTemplate {
	return &regReporterTemplate{
		options,
		db,
		writer,
		template.Must(template.New("template").Funcs(getTemplateFunctions(options)).Parse(getInternalTemplate(options.InternalTemplateName))),
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
