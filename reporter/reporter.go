package reporter

import (
	"io"
	"time"

	"github.com/aquilax/hranoprovod-cli/shared"
)

// ReportType represents the type of the report
type ReportType int8

const (
	// Reg is register report type
	Reg ReportType = iota
	// Bal is balance report type
	Bal
)

// Options contains the options for the reporter
type Options struct {
	CSV               bool
	Color             bool
	TotalsOnly        bool
	Totals            bool
	DateFormat        string
	CaloriesLabel     string
	FatLabel          string
	CarbohydrateLabel string
	ProteinLabel      string
	HasBeginning      bool
	HasEnd            bool
	BeginningTime     time.Time
	EndTime           time.Time
	Unresolved        bool
	SingleElement     string
	SingleFood        string
	CollapseLast      bool
	Collapse          bool
}

// NewDefaultOptions returns the default reporter options
func NewDefaultOptions() *Options {
	return &Options{
		CSV:               false,
		Color:             false,
		DateFormat:        "2006/01/02",
		CaloriesLabel:     "calories",
		FatLabel:          "fat",
		CarbohydrateLabel: "carbohydrate",
		ProteinLabel:      "protein",
		Totals:            true,
		BeginningTime:     time.Now(),
		EndTime:           time.Now(),
		CollapseLast:      false,
		Collapse:          false,
	}
}

// Reporter is the reporting interface
type Reporter interface {
	Process(ln *shared.LogNode) error
	Flush() error
}

// NewReporter creates new response handler
func NewReporter(rt ReportType, options *Options, db *shared.NodeList, writer io.Writer) Reporter {
	if rt == Bal {
		if len(options.SingleElement) > 0 {
			return newBalanceSingleReporter(options, db, writer)
		}
		return newBalanceReporter(options, db, writer)
	}
	if options.Unresolved {
		return newUnsolvedReporter(options, db, writer)
	}
	if len(options.SingleElement) > 0 {
		return newSingleReporter(options, db, writer)
	}
	if len(options.SingleFood) > 0 {
		return newSingleFoodReporter(options, db, writer)
	}
	return newRegReporter(options, db, writer)
}
