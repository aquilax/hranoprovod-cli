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
	DateFormat  string
}

// NewDefaultConfig returns the default set of parser configuration
func NewDefaultConfig() Config {
	return Config{DefaultCommentChar, DefaultDateFormat}
}

// Parser is the parser data structure
type Parser struct {
	config Config
	Nodes  chan *shared.ParserNode
	Errors chan error
	Done   chan bool
}

// NewParser returns new parser
func NewParser(c Config) Parser {
	return Parser{
		config: c,
		Nodes:  make(chan *shared.ParserNode),
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
	defer f.Close()
	p.ParseStream(f)
}

// ParseCallback is called on node or error event when parsing the stream
type ParseCallback func(n *shared.ParserNode, err error) (stop bool, cbError error)

func ParseFileCallback(fileName string, c Config, callback ParseCallback) error {
	f, err := os.Open(fileName)
	if err != nil {
		return NewErrorIO(err, fileName)
	}
	defer f.Close()
	return ParseStreamCallback(f, c, callback)
}

// ParseStreamCallback parses stream and calls callback on node or error
func ParseStreamCallback(reader io.Reader, c Config, callback ParseCallback) error {
	var node *shared.ParserNode
	var line string
	var trimmedLine string
	var title string
	var sQty string
	var separatorPos int
	var err error
	var fQty float64
	var mp *shared.MetadataPair

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
			if node != nil {
				// flush complete node
				if stop, err := callback(node, nil); stop {
					return err
				}
			}
			// start new node
			node = shared.NewParserNode(trimmedLine)
			continue
		}

		if node != nil {
			if trimmedLine[0] == c.CommentChar {
				// Metadata
				mp, _ = getMetadataPair(trimmedLine)
				if mp != nil {
					if node.Metadata == nil {
						node.Metadata = &shared.Metadata{*mp}
					} else {
						*node.Metadata = append(*node.Metadata, *mp)
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
				if stop, err := callback(nil, NewErrorConversion(sQty, lineNumber, line)); stop {
					return err
				}
				continue
			}

			node.Elements.Add(title, fQty)
		}
	}
	// push last node
	if node != nil {
		_, err = callback(node, nil)
		return err
	}
	return nil
}

// ParseStream parses the contents of stream
func (p Parser) ParseStream(reader io.Reader) {
	ParseStreamCallback(reader, p.config, func(n *shared.ParserNode, err error) (stop bool, cbError error) {
		if err != nil {
			p.Errors <- err
			return true, err
		}
		p.Nodes <- n
		return false, nil
	})
	p.Done <- true
}

func getMetadataPair(line string) (*shared.MetadataPair, error) {
	trimmedLine := strings.TrimSpace(strings.Trim(line, "#"))
	separatorPos := strings.Index(trimmedLine, ":")
	if separatorPos > -1 {
		return &shared.MetadataPair{
			Name:  strings.Trim(trimmedLine[:separatorPos], "# \t"),
			Value: strings.TrimSpace(trimmedLine[separatorPos+1:]),
		}, nil
	}
	return &shared.MetadataPair{
		Name:  "",
		Value: trimmedLine,
	}, nil
}
