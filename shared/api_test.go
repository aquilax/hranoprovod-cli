package shared

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAccumulator(t *testing.T) {
	Convey("Given APIError", t, func() {
		err := APIError{
			IsError: true,
			Code:    100,
			Status:  "status",
			Message: "message",
		}
		So(err.Error(), ShouldEqual, "status: message")
	})
}
