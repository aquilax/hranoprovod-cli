package resolver

import (
	"testing"

	"github.com/aquilax/hranoprovod-cli/shared"
	. "github.com/smartystreets/goconvey/convey"
)

func TestResolver(t *testing.T) {
	Convey("Given nodes database and reslover", t, func() {
		nl := shared.NewDBNodeList()
		node1 := shared.NewParserNode("node1")
		node1.Elements.Add("element1", 100)
		node1.Elements.Add("element2", 200)
		nl.Push(shared.NewDBNodeFromNode(node1))
		node2 := shared.NewParserNode("node2")
		node2.Elements.Add("node1", 2)
		nl.Push(shared.NewDBNodeFromNode(node2))
		resolver := NewResolver(nl, 1)
		Convey("Resolve resolves the database", func() {
			resolver.Resolve()
			Convey("Elements are resolved", func() {
				n1 := nl["node1"]
				So(n1.Elements[0].Name, ShouldEqual, "element1")
				So(n1.Elements[0].Val, ShouldEqual, 100)
				So(n1.Elements[1].Name, ShouldEqual, "element2")
				So(n1.Elements[1].Val, ShouldEqual, 200)
				n2 := nl["node2"]
				So(n2.Elements[0].Name, ShouldEqual, "element1")
				So(n2.Elements[0].Val, ShouldEqual, 200)
				So(n2.Elements[1].Name, ShouldEqual, "element2")
				So(n2.Elements[1].Val, ShouldEqual, 400)
			})
		})
	})
}
