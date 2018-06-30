package parser

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBreakingError(t *testing.T) {
	Convey("Given errors", t, func() {
		Convey("IO Error works", func() {
			err := NewErrorIO(errors.New("test"), "file_name")
			So(err.FileName, ShouldEqual, "file_name")
			So(err.Error(), ShouldEqual, "test")
		})
		Convey("Bad Syntax error works", func() {
			err := NewErrorBadSyntax(3, "test line")
			So(err.Error(), ShouldEqual, "Bad syntax on line 3, \"test line\".")
			So(err.LineNumber, ShouldEqual, 3)
			So(err.Line, ShouldEqual, "test line")
		})
		Convey("Conversion error works", func() {
			err := NewErrorConversion("bibip", 5, "line string")
			So(err.Error(), ShouldEqual, "Error converting \"bibip\" to float on line 5 \"line string\".")
			So(err.Text, ShouldEqual, "bibip")
			So(err.LineNumber, ShouldEqual, 5)
			So(err.Line, ShouldEqual, "line string")
		})
	})
}
