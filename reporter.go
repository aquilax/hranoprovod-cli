package main

import (
	"io"
	"fmt"
)

type Reporter struct {
	output  io.Writer
}

func NewReporter(writer io.Writer) *Reporter{
	return &Reporter{
		writer,
	}
}

func (r *Reporter) PrintAPISearchResult(nl APINodeList) error {
	for _, n := range nl {
		err := r.PrintAPINode(n)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Reporter) PrintAPINode(n APINode) error {
	_, err := fmt.Fprintln(r.output, n.Name)
	return err
}