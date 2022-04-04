package register

import (
	"bytes"
	"testing"
	"time"

	"github.com/aquilax/hranoprovod-cli/v2/cmd/hranoprovod-cli/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/v2/lib/shared"
	"github.com/stretchr/testify/assert"
)

const expected = `2019/10/10
	test1                       :      3.10
		               test1       3.10
	test1                       :      3.10
		               test1       3.10
	test2                       :      3.20
		                 el1       3.52
		                 el2       3.84
		                 el3       4.16
	test3/test3/t…3/test3/test3 :      3.30
		test3/tes…est3/test3       3.30
	-- TOTAL  ----------------------------------------------------
		                 el1       3.52       0.00 =      3.52
		                 el2       3.84       0.00 =      3.84
		                 el3       4.16       0.00 =      4.16
		               test1       6.20       0.00 =      6.20
		test3/tes…est3/test3       3.30       0.00 =      3.30
`

func TestRegReporterTemplate(t *testing.T) {
	t.Run("Given template reg reporter", func(t *testing.T) {
		var b bytes.Buffer
		db := shared.NewDBNodeMap()
		dbElements := shared.NewElements()
		dbElements.Add("el1", 1.1)
		dbElements.Add("el2", 1.2)
		dbElements.Add("el3", 1.3)
		dbNode := shared.DBNode{Header: "test2", Elements: dbElements}
		db.Push(&dbNode)

		t.Run("Prints list of unresolved items", func(t *testing.T) {
			c := reporter.NewDefaultConfig()
			c.ShortenStrings = true
			c.Color = false
			c.Output = &b
			rp := NewRegReporter(c, db)

			el := shared.NewElements()
			el.Add("test1", 3.1)
			el.Add("test1", 3.1)
			el.Add("test2", 3.2)
			el.Add("test3/test3/test3/test3/test3/test3/test3/test3/test3", 3.3)
			ln := shared.NewLogNode(time.Date(2019, 10, 10, 0, 0, 0, 0, time.UTC), el, nil)
			err := rp.Process(ln)
			rp.Flush()
			assert.Nil(t, err)
			assert.Equal(t, string(expected), b.String())
		})
	})
}
