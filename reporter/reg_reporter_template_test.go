package reporter

import (
	"bytes"
	"io/ioutil"
	"testing"
	"time"

	"github.com/aquilax/hranoprovod-cli/shared"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRegReporterTemplate(t *testing.T) {
	Convey("Given template reg reporter", t, func() {
		var b bytes.Buffer
		db := shared.NewDBNodeList()
		dbel := shared.NewElements()
		dbel.Add("el1", 1.1)
		dbel.Add("el2", 1.2)
		dbel.Add("el3", 1.3)
		dbNode := shared.DBNode{"test2", dbel}
		db.Push(&dbNode)

		o := NewDefaultOptions()
		o.ShortenStrings = true
		rp := NewRegReporter(o, db, &b)
		Convey("Prints list of unresolved items", func() {
			el := shared.NewElements()
			el.Add("test1", 3.1)
			el.Add("test1", 3.1)
			el.Add("test2", 3.2)
			el.Add("test3/test3/test3/test3/test3/test3/test3/test3/test3", 3.3)
			ln := shared.NewLogNode(time.Date(2019, 10, 10, 0, 0, 0, 0, time.UTC), el)
			expected, _ := ioutil.ReadFile("testdata/TestRegReporterTemplate.txt")
			err := rp.Process(ln)
			rp.Flush()
			So(err, ShouldBeNil)
			So(b.String(), ShouldEqual, string(expected))
		})
	})
}
