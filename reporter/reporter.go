package reporter

import (
	"io"
	"os"

	"github.com/aquilax/hranoprovod-cli/v2/shared"
)

// Config contains the options for the reporter
type Config struct {
	CSV                  bool
	Color                bool
	TotalsOnly           bool
	Totals               bool
	DateFormat           string
	CaloriesLabel        string
	FatLabel             string
	CarbohydrateLabel    string
	ProteinLabel         string
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
	Output               io.Writer
}

// NewDefaultConfig returns the default reporter config
func NewDefaultConfig() Config {
	return Config{
		CSV:                  false,
		Color:                false,
		DateFormat:           "2006/01/02",
		CaloriesLabel:        "calories",
		FatLabel:             "fat",
		CarbohydrateLabel:    "carbohydrate",
		ProteinLabel:         "protein",
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
func NewRegReporter(c Config, db shared.DBNodeList, writer io.Writer) Reporter {
	if c.Unresolved {
		return NewUnsolvedReporter(c, db, writer)
	}
	if len(c.SingleElement) > 0 {
		if c.ElementGroupByFood {
			return newElementByFoodReporter(c, db, writer)
		}
		return newSingleReporter(c, db, writer)
	}
	if len(c.SingleFood) > 0 {
		return newSingleFoodReporter(c, db, writer)
	}
	if c.UseOldRegReporter {
		return newRegReporter(c, db, writer)
	}
	return newRegReporterTemplate(c, db, writer)
}

// NewBalanceReporter returns balance reporter
func NewBalanceReporter(options Config, db shared.DBNodeList, writer io.Writer) Reporter {
	if len(options.SingleElement) > 0 {
		return newBalanceSingleReporter(options, db, writer)
	}
	return newBalanceReporter(options, db, writer)
}
