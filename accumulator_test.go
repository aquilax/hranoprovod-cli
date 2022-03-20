package hranoprovod

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccumulator_Add(t *testing.T) {
	tests := []struct {
		name   string
		caller func() Accumulator
		want   Accumulator
	}{
		{
			"add positive number to the right register",
			func() Accumulator {
				acc := NewAccumulator()
				acc.Add("test", 2.22)
				return acc
			},
			Accumulator{
				"test": AccValues{0, 2.22},
			},
		},
		{
			"accumulates positive number to the right register",
			func() Accumulator {
				acc := NewAccumulator()
				acc.Add("test", 2.00)
				acc.Add("test", 3.00)
				return acc
			},
			Accumulator{
				"test": AccValues{0, 5.00},
			},
		},
		{
			"adds negative number to the right register",
			func() Accumulator {
				acc := NewAccumulator()
				acc.Add("test", -3.00)
				acc.Add("test", 2.00)
				return acc
			},
			Accumulator{
				"test": AccValues{-3.00, 2.00},
			},
		},
		{
			"accumulates correctly to both registers",
			func() Accumulator {
				acc := NewAccumulator()
				acc.Add("test", 1.00)
				acc.Add("test", -1.00)
				acc.Add("test2", 2.00)
				acc.Add("test2", 2.00)
				acc.Add("test2", -2.00)
				acc.Add("test2", -2.00)
				acc.Add("test3", 0)
				return acc
			},
			Accumulator{
				"test":  AccValues{-1.00, 1.00},
				"test2": AccValues{-4.00, 4.00},
				"test3": AccValues{0, 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.caller())
		})
	}
}
