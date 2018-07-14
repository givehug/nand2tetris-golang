package analizer

import (
	pt "nand2tetris-golang/common/parsetree"
	"nand2tetris-golang/compiler/tokenizer"
	vld "nand2tetris-golang/compiler/validators"
	"strings"
)

// rule types
const (
	RuleTypeClass         = "class"
	RuleTypeKeyword       = "keyword"
	RuleTypeIdentifier    = "identifier"
	RuleTypeSymbol        = "symbol"
	RuleTypeClassVarDec   = "classVarDec"
	RuleTypeSubroutineDec = "subroutineDec"
)

// Analizer type
type Analizer struct {
	ti int                // current token index
	tl *[]tokenizer.Token // token list
}

// New constructs new Analizer
func New(tl *[]tokenizer.Token) *Analizer {
	return &Analizer{-1, tl}
}

// CompileClass ...
func (a *Analizer) CompileClass() *pt.ParseTree {
	// Grammar: 'class' className '{' classVarDec* subroutineDec* '}'
	leaf := pt.New(RuleTypeClass, "")

	// 'class'
	leaf.AddChildren(pt.New(RuleTypeKeyword, a.eat(vld.Identity("class"))))
	// className
	leaf.AddChildren(pt.New(RuleTypeIdentifier, a.eat(vld.IsIdentifier)))
	// '{'
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity("{"))))
	// classVarDec*
	for {
		if !vld.OneOf("static", "field")(a.nextToken().S) {
			break // no more var decs
		}
		addIfHasChildren(leaf, a.CompileClassVarDec())
	}
	// subroutineDec*
	for {
		if !vld.OneOf("constructor", "function", "method")(a.nextToken().S) {
			break // no more subroutines
		}
		addIfHasChildren(leaf, a.CompileClassSubroutineDec())
	}
	// '}'
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity("}"))))

	return leaf
}

// CompileClassVarDec ...
func (a *Analizer) CompileClassVarDec() *pt.ParseTree {
	// Grammar: ('static' | 'field') type varName (',' varName)* ';'
	leaf := pt.New(RuleTypeClassVarDec, "")

	// ('static' | 'field')
	leaf.AddChildren(pt.New(RuleTypeKeyword, a.eat(vld.OneOf("static", "field"))))
	// type
	leaf.AddChildren(pt.New(RuleTypeKeyword, a.eat(
		vld.OneOf("int", "char", "boolean"),
		vld.IsIdentifier,
	)))
	// varName (',' varName)*
	for {
		// varName
		leaf.AddChildren(pt.New(RuleTypeKeyword, a.eat(vld.IsIdentifier)))
		if !vld.Identity(",")(a.nextToken().S) {
			break // no more identifiers
		}
		// ','
		leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity(","))))
	}
	// ';'
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity(";"))))

	return leaf
}

// CompileClassSubroutineDec ... TODO
func (a *Analizer) CompileClassSubroutineDec() *pt.ParseTree {
	leaf := pt.New(RuleTypeSubroutineDec, "")
	return leaf
}

// get current token value
func (a *Analizer) currentToken() tokenizer.Token {
	t := *a.tl
	return t[a.ti]
}

// get next token
func (a *Analizer) nextToken() tokenizer.Token {
	t := *a.tl
	return t[a.ti+1]
}

// increment ti, return currentToken if valid, panic if not
// one of provided rules should pass
func (a *Analizer) eat(rules ...vld.Rule) string {
	a.ti++
	s := a.currentToken().S
	for _, r := range rules {
		if r(s) {
			return s
		}
	}
	panic("Rule not met for: " + s)
}

// add 'leaf' to 'to' as child if 'leaf' has children
func addIfHasChildren(to, leaf *pt.ParseTree) {
	if leaf.HasChildren() {
		to.AddChildren(leaf)
	}
}

// ToXML ...
func ToXML(pt *pt.ParseTree, indent int) string {
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

func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}
