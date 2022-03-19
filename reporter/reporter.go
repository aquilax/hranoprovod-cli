package reporter

import (
	"io"
	"os"

	"github.com/aquilax/hranoprovod-cli/v2/shared"
)

type CommonConfig struct {
	Output io.Writer
	Color  bool
}

func NewCommonConfig(color bool) CommonConfig {
	return CommonConfig{
		Output: os.Stdout,
		Color:  color,
	}
}

// Config contains the options for the reporter
type Config struct {
	Output               io.Writer
	CSV                  bool
	Color                bool
	TotalsOnly           bool
	Totals               bool
	DateFormat           string
	Unresolved           bool
	SingleElement        string
	SingleFood           string
	CollapseLast         bool
	Collapse             bool
	ElementGroupByFood   bool
	ShortenStrings       bool
	UseOldRegReporter    bool
	InternalTemplateName string
	CSVSeparator         rune
}

// NewDefaultConfig returns the default reporter config
func NewDefaultConfig() Config {
	return Config{
		CSV:                  false,
		Color:                true,
		DateFormat:           "2006/01/02",
		Totals:               true,
		CollapseLast:         false,
		Collapse:             false,
		ElementGroupByFood:   false,
		ShortenStrings:       false,
		UseOldRegReporter:    false,
		InternalTemplateName: "default",
		CSVSeparator:         ',',
		Output:               os.Stdout,
	}
}

// Reporter is the reporting interface
type Reporter interface {
	Process(ln *shared.LogNode) error
	Flush() error
}

// NewRegReporter creates new response handler
func NewRegReporter(c Config, db shared.DBNodeMap) Reporter {
	if c.Unresolved {
		return NewUnsolvedReporter(c, db)
	}
	if len(c.SingleElement) > 0 {
		if c.ElementGroupByFood {
			return newElementByFoodReporter(c, db)
		}
		return newSingleReporter(c, db)
	}
	if len(c.SingleFood) > 0 {
		return newSingleFoodReporter(c, db)
	}
	if c.UseOldRegReporter {
		return newRegReporter(c, db)
	}
	return newRegReporterTemplate(c, db)
}

// NewBalanceReporter returns balance reporter
func NewBalanceReporter(options Config, db shared.DBNodeMap) Reporter {
	if len(options.SingleElement) > 0 {
		return newBalanceSingleReporter(options, db)
	}
	return newBalanceReporter(options, db)
}
