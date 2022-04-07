package summary

import (
	"bufio"
	"bytes"
	"testing"
	"time"

	"github.com/aquilax/hranoprovod-cli/v3/cmd/hranoprovod-cli/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/v3/lib/shared"
	"github.com/stretchr/testify/assert"
)

func TestSummaryReporterTemplate_Process(t *testing.T) {
	tests := []struct {
		name string
		db   shared.DBNodeMap
		ln   *shared.LogNode
		want string
	}{
		{
			"generates summary report",
			shared.DBNodeMap{
				"test1": &shared.DBNode{
					Header: "test1",
					Elements: shared.Elements{
						shared.Element{Name: "energy", Value: 10},
						shared.Element{Name: "protein", Value: 20},
					},
				},
				"test2": &shared.DBNode{
					Header: "test2",
					Elements: shared.Elements{
						shared.Element{Name: "energy", Value: 20},
						shared.Element{Name: "protein", Value: 30},
					},
				},
			},
			shared.NewLogNode(time.Date(2019, 10, 10, 0, 0, 0, 0, time.UTC), shared.Elements{
				shared.NewElement("test1", 10),
				shared.NewElement("test2", 20),
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
			w.Flush()
			assert.Equal(t, tt.want, b.String())
		})
	}
}
