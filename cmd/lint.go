package cmd

import (
	"fmt"
	"io"

	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/urfave/cli/v2"
)

func newLintCommand(cu cmdUtils, lint app.LintCmd) *cli.Command {
	return &cli.Command{
		Name:      "lint",
		Usage:     "Lints file for parsing errors",
		ArgsUsage: "[FILE]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "silent",
				Aliases: []string{"s"},
				Usage:   "stay silent if no errors are found",
			},
		},
		Before: func(c *cli.Context) error {
			if c.Args().First() == "" {
				return fmt.Errorf("no file provided")
			}
			return nil
		},
		Action: func(c *cli.Context) error {
			return cu.withOptions(c, func(o *app.Options) error {
				return cu.withFileReader(c.Args().First(), func(streamToLint io.Reader) error {
					o.ParserConfig.StopOnError = false
					return lint(streamToLint, c.IsSet("silent"), o.ParserConfig, o.ReporterConfig)
				})
			})
		},
	}
}