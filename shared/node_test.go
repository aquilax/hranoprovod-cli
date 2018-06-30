package shared

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNodeList(t *testing.T) {
	Convey("Given NodeList", t, func() {
		nl := NewNodeList()
		Convey("Creates new NodeList", func() {
			So(nl != nil, ShouldBeTrue)
		})
		Convey("Adding new node", func() {
			node := NewNode("test")
			nl.Push(node)
			Convey("Increases the number of nodes in the list", func() {
				So(len(*nl), ShouldEqual, 1)
			})
		})
	})
}
