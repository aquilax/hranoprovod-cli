package accumulator

const (
	// Negative register
	Negative = 0
	// Positive register
	Positive = 1
)

// AccValues contains accumulator values
type AccValues []float64

// Accumulator accumulates element values by name
type Accumulator map[string]AccValues

// NewAccumulator returns new accumulator
func NewAccumulator() Accumulator {
	return Accumulator{}
}

// Add adds name/value to the accumulator
func (acc Accumulator) Add(name string, val float64) {
	sign := Positive
	if val < 0 {
		sign = Negative
	}
	if _, exists := acc[name]; exists {
		acc[name][sign] += val
	} else {
		newVal := make(AccValues, 2)
		newVal[sign] = val
		acc[name] = newVal
	}
}
