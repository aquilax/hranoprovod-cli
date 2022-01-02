package cmd

import (
	"io"
	"strings"

	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/urfave/cli/v2"
)

func getMockApp(cmd *cli.Command) cli.App {
	return cli.App{
		Commands: cli.Commands{cmd},
	}
}

func getMockCu(content string) cmdUtils {
	return cmdUtils{
		func(fileName string, cb func(io.Reader) error) error {
			return cb(strings.NewReader(content))
		},
		func(c *cli.Context, cb func(*app.Options) error) error {
			return cb(app.NewOptions())
		},
	}
}
