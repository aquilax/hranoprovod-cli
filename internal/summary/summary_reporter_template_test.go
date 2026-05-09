package summary

import (
	"bufio"
	"bytes"
	"testing"
	"time"

	"github.com/aquilax/hranoprovod-cli/v3/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/element"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/node"
	"github.com/stretchr/testify/assert"
)

func TestSummaryReporterTemplate_Process(t *testing.T) {
	tests := []struct {
		name string
		db   node.DBNodeMap
		ln   *node.LogNode
		want string
	}{
		{
			"generates summary report",
			node.DBNodeMap{
				"test1": &node.DBNode{
					Header: "test1",
					Elements: element.Elements{
						element.Element{Name: "energy", Value: 10},
						element.Element{Name: "protein", Value: 20},
					},
				},
				"test2": &node.DBNode{
					Header: "test2",
					Elements: element.Elements{
						element.Element{Name: "energy", Value: 20},
						element.Element{Name: "protein", Value: 30},
					},
				},
			},
			node.NewLogNode(time.Date(2019, 10, 10, 0, 0, 0, 0, time.UTC), element.Elements{
				element.NewElement("test1", 10),
				element.NewElement("test2", 20),
			}, nil),
			`2019/10/10 :
    500.00 : energy
    800.00 : protein
------------
     10.00 : test1
     20.00 : test2
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			w := bufio.NewWriter(&b)
			c := reporter.NewDefaultConfig()
			c.Color = false
			c.Output = w
			r := NewSummaryReporterTemplate(c, tt.db)
			if err := r.Process(tt.ln); err != nil {
				t.Errorf("SummaryReporterTemplate.Process() error = %v", err)
			}
			assert.NoError(t, w.Flush())
			assert.Equal(t, tt.want, b.String())
		})
	}
}
