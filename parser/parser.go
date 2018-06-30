// Package parser provides function to parse hrandoprovod formatted files
package parser

import (
	"bufio"
	"github.com/aquilax/hranoprovod-cli/shared"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	runeTab   = '\t'
	runeSpace = ' '
)

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
	Nodes   chan *shared.Node
	Errors  chan error
	Done    chan bool
}

// NewParser returns new parser
func NewParser(options *Options) *Parser {
	return &Parser{
		options,
		make(chan *shared.Node),
		make(chan error),
		make(chan bool),
	}
}

// ParseFile parsers the contents of file
func (p *Parser) ParseFile(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		p.Errors <- NewErrorIO(err, fileName)
		return
	}
	defer f.Close()
	p.ParseStream(f)
}

// ParseStream parses the contents of stream
func (p *Parser) ParseStream(reader io.Reader) {
	var node *shared.Node
	lineNumber := 0
	lineScanner := bufio.NewScanner(reader)
	for lineScanner.Scan() {
		lineNumber++
		line := lineScanner.Text()
		trimmedLine := mytrim(line)

		//skip empty lines and lines starting with #
		if trimmedLine == "" || line[0] == p.options.CommentChar {
			continue
		}

		//new nodes start at the beginning of the line
		if line[0] != runeSpace && line[0] != runeTab {
			if node != nil {
				p.Nodes <- node
			}
			node = shared.NewNode(trimmedLine)
			continue
		}

		if node != nil {
			separator := strings.LastIndexAny(trimmedLine, "\t ")

			if separator == -1 {
				p.Errors <- NewErrorBadSyntax(lineNumber, line)
				return
			}
			ename := mytrim(trimmedLine[0:separator])

			//get element value
			snum := mytrim(trimmedLine[separator:])
			enum, err := strconv.ParseFloat(snum, 32)
			if err != nil {
				p.Errors <- NewErrorConversion(snum, lineNumber, line)
				return
			}

			if ndx, exists := node.Elements.Index(ename); exists {
				(*node.Elements)[ndx].Val += float32(enum)
			} else {
				node.Elements.Add(ename, float32(enum))
			}
		}
	}
	// push last node
	if node != nil {
		p.Nodes <- node
	}
	p.Done <- true
}
