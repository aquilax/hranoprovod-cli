package lint

import (
	"fmt"
	"io"

	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/options"
	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/utils"
	shared "github.com/aquilax/hranoprovod-cli/v3"
	"github.com/aquilax/hranoprovod-cli/v3/parser"
	"github.com/urfave/cli/v2"
)

type lintCmd func(stream io.Reader, lc LintConfig) error

func Command() *cli.Command {
	return newLintCommand(utils.NewCmdUtils(), Lint)
}

func newLintCommand(cu utils.CmdUtils, lint lintCmd) *cli.Command {
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
			return cu.WithOptions(c, func(o *options.Options) error {
				return cu.WithFileReaders([]string{c.Args().First()}, func(streams []io.Reader) error {
					streamToLint := streams[0]
					return lint(streamToLint, LintConfig{
						Silent:         c.IsSet("silent"),
						ParserConfig:   o.ParserConfig,
						ReporterConfig: o.ReporterConfig,
					})
				})
			})
		},
	}
}

type LintConfig struct {
	Silent         bool
	ParserConfig   parser.Config
	ReporterConfig reporter.Config
}

// Lint lints file
func Lint(stream io.Reader, lc LintConfig) error {
	err := parser.ParseStreamCallback(stream, lc.ParserConfig, func(node *shared.ParserNode, err error) (stop bool, cbError error) {
		if err != nil {
			fmt.Fprintln(lc.ReporterConfig.Output, err)
		}
		return false, nil
	})
	if err != nil {
		return err
	}
	if !lc.Silent {
		fmt.Fprintln(lc.ReporterConfig.Output, "No errors found")
	}
	return nil
}
