package cmd

import (
	"io"
	"os"
)

func withFileReader(fileName string, cb func(io.Reader) error) error {
	if f, err := os.Open(fileName); err != nil {
		return err
	} else {
		defer f.Close()
		return cb(f)
	}
}
