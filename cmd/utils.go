package cmd

import (
	"io"
	"os"
)

func getFileReader(fileName string) (io.Reader, error) {
	f, err := os.Open(fileName)
	defer f.Close()
	return f, err
}
