package analyzer

// TODO
// - vld.IsIdentifier -> t.T == tokenizer.TokenTypeIdentifier ?
// - check token type always, because symbols/identifiers etc may be strings (pass type to eat) ?
// - check why utils are not compiled to this module

import (
	pt "nand2tetris-golang/common/parsetree"
	"nand2tetris-golang/compiler/tokenizer"
	vld "nand2tetris-golang/compiler/validators"
	"strings"
)

// rule types
const (
	RuleTypeClass           = "class"
	RuleTypeKeyword         = "keyword"
	RuleTypeIdentifier      = "identifier"
	RuleTypeSymbol          = "symbol"
	RuleTypeClassVarDec     = "classVarDec"
	RuleTypeSubroutineDec   = "subroutineDec"
	RuleTypeParameterList   = "parameterList"
	RuleTypeSubroutineBody  = "subroutineBody"
	RuleTypeVarDec          = "varDec"
	RuleTypeStatements      = "statements"
	RuleTypeLetStatement    = "letStatement"
	RuleTypeDoStatement     = "doStatement"
	RuleTypeReturnStatement = "returnStatement"
	RuleTypeWhileStatement  = "whileStatement"
	RuleTypeIfStatement     = "ifStatement"
	RuleTypeExpression      = "expression"
	RuleTypeExpressionList  = "expressionList"
	RuleTypeTerm            = "term"
	RuleTypeIntegerConstant = "integerConstant"
	RuleTypeStringConstant  = "stringConstant"
)

// Analyzer type
type Analyzer struct {
	ti int                // current token index
	tl *[]tokenizer.Token // token list
}

// New returns new Analyzer
func New(tl *[]tokenizer.Token) *Analyzer {
	return &Analyzer{-1, tl}
}

// CompileClass ...
func (a *Analyzer) CompileClass() *pt.ParseTree {
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

func (a *Analyzer) compileClassVarDec() *pt.ParseTree {
	// Grammar: ('static' | 'field') type varName (',' varName)* ';'
	leaf := pt.New(RuleTypeClassVarDec, "")

	// ('static' | 'field')
	leaf.AddChildren(pt.New(RuleTypeKeyword, a.eat(vld.OneOf("static", "field"))))
	// type
	a.compileType(leaf, false)
	// varName (',' varName)*
	for {
		// varName
		leaf.AddChildren(pt.New(RuleTypeIdentifier, a.eat(vld.IsIdentifier)))
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

func (a *Analyzer) compileClassSubroutineDec() *pt.ParseTree {
	// Grammar: ("constructor" | "function" | "method") ('void' | type)
	// subroutineName '(' parameterList ')' subroutineBody
	leaf := pt.New(RuleTypeSubroutineDec, "")

	// ("constructor" | "function" | "method")
	leaf.AddChildren(pt.New(RuleTypeKeyword, a.eat(vld.OneOf("constructor", "function", "method"))))
	// ('void' | type)
	a.compileType(leaf, true)
	// subroutineName
	leaf.AddChildren(pt.New(RuleTypeIdentifier, a.eat(vld.IsIdentifier)))
	// '('
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity("("))))
	// parameterList
	leaf.AddChildren(a.compileParameterList())
	// ')'
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity(")"))))
	// subroutineBody
	leaf.AddChildren(a.compileSubroutineBody())

	return leaf
}

func (a *Analyzer) compileParameterList() *pt.ParseTree {
	// Grammar: ((type varName)(',' type varName)*)?
	leaf := pt.New(RuleTypeParameterList, "")

	// no parameters
	if vld.OneOf("int", "char", "boolean")(a.getNextToken().S) {
		for {
			// type
			a.compileType(leaf, false)
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

func (a *Analyzer) compileType(leaf *pt.ParseTree, includeVoid bool) {
	if a.getNextToken().T == tokenizer.TokenTypeKeyword {
		ops := []string{"int", "char", "boolean"}
		if includeVoid {
			ops = append(ops, "void")
		}
		leaf.AddChildren(pt.New(RuleTypeKeyword, a.eat(vld.OneOf(ops...))))
	} else {
		leaf.AddChildren(pt.New(RuleTypeIdentifier, a.eat(vld.IsIdentifier)))
	}
}

func (a *Analyzer) compileSubroutineBody() *pt.ParseTree {
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
	leaf.AddChildren(a.compileStatements())
	// '}'
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity("}"))))

	return leaf
}

func (a *Analyzer) compileVarDec() *pt.ParseTree {
	// Grammar: 'var' type varName (',' varName)* ';'
	leaf := pt.New(RuleTypeVarDec, "")

	// 'var'
	leaf.AddChildren(pt.New(RuleTypeKeyword, a.eat(vld.Identity("var"))))
	// type
	a.compileType(leaf, false)
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
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity(";"))))

	return leaf
}

func (a *Analyzer) compileStatements() *pt.ParseTree {
	// Grammar: letStatement | ifStatement | whileStatement | doStatement | return Statement
	leaf := pt.New(RuleTypeStatements, "")

	for {
		next := a.getNextToken()
		if !(next.T == tokenizer.TokenTypeKeyword && vld.OneOf("let", "if", "while", "do", "return")(next.S)) {
			break // no more var decs
		}
		switch next.S {
		case "let":
			leaf.AddChildren(a.compileLetStatement())
		case "if":
			leaf.AddChildren(a.compileIfStatement())
		case "while":
			leaf.AddChildren(a.compileWhileStatement())
		case "do":
			leaf.AddChildren(a.compileDoStatement())
		case "return":
			leaf.AddChildren(a.compileReturnStatement())
		}
	}

	return leaf
}

func (a *Analyzer) compileDoStatement() *pt.ParseTree {
	// Grammar: 'do' subroutineCall ';'
	leaf := pt.New(RuleTypeDoStatement, "")

	// 'do'
	leaf.AddChildren(pt.New(RuleTypeKeyword, a.eat(vld.Identity("do"))))
	// subroutineCall
	a.compileSubroutineCall(leaf)
	// ';'
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity(";"))))

	return leaf
}

func (a *Analyzer) compileReturnStatement() *pt.ParseTree {
	// Grammar: 'return' expression? ';'
	leaf := pt.New(RuleTypeReturnStatement, "")

	// 'return'
	leaf.AddChildren(pt.New(RuleTypeKeyword, a.eat(vld.Identity("return"))))
	// expression?
	if !vld.Identity(";")(a.getNextToken().S) {
		leaf.AddChildren(a.compileExpression())
	}
	// ';'
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity(";"))))

	return leaf
}

func (a *Analyzer) compileWhileStatement() *pt.ParseTree {
	// Grammar: 'while' '(' expression ')' '{' statements '}'
	leaf := pt.New(RuleTypeWhileStatement, "")

	// 'while'
	leaf.AddChildren(pt.New(RuleTypeKeyword, a.eat(vld.Identity("while"))))
	// '('
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity("("))))
	// expression
	leaf.AddChildren(a.compileExpression())
	// ')'
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity(")"))))
	// '{'
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity("{"))))
	// statements
	leaf.AddChildren(a.compileStatements())
	// '}'
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity("}"))))

	return leaf
}

func (a *Analyzer) compileIfStatement() *pt.ParseTree {
	// Grammar: 'if' '(' expression ')' '{' statements '}'
	// ('else' '{' statements '}')?
	leaf := pt.New(RuleTypeIfStatement, "")

	// 'if'
	leaf.AddChildren(pt.New(RuleTypeKeyword, a.eat(vld.Identity("if"))))
	// '('
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity("("))))
	// expression
	leaf.AddChildren(a.compileExpression())
	// ')'
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity(")"))))
	// '{'
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity("{"))))
	// statements
	leaf.AddChildren(a.compileStatements())
	// '}'
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity("}"))))
	// ('else' '{' statements '}')?
	if vld.Identity("else")(a.getNextToken().S) {
		// 'else'
		leaf.AddChildren(pt.New(RuleTypeKeyword, a.eat(vld.Identity("else"))))
		// '{'
		leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity("{"))))
		// statements
		leaf.AddChildren(a.compileStatements())
		// '}'
		leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity("}"))))
	}

	return leaf
}

func (a *Analyzer) compileLetStatement() *pt.ParseTree {
	// Grammar: 'let' varName ('['expression']')? '=' expression ';'
	leaf := pt.New(RuleTypeLetStatement, "")

	// 'let'
	leaf.AddChildren(pt.New(RuleTypeKeyword, a.eat(vld.Identity("let"))))
	// varName
	leaf.AddChildren(pt.New(RuleTypeIdentifier, a.eat(vld.IsIdentifier)))
	// ('['expression']')?
	if vld.Identity("[")(a.getNextToken().S) {
		// '['
		leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity("["))))
		// 'expression'
		leaf.AddChildren(a.compileExpression())
		// ']'
		leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity("]"))))
	}
	// '='
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity("="))))
	// expression
	leaf.AddChildren(a.compileExpression())
	// ';'
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity(";"))))

	return leaf
}

func (a *Analyzer) compileSubroutineCall(leaf *pt.ParseTree) {
	// Grammar: subroutineName '(' expressionList ')' |
	// (className|varName) '.' subroutineName '(' expressionList ')'

	// subroutineName | (className|varName)
	leaf.AddChildren(pt.New(RuleTypeIdentifier, a.eat(vld.IsIdentifier)))
	if vld.Identity(".")(a.getNextToken().S) {
		// '.' subroutineName
		// '.'
		leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity("."))))
		// subroutineName
		leaf.AddChildren(pt.New(RuleTypeIdentifier, a.eat(vld.IsIdentifier)))
	}
	// '('
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity("("))))
	// expressionList
	leaf.AddChildren(a.compileExpressionList())
	// ')'
	leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity(")"))))
}

func (a *Analyzer) compileExpression() *pt.ParseTree {
	// Grammar: term (op term)*
	leaf := pt.New(RuleTypeExpression, "")
	ops := []string{"+", "-", "*", "/", "&", "|", "<", ">", "="}

	for {
		leaf.AddChildren(a.compileTerm())
		if !vld.OneOf(ops...)(a.getNextToken().S) {
			break
		}
		leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.OneOf(ops...))))
	}

	return leaf
}

func (a *Analyzer) compileExpressionList() *pt.ParseTree {
	// Grammar: (expression (',' expression)*)?
	leaf := pt.New(RuleTypeExpressionList, "")

	for {
		// TODO verify this
		if vld.Identity(")")(a.getNextToken().S) {
			break
		}
		// expression
		leaf.AddChildren(a.compileExpression())
		if vld.Identity(",")(a.getNextToken().S) {
			// ','
			leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity(","))))
		}
	}

	return leaf
}

func (a *Analyzer) compileTerm() *pt.ParseTree {
	// Grammar: integetConstant | stringConstant | keywordConstant | varName |
	// varName '['expression']' | subroutineCall | '('expression')' | unaryOp term
	leaf := pt.New(RuleTypeTerm, "")
	next := a.getNextToken()
	afterNext := a.getTokenAfterNext()

	if next.T == tokenizer.TokenTypeIntConstant {
		// integetConstant
		leaf.AddChildren(pt.New(RuleTypeIntegerConstant, a.eat(vld.IsInt)))
	} else if next.T == tokenizer.TokenTypeStringConstant {
		// stringConstant
		leaf.AddChildren(pt.New(RuleTypeStringConstant, a.eat(vld.IsAny())))
	} else if next.T == tokenizer.TokenTypeKeyword && vld.OneOf("true", "false", "null", "this")(next.S) {
		// keywordConstant
		leaf.AddChildren(pt.New(RuleTypeKeyword, a.eat(vld.OneOf("true", "false", "null", "this"))))
	} else if next.T == tokenizer.TokenTypeSymbol && next.S == "(" {
		// '('expression')'
		leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity("("))))
		leaf.AddChildren(a.compileExpression())
		leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity(")"))))
	} else if next.T == tokenizer.TokenTypeSymbol && vld.OneOf("-", "~")(next.S) {
		// unaryOp term
		leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.OneOf("-", "~"))))
		leaf.AddChildren(a.compileTerm())
	} else if afterNext.T == tokenizer.TokenTypeSymbol && vld.OneOf("(", ".")(afterNext.S) {
		// subroutineCall
		a.compileSubroutineCall(leaf)
	} else {
		// varName
		leaf.AddChildren(pt.New(RuleTypeIdentifier, a.eat(vld.IsIdentifier)))
		// '['expression']' ?
		next := a.getNextToken()
		if next.T == tokenizer.TokenTypeSymbol && next.S == "[" {
			leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity("["))))
			leaf.AddChildren(a.compileExpression())
			leaf.AddChildren(pt.New(RuleTypeSymbol, a.eat(vld.Identity("]"))))
		}
	}

	return leaf
}

func (a *Analyzer) getCurrentToken() tokenizer.Token {
	t := *a.tl
	return t[a.ti]
}

func (a *Analyzer) getNextToken() tokenizer.Token {
	t := *a.tl
	return t[a.ti+1]
}

func (a *Analyzer) getTokenAfterNext() tokenizer.Token {
	t := *a.tl
	return t[a.ti+2]
}

// increment ti, return currentToken if valid, panic if not
// one of provided rules should pass
func (a *Analyzer) eat(rules ...vld.Rule) string {
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
func ToXML(tree *pt.ParseTree, indent int) string {
	tab := strings.Repeat("  ", indent)
	xml := ""
	leafType := tree.Type()
	val := tree.Value()
	children := tree.Children()
	open := "<" + leafType + ">"
	close := "</" + leafType + ">"
	childrenLess := []string{
		RuleTypeKeyword, RuleTypeIdentifier, RuleTypeSymbol,
		RuleTypeIntegerConstant, RuleTypeStringConstant,
	}

	if vld.OneOf(childrenLess...)(leafType) {
		xml += tab + open + " " + normalizeSymbol(val) + " " + close + "\n"
	} else {
		xml += tab + open + "\n"
		for _, leaf := range children {
			xml += ToXML(leaf, indent+1)
		}
		xml += tab + close + "\n"
	}

	return xml
}

func normalizeSymbol(s string) string {
	m := map[string]string{
		"<":  "&lt;",
		">":  "&gt;",
		"&":  "&amp;",
		"\"": "&quot;",
	}
	if v, in := m[s]; in {
		return v
	}
	return s
}
