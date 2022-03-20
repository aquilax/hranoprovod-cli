package reporter

import (
	"bufio"
	"sort"
	"text/template"

	"github.com/aquilax/hranoprovod-cli/v2"
	"github.com/aquilax/hranoprovod-cli/v2/accumulator"
)

const defaultTemplate = `{{formatDate .Time}}
{{- if .Elements }}
{{- range $el := .Elements}}
{{ printf "\t%-27s :%s" (shorten $el.Name 27) (formatValue $el.Value) }}
{{- range $ing := $el.Ingredients}}
{{ printf "\t\t%20s %s" (shorten $ing.Name 20) (formatValue $ing.Value) }}
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
{{ printf "  %s  %s" (formatValue $el.Value) $el.Name}}
{{- range $ing := $el.Ingredients}}
{{ printf "  %s    %s" (formatValue $ing.Value) $ing.Name }}
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
	config   Config
	db       hranoprovod.DBNodeMap
	output   *bufio.Writer
	template *template.Template
}

func getInternalTemplate(internalTemplateName string) string {
	if internalTemplateName == "left-aligned" {
		return leftAlignedTemplate
	}
	return defaultTemplate
}

func newRegReporterTemplate(config Config, db hranoprovod.DBNodeMap) *regReporterTemplate {
	return &regReporterTemplate{
		config,
		db,
		bufio.NewWriter(config.Output),
		template.Must(template.New("template").Funcs(getTemplateFunctions(config)).Parse(getInternalTemplate(config.InternalTemplateName))),
	}
}

func (r *regReporterTemplate) Process(ln *hranoprovod.LogNode) error {
	return r.template.Execute(r.output, getReportItem(ln, r.db, r.config))
}

func (r *regReporterTemplate) Flush() error {
	return r.output.Flush()
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
