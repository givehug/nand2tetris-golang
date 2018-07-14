package analizer

import (
	"nand2tetris-golang/common/parsetree"
	"nand2tetris-golang/compiler/tokenizer"
	"strings"
)

// rule types
const (
	RuleTypeClass      = "class"
	RuleTypeKeyword    = "keyword"
	RuleTypeIdentifier = "identifier"
	RuleTypeSymbol     = "symbol"
)

// Analizer type
type Analizer struct {
	pt *parsetree.ParseTree // parse tree
	ti *int                 // current token index
	tl *[]tokenizer.Token   // token list
}

// New constructs new Analizer
func New(tl *[]tokenizer.Token) *Analizer {
	ti := -1
	pt := parsetree.New(RuleTypeClass, "")
	return &Analizer{pt, &ti, tl}
}

// get current token value
func (a *Analizer) currentToken() tokenizer.Token {
	t := *a.tl
	return t[*a.ti]
}

// increment ti, check if rule is met
func (a *Analizer) eat(rule string) {
	*a.ti++
	t := a.currentToken()
	if t.S != rule {
		panic("bad eat: " + rule + " - " + t.S + " - " + t.T)
	}
}

// CompileClass ...
func (a *Analizer) CompileClass() *parsetree.ParseTree {
	a.eat("class")
	a.pt.AddChildren(parsetree.New(RuleTypeKeyword, "class"))
	// todo get class name
	a.eat("{")
	a.pt.AddChildren(parsetree.New(RuleTypeSymbol, "{"))
	// todo get class var dec
	// todo get class subroutine dec
	a.eat("}")
	a.pt.AddChildren(parsetree.New(RuleTypeSymbol, "}"))

	return a.pt
}

// ToXML ...
func ToXML(pt *parsetree.ParseTree, indent int) string {
	tab := strings.Repeat("  ", indent)
	xml := ""
	t := pt.Type()
	v := pt.Value()
	c := pt.Children()
	open := "<" + t + ">"
	close := "</" + t + ">"

	if len(c) == 0 {
		xml += tab + open + " " + v + " " + close + "\n"
	} else {
		xml += tab + open + "\n"
		for _, ch := range c {
			xml += ToXML(ch, indent+1)
		}
		xml += tab + close + "\n"
	}

	return xml
}
