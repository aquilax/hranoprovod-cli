package cmd

import (
	"fmt"

	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/urfave/cli/v2"
)

func newLintCommand(ol optionLoader) *cli.Command {
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
			o, err := ol(c)
			if err != nil {
				return err
			}
			o.ParserConfig.StopOnError = false
			if streamToLint, err := getFileReader(c.Args().First()); err == nil {
				return app.Lint(streamToLint, c.IsSet("silent"), o.ParserConfig, o.ReporterConfig)
			} else {
				return err
			}
		},
	}
}
