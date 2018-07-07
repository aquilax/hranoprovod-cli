package reporter

import (
	"bytes"
	"testing"
	"time"

	"github.com/aquilax/hranoprovod-cli/shared"
	. "github.com/smartystreets/goconvey/convey"
)

func TestReporter(t *testing.T) {
	Convey("Given reporter", t, func() {
		var b bytes.Buffer
		nl := &shared.NodeList{}
		o := NewDefaultOptions()
		o.Unresolved = true
		rp := NewReporter(Reg, o, nl, &b)
		Convey("Prints list of unresolved items", func() {
			el := shared.NewElements()
			el.Add("test", 3.55)
			ln := shared.NewLogNode(time.Now(), el)
			expected := `test
`
			err := rp.Process(ln)
			rp.Flush()
			So(err, ShouldBeNil)
			So(b.String(), ShouldEqual, expected)
		})
	})
}
