package reporter

import (
	"fmt"
	"io"

	"github.com/aquilax/hranoprovod-cli/shared"
)

// APIReporter outputs api search results
type APIReporter struct {
	options *Options
	output  io.Writer
}

// NewAPIReporter creates new API result reporter
func NewAPIReporter(ro *Options, writer io.Writer) *APIReporter {
	return &APIReporter{
		ro,
		writer,
	}
}

// PrintAPISearchResult prints a list of search results
func (r *APIReporter) PrintAPISearchResult(nl shared.APINodeList) error {
	for _, n := range nl {
		err := r.PrintAPINode(n)
		if err != nil {
			return err
		}
	}
	return nil
}

// PrintAPINode prints single API result
func (r *APIReporter) PrintAPINode(n shared.APINode) error {
	fmt.Fprintln(r.output, n.Name+":")
	fmt.Fprintf(r.output, "  %s: %0.3f\n", r.options.CaloriesLabel, n.Calories)
	fmt.Fprintf(r.output, "  %s: %0.3f\n", r.options.FatLabel, n.Fat)
	fmt.Fprintf(r.output, "  %s: %0.3f\n", r.options.CarbohydrateLabel, n.Carbohydrate)
	fmt.Fprintf(r.output, "  %s: %0.3f\n", r.options.ProteinLabel, n.Protein)
	return nil
}
