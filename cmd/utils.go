package cmd

import (
	"io"
	"os"

	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/urfave/cli/v2"
)

type cmdUtils struct {
	withFileReader func(fileName string, cb func(io.Reader) error) error
	withOptions    func(c *cli.Context, cb func(*app.Options) error) error
}

func NewCmdUtils() cmdUtils {
	return cmdUtils{
		withFileReader: func(fileName string, cb func(io.Reader) error) error {
			return withFileReader(fileName, cb)
		},
		withOptions: func(c *cli.Context, cb func(*app.Options) error) error {
			o := app.NewOptions()
			if err := o.Load(c); err != nil {
				return err
			}
			return cb(o)
		},
	}
}

func withFileReader(fileName string, cb func(io.Reader) error) error {
	if f, err := os.Open(fileName); err != nil {
		return err
	} else {
		defer f.Close()
		return cb(f)
	}
}
