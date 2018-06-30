package shared

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestNewLogNode(t *testing.T) {
	Convey("Given NewLogNode", t, func() {
		now := time.Now()
		elements := NewElements()
		elements.Add("test", 1.22)
		logNode := NewLogNode(now, elements)
		Convey("Creates new log node with the proper fields", func() {
			So(logNode.Time.Equal(now), ShouldBeTrue)
			So(logNode.Elements, ShouldEqual, elements)
			So((*logNode.Elements)[0].Name, ShouldEqual, "test")
			So((*logNode.Elements)[0].Val, ShouldEqual, 1.22)
		})
	})
	Convey("Given Node", t, func() {
		Convey("Creates new node on valid date", nil)
		Convey("Generates error on invalid date", nil)
	})
}
