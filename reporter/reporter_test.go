package reporter

import (
	"bytes"
	"github.com/aquilax/hranoprovod-cli/shared"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestReporter(t *testing.T) {
	Convey("Given reporter", t, func() {
		var b bytes.Buffer
		rp := NewReporter(NewDefaultOptions(), &b)
		Convey("Prints list of API results", func() {
			nl := shared.APINodeList{
				shared.APINode{
					Name: "test1",
				},
				shared.APINode{
					Name: "test2",
				},
			}
			expected := `test1:
  calories: 0.000
  fat: 0.000
  carbohydrate: 0.000
  protein: 0.000
test2:
  calories: 0.000
  fat: 0.000
  carbohydrate: 0.000
  protein: 0.000
`
			err := rp.PrintAPISearchResult(nl)
			So(err, ShouldBeNil)
			So(b.String(), ShouldEqual, expected)
		})
	})
}
