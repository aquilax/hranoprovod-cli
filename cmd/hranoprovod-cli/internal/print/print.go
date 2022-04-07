package print

import (
	"io"

	"github.com/aquilax/hranoprovod-cli/v3/cmd/hranoprovod-cli/internal/options"
	"github.com/aquilax/hranoprovod-cli/v3/cmd/hranoprovod-cli/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/v3/cmd/hranoprovod-cli/internal/utils"
	"github.com/aquilax/hranoprovod-cli/v3/lib/filter"
	"github.com/aquilax/hranoprovod-cli/v3/lib/parser"
	"github.com/urfave/cli/v2"
)

type printCmd func(logStream io.Reader, pc PrintConfig) error

func Command() *cli.Command {
	return NewPrintCommand(utils.NewCmdUtils(), Print)
}

func NewPrintCommand(cu utils.CmdUtils, printCb printCmd) *cli.Command {
	return &cli.Command{
		Name:  "print",
		Usage: "Print log",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "begin",
				Aliases: []string{"b"},
				Usage:   "Beginning of period `DATE`",
			},
			&cli.StringFlag{
				Name:    "end",
				Aliases: []string{"e"},
				Usage:   "End of period `DATE`",
			},
		},
		Action: func(c *cli.Context) error {
			return cu.WithOptions(c, func(o *options.Options) error {
				return cu.WithFileReaders([]string{o.GlobalConfig.LogFileName}, func(streams []io.Reader) error {
					logStream := streams[0]
					return printCb(logStream, PrintConfig{
						DateFormat:     o.GlobalConfig.DateFormat,
						ParserConfig:   o.ParserConfig,
						ReporterConfig: o.ReporterConfig,
						FilterConfig:   o.FilterConfig,
					})
				})
			})
		},
	}
}

type PrintConfig struct {
	DateFormat     string
	ParserConfig   parser.Config
	ReporterConfig reporter.Config
	FilterConfig   filter.Config
}

// Print reads and prints back out the log file
func Print(logStream io.Reader, pc PrintConfig) error {
	r := NewPrintReporter(pc.ReporterConfig)
	f := filter.GetIntervalNodeFilter(pc.FilterConfig)
	return utils.WalkNodesInStream(logStream, pc.DateFormat, pc.ParserConfig, f, r)
}
