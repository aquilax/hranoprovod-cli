package reporter

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReporter(t *testing.T) {
	t.Run("Given reporter", func(t *testing.T) {
		var b bytes.Buffer
		nl := hranoprovod.NewDBNodeMap()
		c := NewDefaultConfig()
		c.Unresolved = true
		c.Output = &b
		rp := NewRegReporter(c, nl)
		t.Run("Prints list of unresolved items", func(t *testing.T) {
			el := hranoprovod.NewElements()
			el.Add("test", 3.55)
			ln := hranoprovod.NewLogNode(time.Now(), el, nil)
			expected := `test
`
			err := rp.Process(ln)
			rp.Flush()
			assert.Nil(t, err)
			assert.Equal(t, expected, b.String())
		})
	})
}
