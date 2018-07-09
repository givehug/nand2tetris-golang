package parser

import (
	"bufio"
	"nand2tetris-golang/common/utils"
	"os"
	"strings"
)

const eol = "\n"

// Parser struct
type Parser struct {
	file    *os.File
	scanner *bufio.Scanner
}

// New creates new Parser
func New(path string) *Parser {
	f, err := os.Open(path)
	utils.HandleErr(err)
	reader := bufio.NewReader(f)
	scanner := bufio.NewScanner(reader)
	return &Parser{f, scanner}
}

// Parse reads next file line
// returns:
// - line text content, excluding white spaces and comments
// - boolean, true if file has more unparsed lines
func (p *Parser) Parse() (string, bool) {
	s := "" // next line text content
	ok := p.scanner.Scan()
	if !ok {
		return s, false // error or end of file
	}
	s = p.scanner.Text() // get next line text
	// trim comments and whitespaces
	if comment := strings.Index(s, "//"); comment > -1 {
		s = strings.TrimSpace(s[:comment])
	} else {
		s = strings.TrimSpace(s)
	}
	// if line is whitespace/comment, prase again
	if len(s) == 0 {
		return p.Parse()
	}
	return s, true
}

// Close file
func (p *Parser) Close() {
	p.file.Close()
}
