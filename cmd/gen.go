package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func newGenCommand(ol optionLoader, a *cli.App) *cli.Command {
	return &cli.Command{
		Name:  "gen",
		Usage: "Generate documentation",
		Subcommands: []*cli.Command{
			{
				Name:  "man",
				Usage: "Generate man page",
				Action: func(c *cli.Context) error {
					o, err := ol(c)
					if err != nil {
						return err
					}
					man, err := a.ToMan()
					if err != nil {
						return err
					}
					_, err = fmt.Fprint(o.ReporterConfig.Output, man)
					return err
				},
			},
			{
				Name:  "markdown",
				Usage: "Generate markdown page",
				Action: func(c *cli.Context) error {
					o, err := ol(c)
					if err != nil {
						return err
					}
					markdown, err := a.ToMarkdown()
					if err != nil {
						return err
					}
					_, err = fmt.Fprint(o.ReporterConfig.Output, markdown)
					return err
				},
			},
		},
	}
}
