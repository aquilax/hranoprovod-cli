package parser_test

import (
	"fmt"
	"strings"

	"github.com/aquilax/hranoprovod-cli/v2/parser"
)

func ExampleParser_ParseStream() {
	p := parser.NewParser(parser.NewDefaultOptions())
	file := `2011/07/17:
  el1: 1.22
  ел 2:  4
  el/3:  3

2011/07/18:
  el1: 1.33
  ел 5:  5
  el/7:  4
  el1: 1.35
`
	go p.ParseStream(strings.NewReader(file))
	func() {
		for {
			select {
			case node := <-p.Nodes:
				fmt.Println(node.Header)
			case <-p.Errors:
				return
			case <-p.Done:
				return
			}
		}
	}()
	// Output:
	// 2011/07/17
	// 2011/07/18
}
