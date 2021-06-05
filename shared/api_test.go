package shared

import (
	"testing"

	"github.com/tj/assert"
)

func TestAccumulator(t *testing.T) {
	t.Run("Given APIError", func(t *testing.T) {
		err := APIError{
			IsError: true,
			Code:    100,
			Status:  "status",
			Message: "message",
		}
		assert.Equal(t, "status: message", err.Error())
	})
}
