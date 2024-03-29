package balance

import (
	"bufio"
	"bytes"
	"io"
	"testing"

	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/reporter"
	shared "github.com/aquilax/hranoprovod-cli/v3"
)

func getSimpleTree() *shared.TreeNode {
	root := shared.NewTreeNode("test", 10.0)
	root.Add(shared.NewTreeNode("child1", 10.0))
	child2 := root.Add(shared.NewTreeNode("child2", 10.0))
	child2.Add(shared.NewTreeNode("child2.1", 10.0)).Add(shared.NewTreeNode("child2.1.1", 10.0))
	return root
}

func Test_balanceReporter_printNode(t *testing.T) {
	buffer := bytes.NewBufferString("")
	config := reporter.NewDefaultConfig()
	config.CollapseLast = true

	type fields struct {
		config reporter.Config
		db     shared.DBNodeMap
		output io.Writer
		root   *shared.TreeNode
	}
	type args struct {
		node  *shared.TreeNode
		level int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Test simple tree",
			fields: fields{
				config: config,
				db:     nil,
				output: buffer,
				root:   nil,
			},
			args: args{
				node:  getSimpleTree(),
				level: 0,
			},
			want: `     10.00 | child1
     10.00 | child2/child2.1/child2.1.1
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := balanceReporter{
				db:           tt.fields.db,
				output:       bufio.NewWriter(tt.fields.output),
				root:         tt.fields.root,
				collapseLast: tt.fields.config.CollapseLast,
			}
			if err := printNodeCollapsed(tt.args.node, tt.args.level, r.output); err != nil {
				t.Error(err)
			}
			r.output.Flush()
			got := buffer.String()
			if got != tt.want {
				t.Errorf("Output = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}
