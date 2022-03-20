package reporter

import (
	"bufio"
	"bytes"
	"io"
	"testing"

	"github.com/aquilax/hranoprovod-cli/v2"
)

func getSimpleTree() *hranoprovod.TreeNode {
	root := hranoprovod.NewTreeNode("test", 10.0)
	root.Add(hranoprovod.NewTreeNode("child1", 10.0))
	child2 := root.Add(hranoprovod.NewTreeNode("child2", 10.0))
	child2.Add(hranoprovod.NewTreeNode("child2.1", 10.0)).Add(hranoprovod.NewTreeNode("child2.1.1", 10.0))
	return root
}

func Test_balanceReporter_printNode(t *testing.T) {
	buffer := bytes.NewBufferString("")
	config := NewDefaultConfig()
	config.CollapseLast = true

	type fields struct {
		config Config
		db     hranoprovod.DBNodeMap
		output io.Writer
		root   *hranoprovod.TreeNode
	}
	type args struct {
		node  *hranoprovod.TreeNode
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
			r := &balanceReporter{
				config: tt.fields.config,
				db:     tt.fields.db,
				output: bufio.NewWriter(tt.fields.output),
				root:   tt.fields.root,
			}
			print(tt.args.node.Name)
			r.printNodeCollapsed(tt.args.node, tt.args.level)
			r.output.Flush()
			got := buffer.String()
			if got != tt.want {
				t.Errorf("Output = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}
