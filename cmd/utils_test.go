package cmd

import (
	"io"
	"strings"

	"github.com/aquilax/hranoprovod-cli/v2/options"
	"github.com/urfave/cli/v2"
)

func getMockApp(cmd *cli.Command) cli.App {
	return cli.App{
		Commands: cli.Commands{cmd},
	}
}

func getMockCu(contents []string) cmdUtils {
	return cmdUtils{
		func(fileNames []string, cb func([]io.Reader) error) error {
			streams := make([]io.Reader, len(fileNames))
			for i := range fileNames {
				streams[i] = strings.NewReader(contents[i])
			}
			return cb(streams)
		},
		func(c *cli.Context, cb func(*options.Options) error) error {
			return cb(options.New())
		},
	}
}
