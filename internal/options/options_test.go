package options

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	t.Run("Given options", func(t *testing.T) {
		o := New()
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

	t.Run("loadFromConfigFile", func(t *testing.T) {
		o := New()
		r := strings.NewReader(`
[Global]
Now=2020-01-01T01:00:00Z
DbFileName=/tmp/db.yaml
LogFileName=/tmp/log.yaml
DateFormat=2006-01-02
[Resolver]
MaxDepth=10
`)
		err := loadFromConfigFile(o, r)
		assert.Nil(t, err)
		expectedNow, _ := time.Parse(time.RFC3339, "2020-01-01T01:00:00Z")
		assert.Equal(t, o.GlobalConfig.Now, expectedNow)
		assert.Equal(t, o.GlobalConfig.DbFileName, "/tmp/db.yaml")
		assert.Equal(t, o.GlobalConfig.LogFileName, "/tmp/log.yaml")
		assert.Equal(t, o.GlobalConfig.DateFormat, "2006-01-02")

		assert.Equal(t, o.ResolverConfig.MaxDepth, 10)

	})
}
