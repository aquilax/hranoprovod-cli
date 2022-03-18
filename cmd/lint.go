package cmd

import (
	"fmt"
	"io"

	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/aquilax/hranoprovod-cli/v2/options"
	"github.com/urfave/cli/v2"
)

type LintCmd func(stream io.Reader, lc app.LintConfig) error

func newLintCommand(cu cmdUtils, lint LintCmd) *cli.Command {
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
			return cu.withOptions(c, func(o *options.Options) error {
				return cu.withFileReaders([]string{c.Args().First()}, func(streams []io.Reader) error {
					streamToLint := streams[0]
					return lint(streamToLint, app.LintConfig{
						Silent:         c.IsSet("silent"),
						ParserConfig:   o.ParserConfig,
						ReporterConfig: o.ReporterConfig,
					})
				})
			})
		},
	}
}
