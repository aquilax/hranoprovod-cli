package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func newGenCommand(a *cli.App) *cli.Command {
	return &cli.Command{
		Name:  "gen",
		Usage: "Generate documentation",
		Subcommands: []*cli.Command{
			{
				Name:  "man",
				Usage: "Generate man page",
				Action: func(c *cli.Context) error {
					man, err := a.ToMan()
					if err != nil {
						return err
					}
					_, err = fmt.Fprint(os.Stdout, man)
					return err
				},
			},
			{
				Name:  "markdown",
				Usage: "Generate markdown page",
				Action: func(c *cli.Context) error {
					markdown, err := a.ToMarkdown()
					if err != nil {
						return err
					}
					_, err = fmt.Fprint(os.Stdout, markdown)
					return err
				},
			},
		},
	}
}
