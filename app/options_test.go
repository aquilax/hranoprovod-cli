package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	t.Run("Given options", func(t *testing.T) {
		o := NewOptions()
		t.Run("New options is created", func(t *testing.T) {
			assert.NotNil(t, o)
			assert.NotNil(t, o.ReporterConfig)
			assert.True(t, o.ReporterConfig.Color)
			assert.NotNil(t, o.ParserConfig)
		})
	})
	t.Run("Given fileExists", func(t *testing.T) {
		t.Run("Returns false if file does not exit", func(t *testing.T) {
			ex, err := fileExists("file_does_not_exist")
			assert.False(t, ex)
			assert.Nil(t, err)
		})
	})
}
