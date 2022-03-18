package cmd

import (
	"io"

	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/urfave/cli/v2"
)

func newPrintCommand(cu cmdUtils, printCb app.PrintCmd) *cli.Command {
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
			return cu.withOptions(c, func(o *app.Options) error {
				return cu.withFileReaders([]string{c.Args().First()}, func(streams []io.Reader) error {
					logStream := streams[0]
					return printCb(logStream, o.GlobalConfig.DateFormat, o.ParserConfig, o.ReporterConfig, o.FilterConfig)
				})
			})
		},
	}
}
