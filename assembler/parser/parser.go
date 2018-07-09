package parser

import (
	p "nand2tetris-golang/common/parser"
	"strconv"
	"strings"
)

// Line types
const (
	CmdTypeA = "a-instruction" // A instruction
	CmdTypeC = "c-instruction" // C instruction
	CmdTypeL = "label"         // label
)

// New returns new Parser
func New(sourceFile string) *p.Parser {
	return p.New(sourceFile)
}

// CommandType returns type of line constant
func CommandType(c string) string {
	switch {
	case strings.HasPrefix(c, "(") && strings.HasSuffix(c, ")"):
		return CmdTypeL // label lines
	case strings.HasPrefix(c, "@"):
		return CmdTypeA // a-instruction
	default:
		return CmdTypeC // c-instruction
	}
}

// CommandArgs returns command args
// (dest, comp, jump) in case of c-instruction
func CommandArgs(s string) (a string, b string, c string) {
	a, b, c = "", "", ""
	switch CommandType(s) {
	case CmdTypeL:
		a = s[1 : len(s)-1]
	case CmdTypeA:
		a = s[1:]
	case CmdTypeC:
		compInd := strings.Index(s, "=")
		jumpInd := strings.Index(s, ";")
		if jumpInd != -1 {
			c = s[jumpInd+1:]
		} else {
			jumpInd = len(s)
		}
		if compInd == -1 {
			b = s[:jumpInd]
		} else {
			a = s[:compInd]
			b = s[compInd+1 : jumpInd]
		}
	}
	return
}

// IsVariable tests if string is variable or value
// returns true if variable
func IsVariable(s string) bool {
	_, err := strconv.ParseInt(s, 10, 16)
	return err != nil
}
