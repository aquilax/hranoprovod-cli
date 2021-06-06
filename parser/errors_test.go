package parser

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBreakingError(t *testing.T) {
	t.Run("Given errors", func(t *testing.T) {
		t.Run("IO Error works", func(t *testing.T) {
			err := NewErrorIO(errors.New("test"), "file_name")
			assert.Equal(t, "file_name", err.FileName)
			assert.Equal(t, "test", err.Error())
		})
		t.Run("Bad Syntax error works", func(t *testing.T) {
			err := NewErrorBadSyntax(3, "test line")
			assert.Equal(t, "Bad syntax on line 3, \"test line\".", err.Error())
			assert.Equal(t, 3, err.LineNumber)
			assert.Equal(t, "test line", err.Line)
		})
		t.Run("Conversion error works", func(t *testing.T) {
			err := NewErrorConversion("bibip", 5, "line string")
			assert.Equal(t, "Error converting \"bibip\" to float on line 5 \"line string\".", err.Error())
			assert.Equal(t, "bibip", err.Text)
			assert.Equal(t, 5, err.LineNumber)
			assert.Equal(t, "line string", err.Line)
		})
	})
}
