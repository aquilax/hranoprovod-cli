package reporter

import (
	"io"
	"os"

	shared "github.com/aquilax/hranoprovod-cli/v3"
)

type CommonConfig struct {
	Output io.Writer
	Color  bool
}

func NewCommonConfig(output io.Writer, color bool) CommonConfig {
	return CommonConfig{
		Output: output,
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
