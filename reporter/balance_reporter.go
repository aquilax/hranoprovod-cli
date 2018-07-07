package reporter

import (
	"io"

	"github.com/aquilax/hranoprovod-cli/shared"
)

type BalanceReporter struct {
	options *Options
	db      *shared.NodeList
	output  io.Writer
}

func NewBalanceReporter(options *Options, db *shared.NodeList, writer io.Writer) *BalanceReporter {
	return &BalanceReporter{
		options,
		db,
		writer,
	}
}

func (r *BalanceReporter) Process(ln *shared.LogNode) error {
	return nil
}

func (r *BalanceReporter) Flush() error {
	return nil
}
