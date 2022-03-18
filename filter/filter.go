package filter

import "time"

type Config struct {
	BeginningTime *time.Time
	EndTime       *time.Time
}
