package reporter

import (
	"bytes"
	"io"
	"testing"

	"github.com/aquilax/hranoprovod-cli/v2/accumulator"
	"github.com/aquilax/hranoprovod-cli/v2/shared"
)

func getSimpleTree() *accumulator.TreeNode {
	root := accumulator.NewTreeNode("test", 10.0)
	root.Add(accumulator.NewTreeNode("child1", 10.0))
	child2 := root.Add(accumulator.NewTreeNode("child2", 10.0))
	child2.Add(accumulator.NewTreeNode("child2.1", 10.0)).Add(accumulator.NewTreeNode("child2.1.1", 10.0))
	return root
}

func Test_balanceReporter_printNode(t *testing.T) {
	buffer := bytes.NewBufferString("")
	options := NewDefaultOptions()
	options.CollapseLast = true

	type fields struct {
		options *Options
		db      shared.DBNodeList
		output  io.Writer
		root    *accumulator.TreeNode
	}
	type args struct {
		node  *accumulator.TreeNode
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
				options: options,
				db:      nil,
				output:  buffer,
				root:    nil,
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
			r := &balanceReporter{
				options: tt.fields.options,
				db:      tt.fields.db,
				output:  tt.fields.output,
				root:    tt.fields.root,
			}
			print(tt.args.node.Name)
			r.printNodeCollapsed(tt.args.node, tt.args.level)
			got := buffer.String()
			if got != tt.want {
				t.Errorf("Output = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}
