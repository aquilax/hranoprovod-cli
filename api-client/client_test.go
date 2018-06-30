package client

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type URLTestCase struct {
	base   string
	path   string
	params URLParams
	result string
	err    bool
}

func TestAPIClient(t *testing.T) {
	Convey("Given api client", t, func() {
		Convey("URL", func() {
			var testCases = []URLTestCase{
				{
					base:   "http://example.com",
					path:   "test",
					params: nil,
					result: "http://example.com/test",
					err:    false,
				},
				{
					base:   "http://example.com/",
					path:   "test",
					params: nil,
					result: "http://example.com/test",
					err:    false,
				},
				{
					base:   "http://example.com/",
					path:   "/test/test2/",
					params: nil,
					result: "http://example.com/test/test2",
					err:    false,
				},
				{
					base:   ":::",
					path:   "test",
					params: nil,
					result: "",
					err:    true,
				},
				{
					base:   "http://example.com/",
					path:   "/test/test2/",
					params: URLParams{"test": "one"},
					result: "http://example.com/test/test2?test=one",
					err:    false,
				},
				{
					base:   "http://example.com/",
					path:   "/test/test2/",
					params: URLParams{"test": "one", "test2": "two"},
					result: "http://example.com/test/test2?test=one&test2=two",
					err:    false,
				},
			}
			for _, tc := range testCases {
				c := NewAPIClient(&Options{tc.base})
				result, err := c.buildURL(tc.path, tc.params)
				So(result, ShouldEqual, tc.result)
				if tc.err {
					So(err, ShouldNotBeNil)
				} else {
					So(err, ShouldBeNil)
				}
			}
		})
	})
}
