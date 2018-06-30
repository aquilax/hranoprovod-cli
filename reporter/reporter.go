package reporter

import (
	"fmt"
	"github.com/aquilax/hranoprovod-cli/shared"
	"io"
	"strings"
	"time"
)

const (
	reset       = "\x1B[0m"
	bold        = "\x1B[1m"
	dim         = "\x1B[2m"
	under       = "\x1B[4m"
	reverse     = "\x1B[7m"
	hide        = "\x1B[8m"
	clearscreen = "\x1B[2J"
	clearline   = "\x1B[2K"
	black       = "\x1B[30m"
	red         = "\x1B[31m"
	green       = "\x1B[32m"
	yellow      = "\x1B[33m"
	blue        = "\x1B[34m"
	magenta     = "\x1B[35m"
	cyan        = "\x1B[36m"
	white       = "\x1B[37m"
	bblack      = "\x1B[40m"
	bred        = "\x1B[41m"
	bgreen      = "\x1B[42m"
	byellow     = "\x1B[43m"
	bblue       = "\x1B[44m"
	bmagenta    = "\x1B[45m"
	bcyan       = "\x1B[46m"
	bwhite      = "\x1B[47m"
	newline     = "\r\n\x1B[0m"
)

// Options contains the options for the reporter
type Options struct {
	CSV               bool
	Color             bool
	DateFormat        string
	CaloriesLabel     string
	FatLabel          string
	CarbohydrateLabel string
	ProteinLabel      string
}

// Reporter is the main report structure
type Reporter struct {
	options *Options
	output  io.Writer
}

// DefaultOptions returns the default reporter options
func NewDefaultOptions() *Options {
	return &Options{
		CSV:               false,
		Color:             false,
		DateFormat:        "2006/01/02",
		CaloriesLabel:     "calories",
		FatLabel:          "fat",
		CarbohydrateLabel: "carbohydrate",
		ProteinLabel:      "protein",
	}
}

// NewReporter creates new reporter
func NewReporter(ro *Options, writer io.Writer) *Reporter {
	return &Reporter{
		ro,
		writer,
	}
}

// PrintAPISearchResult prints a list of search resilts
func (r *Reporter) PrintAPISearchResult(nl shared.APINodeList) error {
	for _, n := range nl {
		err := r.PrintAPINode(n)
		if err != nil {
			return err
		}
	}
	return nil
}

// PrintAPINode prints single API result
func (r *Reporter) PrintAPINode(n shared.APINode) error {
	fmt.Fprintln(r.output, n.Name+":")
	fmt.Fprintf(r.output, "  %s: %0.3f\n", r.options.CaloriesLabel, n.Calories)
	fmt.Fprintf(r.output, "  %s: %0.3f\n", r.options.FatLabel, n.Fat)
	fmt.Fprintf(r.output, "  %s: %0.3f\n", r.options.CarbohydrateLabel, n.Carbohydrate)
	fmt.Fprintf(r.output, "  %s: %0.3f\n", r.options.ProteinLabel, n.Protein)
	return nil
}

func (r *Reporter) cNum(num float32) string {
	if r.options.Color {
		if num > 0 {
			return red + fmt.Sprintf("%10.2f", num) + reset
		}
		if num < 0 {
			return green + fmt.Sprintf("%10.2f", num) + reset
		}
	}
	return fmt.Sprintf("%10.2f", num)
}

func (r *Reporter) PrintDate(ts time.Time) {
	fmt.Fprintf(r.output, "%s\n", ts.Format(r.options.DateFormat))
}

func (r *Reporter) PrintElement(element *shared.Element) {
	fmt.Fprintf(r.output, "\t%-27s :%s\n", element.Name, r.cNum(element.Val))
}

func (r *Reporter) PrintIngredient(name string, value float32) {
	fmt.Fprintf(r.output, "\t\t%20s %s\n", name, r.cNum(value))
}

func (r *Reporter) PrintTotalHeader() {
	fmt.Fprintf(r.output, "\t-- %s %s\n", "TOTAL ", strings.Repeat("-", 52))
}

func (r *Reporter) PrintTotalRow(name string, pos float32, neg float32) {
	fmt.Fprintf(r.output, "\t\t%20s %s %s =%s\n", name, r.cNum(pos), r.cNum(neg), r.cNum(pos+neg))
}

func (r *Reporter) PrintSingleElementRow(ts time.Time, name string, pos float32, neg float32) {
	format := "%s %20s %10.2f %10.2f =%10.2f\n"
	if r.options.CSV {
		format = "%s;\"%s\";%0.2f;%0.2f;%0.2f\n"
	}
	fmt.Fprintf(r.output, format, ts.Format(r.options.DateFormat), name, pos, -1*neg, pos+neg)
}

func (r *Reporter) PrintSingleFoodRow(ts time.Time, name string, val float32) {
	fmt.Fprintf(r.output, "%s\t%s\t%0.2f\n", ts.Format(r.options.DateFormat), name, val)
}

func (r *Reporter) PrintUnresolvedRow(name string) {
	fmt.Fprintln(r.output, name)
}
