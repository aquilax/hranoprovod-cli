package lint

import (
	"fmt"
	"io"

	"github.com/aquilax/hranoprovod-cli/v2/cmd/hranoprovod-cli/internal/options"
	"github.com/aquilax/hranoprovod-cli/v2/cmd/hranoprovod-cli/internal/utils"
	"github.com/aquilax/hranoprovod-cli/v2/lib/parser"
	"github.com/aquilax/hranoprovod-cli/v2/lib/reporter"
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
	parser := parser.NewParser(lc.ParserConfig)
	go parser.ParseStream(stream)
	err := func() error {
		for {
			select {
			case <-parser.Nodes:
			case err := <-parser.Errors:
				fmt.Fprintln(lc.ReporterConfig.Output, err)
			case <-parser.Done:
				return nil
			}
		}
	}()
	if err != nil {
		return err
	}
	if !lc.Silent {
		fmt.Fprintln(lc.ReporterConfig.Output, "No errors found")
	}
	return nil
}
