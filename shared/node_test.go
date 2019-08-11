package shared

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestDBNodeList(t *testing.T) {
	Convey("Given NodeList", t, func() {
		nl := NewDBNodeList()
		Convey("Creates new DBNodeList", func() {
			So(nl != nil, ShouldBeTrue)
		})
		Convey("Adding new node", func() {
			node := NewDBNodeFromNode(NewParserNode("test"))
			nl.Push(node)
			Convey("Increases the number of nodes in the list", func() {
				So(len(nl), ShouldEqual, 1)
			})
		})
	})
}

func TestNewLogNode(t *testing.T) {
	Convey("Given NewLogNode", t, func() {
		now := time.Now()
		elements := NewElements()
		elements.Add("test", 1.22)
		logNode := NewLogNode(now, elements)
		Convey("Creates new log node with the proper fields", func() {
			So(logNode.Time.Equal(now), ShouldBeTrue)
			So((logNode.Elements)[0].Name, ShouldEqual, "test")
			So((logNode.Elements)[0].Val, ShouldEqual, 1.22)
		})
	})
	Convey("Given Node", t, func() {
		Convey("Creates new node on valid date", func() {
			node := NewParserNode("2006/01/02")
			logNode, err := NewLogNodeFromNode(node, "2006/01/02")
			So(logNode, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		Convey("Generates error on invalid date", func() {
			node := NewParserNode("2006/13/02")
			logNode, err := NewLogNodeFromNode(node, "2006/01/02")
			So(logNode, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})
}
