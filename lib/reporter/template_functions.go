package reporter

import (
	"fmt"
	"text/template"
	"time"

	"github.com/aquilax/truncate"
)

const negativeFormat = red + "%10.2f" + reset
const positiveFormat = green + "%10.2f" + reset

func getTemplateFunctions(config Config) template.FuncMap {
	return template.FuncMap{
		"formatDate": func(ts time.Time) string {
			return ts.Format(config.DateFormat)
		},
		"formatValue": getFormatValue(config.Color),
		"shorten":     getShorten(config.ShortenStrings),
	}
}

func getShorten(shortenStrings bool) func(string, int) string {
	if shortenStrings {
		return func(t string, max int) string {
			return truncate.Truncate(t, max, truncate.DEFAULT_OMISSION, truncate.PositionMiddle)
		}
	}
	return func(t string, max int) string {
		return t
	}
}

func getFormatValue(hasColor bool) func(float64) string {
	if hasColor {
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
