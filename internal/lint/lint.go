package lint

import (
	"fmt"
	"io"

	"github.com/aquilax/hranoprovod-cli/v3/internal/options"
	"github.com/aquilax/hranoprovod-cli/v3/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/v3/internal/utils"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/node"
	"github.com/aquilax/hranoprovod-cli/v3/pkg/parser"
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
	if err := parser.ParseStreamCallback(stream, lc.ParserConfig, func(node *node.ParserNode, err error) (stop bool, cbError error) {
		if err != nil {
			if _, err := fmt.Fprintln(lc.ReporterConfig.Output, err); err != nil {
				panic(err)
			}
		}
		return false, nil
	}); err != nil {
		return err
	}
	if !lc.Silent {
		if _, err := fmt.Fprintln(lc.ReporterConfig.Output, "No errors found"); err != nil {
			panic(err)
		}
	}
	return nil
}
