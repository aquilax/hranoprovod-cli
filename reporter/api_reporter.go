package reporter

import (
	"fmt"
	"io"

	"github.com/aquilax/hranoprovod-cli/v2/shared"
)

// APIReporter outputs api search results
type APIReporter struct {
	config Config
	output io.Writer
}

// NewAPIReporter creates new API result reporter
func NewAPIReporter(rc Config, writer io.Writer) *APIReporter {
	return &APIReporter{
		rc,
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
	fmt.Fprintf(r.output, "  %s: %0.3f\n", r.config.CaloriesLabel, n.Calories)
	fmt.Fprintf(r.output, "  %s: %0.3f\n", r.config.FatLabel, n.Fat)
	fmt.Fprintf(r.output, "  %s: %0.3f\n", r.config.CarbohydrateLabel, n.Carbohydrate)
	fmt.Fprintf(r.output, "  %s: %0.3f\n", r.config.ProteinLabel, n.Protein)
	return nil
}
