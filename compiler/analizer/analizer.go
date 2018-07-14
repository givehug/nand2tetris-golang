package analizer

import (
	pt "nand2tetris-golang/common/parsetree"
	"nand2tetris-golang/compiler/tokenizer"
	vld "nand2tetris-golang/compiler/validators"
	"strings"
)

// rule types
const (
	RuleTypeClass          = "class"
	RuleTypeKeyword        = "keyword"
	RuleTypeIdentifier     = "identifier"
	RuleTypeSymbol         = "symbol"
	RuleTypeClassVarDec    = "classVarDec"
	RuleTypeSubroutineDec  = "subroutineDec"
	RuleTypeParameterList  = "parameterList"
	RuleTypeSubroutineBody = "subroutineBody"
	RuleTypeVarDec         = "varDec"
	RuleTypeStatements     = "statements"
)

// Analizer type
type Analizer struct {
	ti int                // current token index
	tl *[]tokenizer.Token // token list
}

// New returns new Analizer
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
		if !vld.OneOf("static", "field")(a.getNextToken().S) {
			break // no more var decs
		}
		addIfHasChildren(leaf, a.compileClassVarDec())
	}
	// subroutineDec*
	for {
		if !vld.OneOf("constructor", "function", "method")(a.getNextToken().S) {
			break // no more subroutines
		}
		addIfHasChildren(leaf, a.compileClassSubroutineDec())
	}
	// '}'
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity("}"))))

	return leaf
}

func (a *Analizer) compileClassVarDec() *pt.ParseTree {
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
		if !vld.Identity(",")(a.getNextToken().S) {
			break // no more identifiers
		}
		// ','
		leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity(","))))
	}
	// ';'
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity(";"))))

	return leaf
}

func (a *Analizer) compileClassSubroutineDec() *pt.ParseTree {
	// Grammar: ("constructor" | "function" | "method") ('void' | type)
	// subroutineName '(' parameterList ')' subroutineBody
	leaf := pt.New(RuleTypeSubroutineDec, "")

	// ("constructor" | "function" | "method")
	leaf.AddChildren(pt.New(RuleTypeKeyword, a.eat(vld.OneOf("constructor", "function", "method"))))
	// ('void' | type)
	leaf.AddChildren(pt.New(RuleTypeKeyword, a.eat(
		vld.OneOf("int", "char", "boolean", "void"),
		vld.IsIdentifier,
	)))
	// subroutineName
	leaf.AddChildren(pt.New(RuleTypeIdentifier, a.eat(vld.IsIdentifier)))
	// '('
	leaf.AddChildren(pt.New(RuleTypeIdentifier, a.eat(vld.Identity("("))))
	// parameterList
	leaf.AddChildren(a.compileParameterList())
	// ')'
	leaf.AddChildren(pt.New(RuleTypeIdentifier, a.eat(vld.Identity(")"))))
	// subroutineBody
	leaf.AddChildren(a.compileSubroutineBody())

	return leaf
}

func (a *Analizer) compileParameterList() *pt.ParseTree {
	// Grammar: ((type varName)(',' type varName)*)?
	leaf := pt.New(RuleTypeParameterList, "")

	// no parameters
	if vld.OneOf("int", "char", "boolean")(a.getNextToken().S) {
		for {
			// type
			leaf.AddChildren(pt.New(RuleTypeKeyword, a.eat(
				vld.OneOf("int", "char", "boolean"),
				vld.IsIdentifier,
			)))
			// varName
			leaf.AddChildren(pt.New(RuleTypeIdentifier, a.eat(vld.IsIdentifier)))
			if !vld.Identity(",")(a.getNextToken().S) {
				break // no more parameters
			}
			// ','
			leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity(","))))
		}
	}

	return leaf
}

func (a *Analizer) compileSubroutineBody() *pt.ParseTree {
	// Grammar: '{' varDec* statements '}'
	leaf := pt.New(RuleTypeSubroutineBody, "")

	// '{'
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity("{"))))
	// varDec*
	for {
		if !vld.Identity("var")(a.getNextToken().S) {
			break // no more var decs
		}
		leaf.AddChildren(a.compileVarDec())
	}
	// statements
	for {
		if !vld.OneOf("let", "if", "while", "do", "return")(a.getNextToken().S) {
			break // no more var decs
		}
		leaf.AddChildren(a.compileStatements())
	}
	// '}'
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity("}"))))

	return leaf
}

// TODO
func (a *Analizer) compileStatements() *pt.ParseTree {
	// Grammar: statement*
	leaf := pt.New(RuleTypeStatements, "")

	leaf.AddChildren(pt.New(RuleTypeKeyword, a.eat(vld.Identity("let"))))

	return leaf
}

func (a *Analizer) compileVarDec() *pt.ParseTree {
	// Grammar: 'var' type varName (',' varName)* ';'
	leaf := pt.New(RuleTypeVarDec, "")

	// 'var'
	leaf.AddChildren(pt.New(RuleTypeKeyword, a.eat(vld.Identity("var"))))
	// type
	leaf.AddChildren(pt.New(RuleTypeKeyword, a.eat(
		vld.OneOf("int", "char", "boolean"),
		vld.IsIdentifier,
	)))
	// varName (',' varName)*
	for {
		// varName
		leaf.AddChildren(pt.New(RuleTypeIdentifier, a.eat(vld.IsIdentifier)))
		if !vld.Identity(",")(a.getNextToken().S) {
			break // no more var identifiers
		}
		// ','
		leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity(","))))
	}
	// ';'
	leaf.AddChildren(pt.New(RuleTypeKeyword, a.eat(vld.Identity(";"))))

	return leaf
}

func (a *Analizer) getCurrentToken() tokenizer.Token {
	t := *a.tl
	return t[a.ti]
}

func (a *Analizer) getNextToken() tokenizer.Token {
	t := *a.tl
	return t[a.ti+1]
}

// increment ti, return currentToken if valid, panic if not
// one of provided rules should pass
func (a *Analizer) eat(rules ...vld.Rule) string {
	a.ti++
	s := a.getCurrentToken().S
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
