package gen

import (
	"fmt"

	"github.com/aquilax/hranoprovod-cli/v2/cmd/hranoprovod-cli/internal/options"
	"github.com/aquilax/hranoprovod-cli/v2/cmd/hranoprovod-cli/internal/utils"
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return newGenCommand(utils.NewCmdUtils())
}

func newGenCommand(cu utils.CmdUtils) *cli.Command {
	return &cli.Command{
		Name:  "gen",
		Usage: "Generate documentation",
		Subcommands: []*cli.Command{
			{
				Name:  "man",
				Usage: "Generate man page",
				Action: func(c *cli.Context) error {
					return cu.WithOptions(c, func(o *options.Options) error {
						man, err := c.App.ToMan()
						if err != nil {
							return err
						}
						_, err = fmt.Fprint(o.ReporterConfig.Output, man)
						return err
					})
				},
			},
			{
				Name:  "markdown",
				Usage: "Generate markdown page",
				Action: func(c *cli.Context) error {
					return cu.WithOptions(c, func(o *options.Options) error {
						markdown, err := c.App.ToMarkdown()
						if err != nil {
							return err
						}
						_, err = fmt.Fprint(o.ReporterConfig.Output, markdown)
						return err
					})
				},
			},
		},
	}
}
