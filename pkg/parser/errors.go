package parser

import (
	"fmt"
)

// ErrorIO raised on IO issue
type ErrorIO struct {
	err      error
	FileName string
}

// NewErrorIO creates new IO Error
func NewErrorIO(err error, fileName string) *ErrorIO {
	return &ErrorIO{err, fileName}
}

// Error returns the error message
func (e *ErrorIO) Error() string {
	return e.err.Error()
}

// ErrorBadSyntax used when the stream contains bad syntax
type ErrorBadSyntax struct {
	LineNumber int
	Line       string
}

// NewErrorBadSyntax creates new ErrorBadSyntax error
func NewErrorBadSyntax(lineNumber int, line string) *ErrorBadSyntax {
	return &ErrorBadSyntax{lineNumber, line}
}

// Error returns the error message
func (e *ErrorBadSyntax) Error() string {
	return fmt.Sprintf("bad syntax on line %d, \"%s\".", e.LineNumber, e.Line)
}

// ErrorConversion raised when the element value cannot be parsed as float
type ErrorConversion struct {
	err        error
	Text       string
	LineNumber int
	Line       string
}

// NewErrorConversion creates new ErrorConversion error
func NewErrorConversion(err error, text string, lineNumber int, line string) *ErrorConversion {
	return &ErrorConversion{err, text, lineNumber, line}
}

// Error returns the error message
func (e *ErrorConversion) Error() string {
	return fmt.Sprintf("error converting \"%s\" to float on line %d \"%s\".", e.Text, e.LineNumber, e.Line)
}
