package reporter

import (
	"io"
	"time"

	"github.com/aquilax/hranoprovod-cli/shared"
)

type ReportType int8

const (
	Reg ReportType = iota
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
	}
}

type Reporter interface {
	Process(ln *shared.LogNode) error
}

func NewReporter(rt ReportType, options *Options, db *shared.NodeList, writer io.Writer) Reporter {
	if rt == Bal {
		return NewBalanceReporter(options, db, writer)
	}
	if options.Unresolved {
		return NewUnsolvedReporter(options, db, writer)
	}
	if len(options.SingleElement) > 0 {
		return NewSingleReporter(options, db, writer)
	}
	if len(options.SingleFood) > 0 {
		return NewSingleFoodReporter(options, db, writer)
	}
	return NewRegReporter(options, db, writer)
}
