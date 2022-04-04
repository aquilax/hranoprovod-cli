package utils

import (
	"io"
	"os"

	"github.com/aquilax/hranoprovod-cli/v2/cmd/hranoprovod-cli/internal/options"
	"github.com/urfave/cli/v2"
)

type CmdUtils struct {
	WithFileReaders func(fileNames []string, cb func([]io.Reader) error) error
	WithOptions     func(c *cli.Context, cb func(*options.Options) error) error
}

func NewCmdUtils() CmdUtils {
	return CmdUtils{
		WithFileReaders: func(fileNames []string, cb func([]io.Reader) error) error {
			result := make([]io.Reader, len(fileNames))
			for i, fileName := range fileNames {
				f, err := os.Open(fileName)
				if err != nil {
					return err
				}
				defer f.Close()
				result[i] = f
			}
			return cb(result)
		},
		WithOptions: func(c *cli.Context, cb func(*options.Options) error) error {
			o := options.New()
			if err := o.Load(c, true); err != nil {
				return err
			}
			return cb(o)
		},
	}
}
