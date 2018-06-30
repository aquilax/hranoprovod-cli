package accumulator

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAccumulator(t *testing.T) {
	Convey("Given the acccumulator", t, func() {
		acc := NewAccumulator()
		Convey("When a positive element is added", func() {
			acc.Add("test", 1.22)
			Convey("It should go to the positive acccumulator", func() {
				So((*acc)["test"][Positive], ShouldEqual, 1.22)
			})
			Convey("When a positive value is added to the same key", func() {
				acc.Add("test", 2.33)
				Convey("It is accumulated in the positive register", func() {
					So((*acc)["test"][Positive], ShouldEqual, 3.55)
				})
			})
			Convey("When a negative value is added to the same key", func() {
				acc.Add("test", -2.33)
				Convey("It is accumulated in the positive register", func() {
					So((*acc)["test"][Negative], ShouldEqual, -2.33)
					So((*acc)["test"][Positive], ShouldEqual, 1.22)
				})
			})
			Convey("When negative element is added", func() {
				acc.Add("test2", -1.32)
				Convey("It should go to the negative acccumulator", func() {
					So((*acc)["test2"][Negative], ShouldEqual, -1.32)
				})
			})
		})
	})
}
