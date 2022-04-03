package main

import (
	"io"
	"strings"

	"github.com/urfave/cli/v2"
)

func getMockApp(cmd *cli.Command) cli.App {
	return cli.App{
		Commands: cli.Commands{cmd},
	}
}

func getMockCmdUtils(contents []string, o *Options) cmdUtils {
	return cmdUtils{
		func(fileNames []string, cb func([]io.Reader) error) error {
			streams := make([]io.Reader, len(fileNames))
			for i := range fileNames {
				streams[i] = strings.NewReader(contents[i])
			}
			return cb(streams)
		},
		func(c *cli.Context, cb func(*Options) error) error {
			return cb(o)
		},
	}
}

func getMockCmdUtilsRealOptions(contents []string, output io.Writer) cmdUtils {
	return cmdUtils{
		func(fileNames []string, cb func([]io.Reader) error) error {
			streams := make([]io.Reader, len(fileNames))
			for i := range fileNames {
				streams[i] = strings.NewReader(contents[i])
			}
			return cb(streams)
		},
		func(c *cli.Context, cb func(*Options) error) error {
			o := New()
			if err := o.Load(c, false); err != nil {
				return err
			}
			o.ReporterConfig.Color = false
			o.ReporterConfig.Output = output
			return cb(o)
		},
	}
}
