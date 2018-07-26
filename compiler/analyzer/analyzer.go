package analyzer

// TODO
// - vld.IsIdentifier -> t.T == tokenizer.TokenTypeIdentifier ?
// - check token type always, because symbols/identifiers etc may be strings (pass type to eat) ?

import (
	pt "nand2tetris-golang/common/parsetree"
	"nand2tetris-golang/common/utils"
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

// CompileClass accepts token list and returns parse tree
func CompileClass(tl *[]tokenizer.Token) *pt.ParseTree {
	// create analyzer struct
	a := &Analyzer{-1, tl}

	// Grammar: 'class' className '{' classVarDec* subroutineDec* '}'
	tree := pt.New(RuleTypeClass, "")

	// 'class'
	tree.AddLeaves(pt.New(RuleTypeKeyword, a.eat(vld.Identity("class"))))
	// className
	tree.AddLeaves(pt.New(RuleTypeIdentifier, a.eat(vld.IsIdentifier)))
	// '{'
	tree.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity("{"))))
	// classVarDec*
	for {
		if !vld.OneOf("static", "field")(a.getNextToken().S) {
			break // no more var decs
		}
		addIfHasChildren(tree, a.compileClassVarDec())
	}
	// subroutineDec*
	for {
		if !vld.OneOf("constructor", "function", "method")(a.getNextToken().S) {
			break // no more subroutines
		}
		addIfHasChildren(tree, a.compileClassSubroutineDec())
	}
	// '}'
	tree.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity("}"))))

	return tree
}

func (a *Analyzer) compileClassVarDec() *pt.ParseTree {
	// Grammar: ('static' | 'field') type varName (',' varName)* ';'
	leaf := pt.New(RuleTypeClassVarDec, "")

	// ('static' | 'field')
	leaf.AddLeaves(pt.New(RuleTypeKeyword, a.eat(vld.OneOf("static", "field"))))
	// type
	a.compileType(leaf, false)
	// varName (',' varName)*
	for {
		// varName
		leaf.AddLeaves(pt.New(RuleTypeIdentifier, a.eat(vld.IsIdentifier)))
		if !vld.Identity(",")(a.getNextToken().S) {
			break // no more identifiers
		}
		// ','
		leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity(","))))
	}
	// ';'
	leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity(";"))))

	return leaf
}

func (a *Analyzer) compileClassSubroutineDec() *pt.ParseTree {
	// Grammar: ("constructor" | "function" | "method") ('void' | type)
	// subroutineName '(' parameterList ')' subroutineBody
	leaf := pt.New(RuleTypeSubroutineDec, "")

	// ("constructor" | "function" | "method")
	leaf.AddLeaves(pt.New(RuleTypeKeyword, a.eat(vld.OneOf("constructor", "function", "method"))))
	// ('void' | type)
	a.compileType(leaf, true)
	// subroutineName
	leaf.AddLeaves(pt.New(RuleTypeIdentifier, a.eat(vld.IsIdentifier)))
	// '('
	leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity("("))))
	// parameterList
	leaf.AddLeaves(a.compileParameterList())
	// ')'
	leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity(")"))))
	// subroutineBody
	leaf.AddLeaves(a.compileSubroutineBody())

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
			leaf.AddLeaves(pt.New(RuleTypeIdentifier, a.eat(vld.IsIdentifier)))
			if !vld.Identity(",")(a.getNextToken().S) {
				break // no more parameters
			}
			// ','
			leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity(","))))
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
		leaf.AddLeaves(pt.New(RuleTypeKeyword, a.eat(vld.OneOf(ops...))))
	} else {
		leaf.AddLeaves(pt.New(RuleTypeIdentifier, a.eat(vld.IsIdentifier)))
	}
}

func (a *Analyzer) compileSubroutineBody() *pt.ParseTree {
	// Grammar: '{' varDec* statements '}'
	leaf := pt.New(RuleTypeSubroutineBody, "")

	// '{'
	leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity("{"))))
	// varDec*
	for {
		if !vld.Identity("var")(a.getNextToken().S) {
			break // no more var decs
		}
		leaf.AddLeaves(a.compileVarDec())
	}
	// statements
	leaf.AddLeaves(a.compileStatements())
	// '}'
	leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity("}"))))

	return leaf
}

func (a *Analyzer) compileVarDec() *pt.ParseTree {
	// Grammar: 'var' type varName (',' varName)* ';'
	leaf := pt.New(RuleTypeVarDec, "")

	// 'var'
	leaf.AddLeaves(pt.New(RuleTypeKeyword, a.eat(vld.Identity("var"))))
	// type
	a.compileType(leaf, false)
	// varName (',' varName)*
	for {
		// varName
		leaf.AddLeaves(pt.New(RuleTypeIdentifier, a.eat(vld.IsIdentifier)))
		if !vld.Identity(",")(a.getNextToken().S) {
			break // no more var identifiers
		}
		// ','
		leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity(","))))
	}
	// ';'
	leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity(";"))))

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
			leaf.AddLeaves(a.compileLetStatement())
		case "if":
			leaf.AddLeaves(a.compileIfStatement())
		case "while":
			leaf.AddLeaves(a.compileWhileStatement())
		case "do":
			leaf.AddLeaves(a.compileDoStatement())
		case "return":
			leaf.AddLeaves(a.compileReturnStatement())
		}
	}

	return leaf
}

func (a *Analyzer) compileDoStatement() *pt.ParseTree {
	// Grammar: 'do' subroutineCall ';'
	leaf := pt.New(RuleTypeDoStatement, "")

	// 'do'
	leaf.AddLeaves(pt.New(RuleTypeKeyword, a.eat(vld.Identity("do"))))
	// subroutineCall
	a.compileSubroutineCall(leaf)
	// ';'
	leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity(";"))))

	return leaf
}

func (a *Analyzer) compileReturnStatement() *pt.ParseTree {
	// Grammar: 'return' expression? ';'
	leaf := pt.New(RuleTypeReturnStatement, "")

	// 'return'
	leaf.AddLeaves(pt.New(RuleTypeKeyword, a.eat(vld.Identity("return"))))
	// expression?
	if !vld.Identity(";")(a.getNextToken().S) {
		leaf.AddLeaves(a.compileExpression())
	}
	// ';'
	leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity(";"))))

	return leaf
}

func (a *Analyzer) compileWhileStatement() *pt.ParseTree {
	// Grammar: 'while' '(' expression ')' '{' statements '}'
	leaf := pt.New(RuleTypeWhileStatement, "")

	// 'while'
	leaf.AddLeaves(pt.New(RuleTypeKeyword, a.eat(vld.Identity("while"))))
	// '('
	leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity("("))))
	// expression
	leaf.AddLeaves(a.compileExpression())
	// ')'
	leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity(")"))))
	// '{'
	leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity("{"))))
	// statements
	leaf.AddLeaves(a.compileStatements())
	// '}'
	leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity("}"))))

	return leaf
}

func (a *Analyzer) compileIfStatement() *pt.ParseTree {
	// Grammar: 'if' '(' expression ')' '{' statements '}'
	// ('else' '{' statements '}')?
	leaf := pt.New(RuleTypeIfStatement, "")

	// 'if'
	leaf.AddLeaves(pt.New(RuleTypeKeyword, a.eat(vld.Identity("if"))))
	// '('
	leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity("("))))
	// expression
	leaf.AddLeaves(a.compileExpression())
	// ')'
	leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity(")"))))
	// '{'
	leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity("{"))))
	// statements
	leaf.AddLeaves(a.compileStatements())
	// '}'
	leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity("}"))))
	// ('else' '{' statements '}')?
	if vld.Identity("else")(a.getNextToken().S) {
		// 'else'
		leaf.AddLeaves(pt.New(RuleTypeKeyword, a.eat(vld.Identity("else"))))
		// '{'
		leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity("{"))))
		// statements
		leaf.AddLeaves(a.compileStatements())
		// '}'
		leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity("}"))))
	}

	return leaf
}

func (a *Analyzer) compileLetStatement() *pt.ParseTree {
	// Grammar: 'let' varName ('['expression']')? '=' expression ';'
	leaf := pt.New(RuleTypeLetStatement, "")

	// 'let'
	leaf.AddLeaves(pt.New(RuleTypeKeyword, a.eat(vld.Identity("let"))))
	// varName
	leaf.AddLeaves(pt.New(RuleTypeIdentifier, a.eat(vld.IsIdentifier)))
	// ('['expression']')?
	if vld.Identity("[")(a.getNextToken().S) {
		// '['
		leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity("["))))
		// 'expression'
		leaf.AddLeaves(a.compileExpression())
		// ']'
		leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity("]"))))
	}
	// '='
	leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity("="))))
	// expression
	leaf.AddLeaves(a.compileExpression())
	// ';'
	leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity(";"))))

	return leaf
}

func (a *Analyzer) compileSubroutineCall(leaf *pt.ParseTree) {
	// Grammar: subroutineName '(' expressionList ')' |
	// (className|varName) '.' subroutineName '(' expressionList ')'

	// subroutineName | (className|varName)
	leaf.AddLeaves(pt.New(RuleTypeIdentifier, a.eat(vld.IsIdentifier)))
	if vld.Identity(".")(a.getNextToken().S) {
		// '.' subroutineName
		// '.'
		leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity("."))))
		// subroutineName
		leaf.AddLeaves(pt.New(RuleTypeIdentifier, a.eat(vld.IsIdentifier)))
	}
	// '('
	leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity("("))))
	// expressionList
	leaf.AddLeaves(a.compileExpressionList())
	// ')'
	leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity(")"))))
}

func (a *Analyzer) compileExpression() *pt.ParseTree {
	// Grammar: term (op term)*
	leaf := pt.New(RuleTypeExpression, "")
	ops := []string{"+", "-", "*", "/", "&", "|", "<", ">", "="}

	for {
		leaf.AddLeaves(a.compileTerm())
		if !vld.OneOf(ops...)(a.getNextToken().S) {
			break
		}
		leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.OneOf(ops...))))
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
		leaf.AddLeaves(a.compileExpression())
		if vld.Identity(",")(a.getNextToken().S) {
			// ','
			leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity(","))))
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
		leaf.AddLeaves(pt.New(RuleTypeIntegerConstant, a.eat(vld.IsInt)))
	} else if next.T == tokenizer.TokenTypeStringConstant {
		// stringConstant
		leaf.AddLeaves(pt.New(RuleTypeStringConstant, a.eat(vld.IsAny())))
	} else if next.T == tokenizer.TokenTypeKeyword && vld.OneOf("true", "false", "null", "this")(next.S) {
		// keywordConstant
		leaf.AddLeaves(pt.New(RuleTypeKeyword, a.eat(vld.OneOf("true", "false", "null", "this"))))
	} else if next.T == tokenizer.TokenTypeSymbol && next.S == "(" {
		// '('expression')'
		leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity("("))))
		leaf.AddLeaves(a.compileExpression())
		leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity(")"))))
	} else if next.T == tokenizer.TokenTypeSymbol && vld.OneOf("-", "~")(next.S) {
		// unaryOp term
		leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.OneOf("-", "~"))))
		leaf.AddLeaves(a.compileTerm())
	} else if afterNext.T == tokenizer.TokenTypeSymbol && vld.OneOf("(", ".")(afterNext.S) {
		// subroutineCall
		a.compileSubroutineCall(leaf)
	} else {
		// varName
		leaf.AddLeaves(pt.New(RuleTypeIdentifier, a.eat(vld.IsIdentifier)))
		// '['expression']' ?
		next := a.getNextToken()
		if next.T == tokenizer.TokenTypeSymbol && next.S == "[" {
			leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity("["))))
			leaf.AddLeaves(a.compileExpression())
			leaf.AddLeaves(pt.New(RuleTypeSymbol, a.eat(vld.Identity("]"))))
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
	if leaf.HasLeaves() {
		to.AddLeaves(leaf)
	}
}

// ToXML returns xml representation of parse tree
func ToXML(tree *pt.ParseTree, indent int) string {
	tab := strings.Repeat("  ", indent)
	childrenLess := []string{
		RuleTypeKeyword, RuleTypeIdentifier, RuleTypeSymbol,
		RuleTypeIntegerConstant, RuleTypeStringConstant,
	}

	val := ""
	if vld.OneOf(childrenLess...)(tree.Type()) {
		val += " " + tree.Value() + " "
	} else {
		val += "\n"
		for _, leaf := range tree.Leaves() {
			val += ToXML(leaf, indent+1)
		}
		val += tab
	}

	return tab + utils.ToXMLTag(tree.Type(), val) + "\n"
}
