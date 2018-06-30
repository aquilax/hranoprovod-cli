package shared

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestElement(t *testing.T) {
	Convey("NewElement", t, func() {
		el := NewElement("test", 10)
		Convey("Creates new element", func() {
			So(el.Name, ShouldEqual, "test")
			So(el.Val, ShouldEqual, 10)
		})
	})
}

func TestElements(t *testing.T) {
	Convey("Given Elements", t, func() {
		el := NewElements()
		Convey("Calling Add", func() {
			el.Add("test", 10)
			Convey("Adds the element to the list", func() {
				So(el.Len(), ShouldEqual, 1)
			})
			Convey("After adding more elements", func() {
				el.Add("test3", 13)
				el.Add("test2", 12)
				el.Add("test1", 11)
				Convey("Calling Index on present element", func() {
					index, found := el.Index("test2")
					Convey("Returns the correct index", func() {
						So(index, ShouldEqual, 2)
					})
					Convey("Returns positive found", func() {
						So(found, ShouldBeTrue)
					})
				})
				Convey("Calling Index on missing element", func() {
					_, found := el.Index("test111")
					Convey("Returns not found", func() {
						So(found, ShouldBeFalse)
					})
				})
				Convey("After Sort", func() {
					el.Sort()
					Convey("Elements are sorted", func() {
						index, _ := el.Index("test3")
						So(index, ShouldEqual, 3)
						index2, _ := el.Index("test1")
						So(index2, ShouldEqual, 1)
					})
				})
				Convey("Having second set of elements", func() {
					el2 := NewElements()
					el2.Add("test3", 113)
					el2.Add("test2", 112)
					el2.Add("test1", 111)
					el2.Add("test4", 444)
					Convey("SumMerge with coef 2", func() {
						el.SumMerge(el2, 2)
						Convey("Returns correct elements", func() {
							index, found := el.Index("test1")
							So(found, ShouldBeTrue)
							So(index, ShouldEqual, 3)
							So((*el)[index].Val, ShouldEqual, 233)
						})
						Convey("New elements are added", func() {
							index, found := el.Index("test4")
							So(found, ShouldBeTrue)
							So((*el)[index].Val, ShouldEqual, 888)
						})
					})
				})
			})
		})
	})
}
