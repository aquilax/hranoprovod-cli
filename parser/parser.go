// Package parser provides function to parse hranoprovod formatted files
package parser

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/aquilax/hranoprovod-cli/v2/shared"
)

const (
	runeTab       = '\t'
	runeSpace     = ' '
	runeArrayItem = '-'
)

func trim(s string) string {
	return strings.Trim(s, "\t \n:\"-")
}

// Options contains the parser related options
type Options struct {
	// CommentChar contains the character used to indicate that the line is a comment
	CommentChar uint8
}

// NewDefaultOptions returns the default set of parser options
func NewDefaultOptions() *Options {
	return &Options{'#'}
}

// Parser is the parser data structure
type Parser struct {
	options *Options
	Nodes   chan *shared.ParserNode
	Errors  chan error
	Done    chan bool
}

// NewParser returns new parser
func NewParser(options *Options) Parser {
	return Parser{
		options,
		make(chan *shared.ParserNode),
		make(chan error),
		make(chan bool),
	}
}

// ParseFile parsers the contents of file
func (p Parser) ParseFile(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		p.Errors <- NewErrorIO(err, fileName)
		return
	}
	defer f.Close()
	p.ParseStream(f)
}

// ParseCallback is called on node or error event when parsing the stream
type ParseCallback func(n *shared.ParserNode, err error) (stop bool)

// ParseStreamCallback parses stream and calls callback on node or error
func ParseStreamCallback(reader io.Reader, commentChar uint8, callback ParseCallback) {
	var node *shared.ParserNode
	var line string
	var trimmedLine string
	var title string
	var sValue string
	var separatorPos int
	var err error
	var fValue float64

	lineNumber := 0
	lineScanner := bufio.NewScanner(reader)
	for lineScanner.Scan() {
		lineNumber++
		line = lineScanner.Text()
		trimmedLine = trim(line)

		//skip empty lines and lines starting with #
		if trimmedLine == "" || line[0] == commentChar {
			continue
		}

		//new nodes start at the beginning of the line
		if line[0] != runeSpace && line[0] != runeTab && line[0] != runeArrayItem {
			if node != nil {
				// flush complete node
				callback(node, nil)
			}
			// start new node
			node = shared.NewParserNode(trimmedLine)
			continue
		}

		if node != nil {
			separatorPos = strings.LastIndexAny(trimmedLine, "\t ")

			if separatorPos == -1 {
				callback(nil, NewErrorBadSyntax(lineNumber, line))
				return
			}
			title = trim(trimmedLine[0:separatorPos])

			//get element value
			sValue = trim(trimmedLine[separatorPos:])
			fValue, err = strconv.ParseFloat(sValue, 64)
			if err != nil {
				callback(nil, NewErrorConversion(sValue, lineNumber, line))
				return
			}

			node.Elements.Add(title, fValue)
		}
	}
	// push last node
	if node != nil {
		callback(node, nil)
	}
}

// ParseStream parses the contents of stream
func (p Parser) ParseStream(reader io.Reader) {
	ParseStreamCallback(reader, p.options.CommentChar, func(n *shared.ParserNode, err error) (stop bool) {
		if err != nil {
			p.Errors <- err
			return true
		}
		p.Nodes <- n
		return false
	})
	p.Done <- true
}
