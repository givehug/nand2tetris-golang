package analizer

import (
	"nand2tetris-golang/common/parsetree"
	"nand2tetris-golang/compiler/tokenizer"
	"nand2tetris-golang/compiler/validators"
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
func (a *Analizer) eat(rule validators.Rule) string {
	*a.ti++
	s := a.currentToken().S
	return validators.Validate(s, rule)
}

// CompileClass ...
func (a *Analizer) CompileClass() *parsetree.ParseTree {
	// class keyword
	a.pt.AddChildren(parsetree.New(RuleTypeKeyword, a.eat(validators.Identity("class"))))
	// class name
	a.pt.AddChildren(parsetree.New(RuleTypeIdentifier, a.eat(validators.IsIdentifier)))
	// open curl
	a.pt.AddChildren(parsetree.New(RuleTypeSymbol, a.eat(validators.Identity("{"))))
	// todo get class var dec
	// todo get class subroutine dec
	// close curl
	a.pt.AddChildren(parsetree.New(RuleTypeSymbol, a.eat(validators.Identity("}"))))

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
