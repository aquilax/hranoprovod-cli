package reporter

import (
	"io"
	"time"

	"github.com/aquilax/hranoprovod-cli/shared"
)

// Options contains the options for the reporter
type Options struct {
	CSV                bool
	Color              bool
	TotalsOnly         bool
	Totals             bool
	DateFormat         string
	CaloriesLabel      string
	FatLabel           string
	CarbohydrateLabel  string
	ProteinLabel       string
	HasBeginning       bool
	HasEnd             bool
	BeginningTime      time.Time
	EndTime            time.Time
	Unresolved         bool
	SingleElement      string
	SingleFood         string
	CollapseLast       bool
	Collapse           bool
	ElementGroupByFood bool
	ShortenStrings     bool
	UseNewRegReporter  bool
}

// NewDefaultOptions returns the default reporter options
func NewDefaultOptions() *Options {
	return &Options{
		CSV:                false,
		Color:              false,
		DateFormat:         "2006/01/02",
		CaloriesLabel:      "calories",
		FatLabel:           "fat",
		CarbohydrateLabel:  "carbohydrate",
		ProteinLabel:       "protein",
		Totals:             true,
		BeginningTime:      time.Now(),
		EndTime:            time.Now(),
		CollapseLast:       false,
		Collapse:           false,
		ElementGroupByFood: false,
		ShortenStrings:     false,
		UseNewRegReporter:  false,
	}
}

// Reporter is the reporting interface
type Reporter interface {
	Process(ln *shared.LogNode) error
	Flush() error
}

// NewRegReporter creates new response handler
func NewRegReporter(options *Options, db shared.DBNodeList, writer io.Writer) Reporter {
	if options.Unresolved {
		return NewUnsolvedReporter(options, db, writer)
	}
	if len(options.SingleElement) > 0 {
		if options.ElementGroupByFood {
			return newElementByFoodReporter(options, db, writer)
		}
		return newSingleReporter(options, db, writer)
	}
	if len(options.SingleFood) > 0 {
		return newSingleFoodReporter(options, db, writer)
	}
	if options.UseNewRegReporter {
		return newRegReporterTemplate(options, db, writer)
	}
	return newRegReporter(options, db, writer)
}

// NewBalanceReporter returns balance reporter
func NewBalanceReporter(options *Options, db shared.DBNodeList, writer io.Writer) Reporter {
	if len(options.SingleElement) > 0 {
		return newBalanceSingleReporter(options, db, writer)
	}
	return newBalanceReporter(options, db, writer)
}
