package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestOptions(t *testing.T) {
	Convey("Given options", t, func() {
		o := NewOptions()
		Convey("New options is created", func() {
			So(o, ShouldNotBeNil)
			So(o.Reporter, ShouldNotBeNil)
			So(o.Reporter.Color, ShouldBeTrue)
			So(o.Processor, ShouldNotBeNil)
			So(o.Parser, ShouldNotBeNil)
			So(o.API, ShouldNotBeNil)
		})
	})
	Convey("Given fileExists", t, func() {
		Convey("Returns false if file does not exit", func() {
			ex, err := fileExists("ASDDD!@!@!@");
			So(ex, ShouldBeFalse)
			So(err, ShouldBeNil)
		})
	})
}
