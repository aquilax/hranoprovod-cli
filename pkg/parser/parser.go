// Package parser provides function to parse hranoprovod formatted files
package parser

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/aquilax/hranoprovod-cli/v3/pkg/node"
)

const (
	DefaultDateFormat  = "2006/01/02"
	DefaultCommentChar = '#'
	runeTab            = '\t'
	runeSpace          = ' '
	runeArrayItem      = '-'
)

const (
	trimText = "\t \n:\"-"
	trimQty  = "\t \n:\""
)

// Config contains the parser configuration
type Config struct {
	// CommentChar contains the character used to indicate that the line is a comment
	CommentChar uint8
}

// NewDefaultConfig returns the default set of parser configuration
func NewDefaultConfig() Config {
	return Config{DefaultCommentChar}
}

// Parser is the parser data structure
type Parser struct {
	config Config
	Nodes  chan *node.ParserNode
	Errors chan error
	Done   chan bool
}

// NewParser returns new parser
func NewParser(c Config) Parser {
	return Parser{
		config: c,
		Nodes:  make(chan *node.ParserNode),
		Errors: make(chan error),
		Done:   make(chan bool),
	}
}

// ParseFile parsers the contents of file
func (p Parser) ParseFile(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		p.Errors <- NewErrorIO(err, fileName)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	p.ParseStream(f)
}

// ParseCallback is called on node or error event when parsing the stream
type ParseCallback func(n *node.ParserNode, err error) (stop bool, cbError error)

func ParseFileCallback(fileName string, c Config, callback ParseCallback) error {
	f, err := os.Open(fileName)
	if err != nil {
		return NewErrorIO(err, fileName)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	return ParseStreamCallback(f, c, callback)
}

// ParseStreamCallback parses stream and calls callback on node or error
func ParseStreamCallback(reader io.Reader, c Config, callback ParseCallback) error {
	var pn *node.ParserNode
	var line string
	var trimmedLine string
	var title string
	var sQty string
	var separatorPos int
	var err error
	var fQty float64
	var mp *node.MetadataPair

	lineNumber := 0
	lineScanner := bufio.NewScanner(reader)
	for lineScanner.Scan() {
		lineNumber++
		line = lineScanner.Text()
		trimmedLine = strings.Trim(line, trimText)

		//skip empty lines and lines starting with #
		if trimmedLine == "" || line[0] == c.CommentChar {
			continue
		}

		//new nodes start at the beginning of the line
		if line[0] != runeSpace && line[0] != runeTab && line[0] != runeArrayItem {
			if pn != nil {
				// flush complete node
				if stop, err := callback(pn, nil); stop {
					return err
				}
			}
			// start new node
			pn = node.NewParserNode(trimmedLine)
			continue
		}

		if pn != nil {
			if trimmedLine[0] == c.CommentChar {
				// Metadata
				mp, _ = getMetadataPair(trimmedLine)
				if mp != nil {
					if pn.Metadata == nil {
						pn.Metadata = &node.Metadata{*mp}
					} else {
						*pn.Metadata = append(*pn.Metadata, *mp)
					}
				}
				continue
			}
			separatorPos = strings.LastIndexAny(trimmedLine, "\t ")

			if separatorPos == -1 {
				if stop, err := callback(nil, NewErrorBadSyntax(lineNumber, line)); stop {
					return err
				}
				continue
			}
			title = strings.Trim(trimmedLine[0:separatorPos], trimText)

			//get element value
			sQty = strings.Trim(trimmedLine[separatorPos:], trimQty)
			fQty, err = strconv.ParseFloat(sQty, 64)
			if err != nil {
				if stop, err := callback(nil, NewErrorConversion(err, sQty, lineNumber, line)); stop {
					return err
				}
				continue
			}

			pn.Elements.Add(title, fQty)
		}
	}
	// push last node
	if pn != nil {
		_, err = callback(pn, nil)
		return err
	}
	return nil
}

// ParseStream parses the contents of stream
func (p Parser) ParseStream(reader io.Reader) {
	if err := ParseStreamCallback(reader, p.config, func(n *node.ParserNode, err error) (stop bool, cbError error) {
		if err != nil {
			p.Errors <- err
			return true, err
		}
		p.Nodes <- n
		return false, nil
	}); err != nil {
		p.Errors <- err
	}
	p.Done <- true
}

func getMetadataPair(line string) (*node.MetadataPair, error) {
	trimmedLine := strings.TrimSpace(strings.Trim(line, "#"))
	separatorPos := strings.Index(trimmedLine, ":")
	if separatorPos > -1 {
		return &node.MetadataPair{
			Name:  strings.Trim(trimmedLine[:separatorPos], "# \t"),
			Value: strings.TrimSpace(trimmedLine[separatorPos+1:]),
		}, nil
	}
	return &node.MetadataPair{
		Name:  "",
		Value: trimmedLine,
	}, nil
}
