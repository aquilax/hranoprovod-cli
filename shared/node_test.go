package shared

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNodeList(t *testing.T) {
	Convey("Given NodeList", t, func() {
		nl := NewNodeList()
		Convey("Creates new NodeList", func() {
			So(nl != nil, ShouldBeTrue)
		})
		Convey("Adding new node", func() {
			node := NewParserNode("test")
			nl.Push(node)
			Convey("Increases the number of nodes in the list", func() {
				So(len(nl), ShouldEqual, 1)
			})
		})
	})
}
