package shared

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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
