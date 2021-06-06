package accumulator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccumulator(t *testing.T) {
	t.Run("Given the accumulator", func(t *testing.T) {
		acc := NewAccumulator()
		t.Run("When a positive element is added", func(t *testing.T) {
			acc.Add("test", 1.22)
			t.Run("It should go to the positive accumulator", func(t *testing.T) {
				assert.Equal(t, 1.22, acc["test"][Positive])
			})
			t.Run("When a positive value is added to the same key", func(t *testing.T) {
				acc.Add("test", 2.33)
				t.Run("It is accumulated in the positive register", func(t *testing.T) {
					assert.Equal(t, 3.55, acc["test"][Positive])
				})
			})
			t.Run("When a negative value is added to the same key", func(t *testing.T) {
				acc.Add("test", -2.33)
				t.Run("It is accumulated in the positive register", func(t *testing.T) {
					assert.Equal(t, -2.33, acc["test"][Negative])
					assert.Equal(t, 3.55, acc["test"][Positive])
				})
			})
			t.Run("When negative element is added", func(t *testing.T) {
				acc.Add("test2", -1.32)
				t.Run("It should go to the negative accumulator", func(t *testing.T) {
					assert.Equal(t, -1.32, acc["test2"][Negative])
				})
			})
		})
	})
}
