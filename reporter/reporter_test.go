package reporter

import (
	"bytes"
	"testing"
	"time"

	"github.com/aquilax/hranoprovod-cli/v2/shared"
	"github.com/tj/assert"
)

func TestReporter(t *testing.T) {
	t.Run("Given reporter", func(t *testing.T) {
		var b bytes.Buffer
		nl := shared.NewDBNodeList()
		o := NewDefaultOptions()
		o.Unresolved = true
		rp := NewRegReporter(o, nl, &b)
		t.Run("Prints list of unresolved items", func(t *testing.T) {
			el := shared.NewElements()
			el.Add("test", 3.55)
			ln := shared.NewLogNode(time.Now(), el)
			expected := `test
`
			err := rp.Process(ln)
			rp.Flush()
			assert.Nil(t, err)
			assert.Equal(t, expected, b.String())
		})
	})
}
