package reporter

import (
	"bufio"
	"bytes"
	"testing"
	"time"

	"github.com/aquilax/hranoprovod-cli/v2"
	"github.com/stretchr/testify/assert"
)

func TestSummaryReporterTemplate_Process(t *testing.T) {
	tests := []struct {
		name string
		db   hranoprovod.DBNodeMap
		ln   *hranoprovod.LogNode
		want string
	}{
		{
			"generates summary report",
			hranoprovod.DBNodeMap{
				"test1": &hranoprovod.DBNode{
					Header: "test1",
					Elements: hranoprovod.Elements{
						hranoprovod.Element{Name: "energy", Value: 10},
						hranoprovod.Element{Name: "protein", Value: 20},
					},
				},
				"test2": &hranoprovod.DBNode{
					Header: "test2",
					Elements: hranoprovod.Elements{
						hranoprovod.Element{Name: "energy", Value: 20},
						hranoprovod.Element{Name: "protein", Value: 30},
					},
				},
			},
			hranoprovod.NewLogNode(time.Date(2019, 10, 10, 0, 0, 0, 0, time.UTC), hranoprovod.Elements{
				hranoprovod.NewElement("test1", 10),
				hranoprovod.NewElement("test2", 20),
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
			c := NewDefaultConfig()
			c.Color = false
			c.Output = w
			r := NewSummaryReporterTemplate(c, tt.db)
			if err := r.Process(tt.ln); err != nil {
				t.Errorf("SummaryReporterTemplate.Process() error = %v", err)
			}
			w.Flush()
			assert.Equal(t, tt.want, b.String())
		})
	}
}
