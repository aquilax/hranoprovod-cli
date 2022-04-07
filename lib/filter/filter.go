package filter

import (
	"time"

	"github.com/aquilax/hranoprovod-cli/v3/lib/shared"
)

// Config contains node filtering configuration
type Config struct {
	BeginningTime *time.Time
	EndTime       *time.Time
}

func NewDefaultConfig() Config {
	return Config{}
}

// LogNodeFilter is a filter callback function that filters nodes based on the filter config
type LogNodeFilter = func(t time.Time, node *shared.ParserNode) (bool, error)

// GetIntervalNodeFilter creates a node filter callback function given filter config
func GetIntervalNodeFilter(fc Config) *LogNodeFilter {
	if fc.BeginningTime == nil && fc.EndTime == nil {
		// no filter if beginning and end time are nil
		return nil
	}

	inInterval := func(t time.Time) bool {
		if (fc.BeginningTime != nil && !isGoodDate(t, *fc.BeginningTime, dateBeginning)) || (fc.EndTime != nil && !isGoodDate(t, *fc.EndTime, dateEnd)) {
			return false
		}
		return true
	}

	filter := func(t time.Time, node *shared.ParserNode) (bool, error) {
		return inInterval(t), nil
	}
	return &filter
}

// compareType identifies the type of date comparison
type compareType bool

const (
	dateBeginning compareType = true
	dateEnd       compareType = false
)

func isGoodDate(time, compareTime time.Time, ct compareType) bool {
	if time.Equal(compareTime) {
		return true
	}
	if ct == dateBeginning {
		return time.After(compareTime)
	}
	return time.Before(compareTime)
}
