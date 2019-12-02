package reporter

import (
	"fmt"
	"text/template"
	"time"

	"github.com/aquilax/truncate"
)

const negativeFormat = red + "%10.2f" + reset
const positiveFormat = green + "%10.2f" + reset

func getTemplateFunctions(options *Options) template.FuncMap {
	return template.FuncMap{
		"formatDate": func(ts time.Time) string {
			return ts.Format(options.DateFormat)
		},
		"formatValue": getFormatValue(options),
		"shorten":     getShorten(options),
	}
}

func getShorten(options *Options) func(string, int) string {
	if options.ShortenStrings {
		return func(t string, max int) string {
			return truncate.Truncate(t, max, truncate.DEFAULT_OMISSION, truncate.PositionMiddle)
		}
	}
	return func(t string, max int) string {
		return t
	}
}

func getFormatValue(options *Options) func(float64) string {
	if options.Color {
		return func(num float64) string {
			if num > 0 {
				return fmt.Sprintf(negativeFormat, num)
			}
			if num < 0 {
				return fmt.Sprintf(positiveFormat, num)
			}
			return fmt.Sprintf("%10.2f", num)
		}
	}
	return func(num float64) string {
		return fmt.Sprintf("%10.2f", num)
	}
}
