package testutils

import (
	"embed"
	"io"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/options"
	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/utils"
)

//go:embed testAssets/*
var content embed.FS

func GetMockApp(cmd *cli.Command) cli.App {
	return cli.App{
		Commands: cli.Commands{cmd},
	}
}

func GetMockCmdUtils(contents []string, o *options.Options) utils.CmdUtils {
	return utils.CmdUtils{
		WithFileReaders: func(fileNames []string, cb func([]io.Reader) error) error {
			streams := make([]io.Reader, len(fileNames))
			for i := range fileNames {
				streams[i] = strings.NewReader(contents[i])
			}
			return cb(streams)
		},
		WithOptions: func(c *cli.Context, cb func(*options.Options) error) error {
			return cb(o)
		},
	}
}

func GetMockCmdUtilsRealOptions(contents []string, output io.Writer) utils.CmdUtils {
	return utils.CmdUtils{
		WithFileReaders: func(fileNames []string, cb func([]io.Reader) error) error {
			streams := make([]io.Reader, len(fileNames))
			for i := range fileNames {
				streams[i] = strings.NewReader(contents[i])
			}
			return cb(streams)
		},
		WithOptions: func(c *cli.Context, cb func(*options.Options) error) error {
			o := options.New()
			if err := o.Load(c, false); err != nil {
				return err
			}
			o.ReporterConfig.Color = false
			o.ReporterConfig.Output = output
			return cb(o)
		},
	}
}

func ReadAsset(fileName string) ([]byte, error) {
	return content.ReadFile(fileName)
}
