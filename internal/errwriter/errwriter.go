package errwriter

import (
	"fmt"
	"io"
)

type ErrWriter struct {
	w   io.Writer
	Err error
}

func New(w io.Writer) *ErrWriter {
	return &ErrWriter{w: w}
}

func (ew *ErrWriter) Fprintf(format string, a ...any) {
	if ew.Err != nil {
		return
	}
	_, ew.Err = fmt.Fprintf(ew.w, format, a...)
}
