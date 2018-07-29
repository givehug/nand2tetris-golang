package engine

// TODO check token type always, symbols/identifiers etc could be strings (pass type to eat) ?

import (
	"nand2tetris-golang/common/utils"
	"nand2tetris-golang/compiler/mapping"
	st "nand2tetris-golang/compiler/symboltable"
	"nand2tetris-golang/compiler/tokenizer"
	vld "nand2tetris-golang/compiler/validators"
	vm "nand2tetris-golang/compiler/vmwriter"
	"os"
	"strconv"
)

// Engine type
type Engine struct {
	tokens         *tokenizer.Tokens // token list
	table          *st.SymbolTable   // symbol table
	className      string            // current class name
	subroutineName string            // current subroutine name
	subroutineKind string            // current subroutine name
	outFile        *os.File          // .vm out file
	labelCounter   int               // unique label name counter
}

// Compile compiles jack to vm file
func Compile(inFile, outFile string) {
	file, err := os.Create(outFile)                  // create out .vm file
	defer file.Close()                               // close file before return
	utils.HandleErr(err)                             // handle possible error
	tokens := tokenizer.GetTokens(inFile)            // get tokens list
	table := st.New()                                // create symbol table
	a := &Engine{tokens, table, "", "", "", file, 0} // create Engine
	a.compileClass()                                 // start compiling class

	// tree := analyzer.CompileClass(&tokens) // create parse tree
	// fmt.Println(tokenizer.ToXML(tokens)) // print tokens XML
	// fmt.Println(analyzer.ToXML(tree, 0)) // print parsed tree XML
}

func (e *Engine) compileClass() {
	// Grammar: 'class' className '{' classVarDec* subroutineDec* '}'
	e.eat(vld.Identity("class"))                     // 'class'
	e.className = e.currentToken().S                 // save class name
	e.eatType(tokenizer.TokenTypeIdentifier)         // className
	vm.WriteComment(e.outFile, "Class "+e.className) // write comment
	e.eat(vld.Identity("{"))                         // '{'
	e.compileClassVarDec()                           // classVarDec*
	e.compileClassSubroutineDec()                    // subroutineDec*
	e.eat(vld.Identity("}"))                         // '}'
}

func (e *Engine) compileClassVarDec() {
	// Grammar: ('static' | 'field') type varName (',' varName)* ';'
	if !vld.OneOf("static", "field")(e.currentToken().S) {
		return // no more var decs
	}
	kind := e.eat(vld.OneOf("static", "field")) // ('static' | 'field')
	varType := e.compileType(false)             // var type
	// varName (',' varName)*
	for {
		name := e.eatType(tokenizer.TokenTypeIdentifier) // var name
		e.table.Define(name, varType, kind)              // add to symbol table
		if e.currentToken().S != "," {
			break // no more identifiers
		}
		e.eat(vld.Identity(",")) // ','
	}
	e.eat(vld.Identity(";")) // ';'
	e.compileClassVarDec()
}

func (e *Engine) compileClassSubroutineDec() {
	// Grammar: ("constructor" | "function" | "method") ('void' | type)
	// subroutineName '(' parameterList ')' subroutineBody
	e.table.StartSubroutine()
	e.subroutineKind = e.eat(vld.OneOf("constructor", "function", "method")) // ("constructor" | "function" | "method")
	e.compileType(true)                                                      // ('void' | type)
	e.subroutineName = e.eatType(tokenizer.TokenTypeIdentifier)              // subroutineName
	e.eat(vld.Identity("("))                                                 // '('
	e.compileParameterList()                                                 // parameterList
	e.eat(vld.Identity(")"))                                                 // ')'
	e.compileSubroutineBody()                                                // subroutineBody
	// check if has more subroutines
	if e.currentToken().S != "}" {
		e.compileClassSubroutineDec()
	}
}

func (e *Engine) compileParameterList() {
	// Grammar: ((type varName)(',' type varName)*)?
	// compile this argument
	if e.subroutineKind == "method" {
		e.table.Define("this", e.className, mapping.IdentifierTypeArg)
	}
	// check if has any params
	if e.currentToken().S == ")" {
		return
	}
	for {
		// add argument to symbol table
		varType := e.compileType(false)                     // type
		varName := e.eatType(tokenizer.TokenTypeIdentifier) // varName
		e.table.Define(varName, varType, mapping.IdentifierTypeArg)
		// check if has more params
		if e.currentToken().S == "," {
			e.eat(vld.Identity(",")) // ','
		} else {
			break
		}
	}
}

func (e *Engine) compileType(includeVoid bool) string {
	if e.currentToken().T == tokenizer.TokenTypeKeyword {
		ops := []string{"int", "char", "boolean"}
		if includeVoid {
			ops = append(ops, "void")
		}
		return e.eat(vld.OneOf(ops...))
	}
	return e.eatType(tokenizer.TokenTypeIdentifier)
}

func (e *Engine) compileSubroutineBody() {
	// Grammar: '{' varDec* statements '}'
	e.eat(vld.Identity("{")) // '{'
	// varDec*
	for {
		if e.currentToken().S != "var" {
			break // no more var decs
		}
		e.compileVarDec()
	}
	// write function
	vm.WriteComment(e.outFile, "Subroutine "+e.subroutineKind+" "+e.subroutineName)
	vm.WriteFunction(e.outFile, e.className, e.subroutineName, e.table.VarCount("local"))
	if e.subroutineKind == "constructor" {
		vm.WritePush(e.outFile, mapping.SegmentCONST, e.table.VarCount("field"))
		vm.WriteCall(e.outFile, "Memory", "alloc", 1)
		vm.WritePop(e.outFile, mapping.SegmentPOINT, 0)
	}
	if e.subroutineKind == "method" {
		vm.WritePush(e.outFile, mapping.SegmentARG, 0)
		vm.WritePop(e.outFile, mapping.SegmentPOINT, 0)
	}
	e.compileStatements()    // statements
	e.eat(vld.Identity("}")) // '}'
}

func (e *Engine) compileVarDec() {
	// Grammar: 'var' type varName (',' varName)* ';'
	e.eat(vld.Identity("var"))      // 'var'
	varType := e.compileType(false) // var type
	// varName (',' varName)*
	for {
		varName := e.eatType(tokenizer.TokenTypeIdentifier) // varName
		e.table.Define(varName, varType, mapping.IdentifierTypeVar)
		if e.currentToken().S != "," {
			break // no more var identifiers
		}
		e.eat(vld.Identity(",")) // ','
	}
	e.eat(vld.Identity(";")) // ';'
}

func (e *Engine) compileStatements() {
	// Grammar: letStatement | ifStatement | whileStatement | doStatement | return Statement
For:
	for {
		switch e.currentToken().S {
		case "let":
			e.compileLetStatement()
		case "if":
			e.compileIfStatement()
		case "while":
			e.compileWhileStatement()
		case "do":
			e.compileDoStatement()
		case "return":
			e.compileReturnStatement()
		default:
			break For // no more statements
		}
	}
}

func (e *Engine) compileDoStatement() {
	// Grammar: 'do' subroutineCall ';'
	vm.WriteComment(e.outFile, "Do statement")

	e.eat(vld.Identity("do")) // 'do'
	e.compileSubroutineCall() // subroutineCall
	e.eat(vld.Identity(";"))  // ';'
}

func (e *Engine) compileReturnStatement() {
	// Grammar: 'return' expression? ';'
	vm.WriteComment(e.outFile, "Return statement")

	e.eat(vld.Identity("return")) // 'return'
	if e.currentToken().S != ";" {
		e.compileExpression() // non void
	} else {
		vm.WritePush(e.outFile, mapping.SegmentCONST, 0) // void
	}
	vm.WriteReturn(e.outFile)
	e.eat(vld.Identity(";")) // ';'
}

func (e *Engine) compileWhileStatement() {
	// Grammar: 'while' '(' expression ')' '{' statements '}'
	vm.WriteComment(e.outFile, "While statement")
	e.labelCounter++

	label := e.className + ".while." + strconv.Itoa(e.labelCounter) + ".L1"
	vm.WriteLabel(e.outFile, label)

	e.eat(vld.Identity("while")) // 'while'
	e.eat(vld.Identity("("))     // '('
	e.compileExpression()        // expression
	e.eat(vld.Identity(")"))     // ')'

	vm.WriteArithmetic(e.outFile, mapping.ArithmCmdNOT)
	gotoLabel := e.className + ".while." + strconv.Itoa(e.labelCounter) + ".L2"
	vm.WriteIf(e.outFile, gotoLabel)

	e.eat(vld.Identity("{")) // '{'
	e.compileStatements()    // statements
	e.eat(vld.Identity("}")) // '}'

	vm.WriteGoto(e.outFile, label)
	vm.WriteLabel(e.outFile, gotoLabel)
}

func (e *Engine) compileIfStatement() {
	// Grammar: 'if' '(' expression ')' '{' statements '}'
	// ('else' '{' statements '}')?
	vm.WriteComment(e.outFile, "If statement")
	e.labelCounter++

	e.eat(vld.Identity("if")) // 'if'
	e.eat(vld.Identity("("))  // '('
	e.compileExpression()     // expression
	e.eat(vld.Identity(")"))  // ')'

	vm.WriteArithmetic(e.outFile, mapping.ArithmCmdNOT)
	label := e.className + ".if." + strconv.Itoa(e.labelCounter) + ".L1"
	vm.WriteIf(e.outFile, label)

	e.eat(vld.Identity("{")) // '{'
	gotoLabel := e.className + ".if." + strconv.Itoa(e.labelCounter) + ".L2"
	e.compileStatements() // statements
	vm.WriteGoto(e.outFile, gotoLabel)
	vm.WriteLabel(e.outFile, label)
	e.eat(vld.Identity("}")) // '}'

	// ('else' '{' statements '}')?
	if e.currentToken().S == "else" {
		vm.WriteComment(e.outFile, "Else statement")
		e.eat(vld.Identity("else")) // 'else'
		e.eat(vld.Identity("{"))    // '{'
		e.compileStatements()       // statements
		e.eat(vld.Identity("}"))    // '}'
	}
	vm.WriteLabel(e.outFile, gotoLabel)
}

func (e *Engine) compileLetStatement() {
	// Grammar: 'let' varName ('['expression']')? '=' expression ';'
	vm.WriteComment(e.outFile, "Let statement")

	e.eat(vld.Identity("let"))                             // 'let'
	identifier := e.eatType(tokenizer.TokenTypeIdentifier) // varName
	segment := e.table.KindOf(identifier)
	index := e.table.IndexOf(identifier)

	// ('['expression']')?
	isArray := e.currentToken().S == "["
	if isArray {
		e.eat(vld.Identity("[")) // '['
		e.compileExpression()    // 'expression'
		e.eat(vld.Identity("]")) // ']'
		vm.WritePush(e.outFile, segment, index)
		vm.WriteArithmetic(e.outFile, mapping.ArithmCmdADD)
	}

	e.eat(vld.Identity("=")) // '='
	e.compileExpression()    // expression
	e.eat(vld.Identity(";")) // ';'

	if isArray {
		vm.WritePop(e.outFile, mapping.SegmentTEMP, 0)
		vm.WritePop(e.outFile, mapping.SegmentPOINT, 1)
		vm.WritePush(e.outFile, mapping.SegmentTEMP, 0)
		vm.WritePop(e.outFile, mapping.SegmentTHAT, 0)
	} else {
		vm.WritePop(e.outFile, segment, index)
	}
}

func (e *Engine) compileSubroutineCall() {
	// Grammar: subroutineName '(' expressionList ')' |
	// (className|varName) '.' subroutineName '(' expressionList ')'

	// subroutineName | (className|varName)
	identifier := e.eatType(tokenizer.TokenTypeIdentifier)
	nArgs := 0
	funcName := ""
	className := identifier
	identifierType := ""

	if e.currentToken().S == "." {
		// '.' subroutineName
		e.eat(vld.Identity("."))                            // '.'
		funcName = e.eatType(tokenizer.TokenTypeIdentifier) // subroutineName
		identifierType = e.table.TypeOf(identifier)
	} else {
		className = e.className
		funcName = identifier
		nArgs++
		vm.WritePush(e.outFile, mapping.SegmentPOINT, 0)
	}

	if identifierType != "" {
		segment := e.table.KindOf(identifier)
		index := e.table.IndexOf(identifier)
		vm.WritePush(e.outFile, segment, index)
		nArgs++
		className = identifierType
	}

	e.eat(vld.Identity("("))           // '('
	nArgs += e.compileExpressionList() // expressionList
	e.eat(vld.Identity(")"))           // ')'

	vm.WriteCall(e.outFile, className, funcName, nArgs)
	vm.WritePop(e.outFile, mapping.SegmentTEMP, 0)
}

func (e *Engine) compileExpression() {
	// Grammar: term (op term)*
	ops := []string{"+", "-", "*", "/", "&", "|", "<", ">", "="}
	commands := make([]string, 0)
	for {
		e.compileTerm()
		if !vld.OneOf(ops...)(e.currentToken().S) {
			break
		}
		cmd, _ := mapping.ArithmSymbols[e.currentToken().S]
		commands = append(commands, cmd)
		e.eat(vld.OneOf(ops...))
	}
	for _, cmd := range commands {
		vm.WriteArithmetic(e.outFile, cmd)
	}
}

func (e *Engine) compileExpressionList() (nArgs int) {
	// Grammar: (expression (',' expression)*)?
	nArgs = 0
	for {
		// TODO verify this
		if e.currentToken().S == ")" {
			break
		}
		e.compileExpression() // expression
		if e.currentToken().S == "," {
			e.eat(vld.Identity(",")) // ','
		}
		nArgs++
	}
	return
}

func (e *Engine) compileTerm() {
	// Grammar: integetConstant | stringConstant | keywordConstant | varName |
	// varName '['expression']' | subroutineCall | '('expression')' | unaryOp term
	next := e.currentToken()
	afterNext := e.nextToken()

	if next.T == tokenizer.TokenTypeIntConstant {
		// integetConstant
		val, _ := strconv.Atoi(next.S)
		vm.WritePush(e.outFile, mapping.SegmentCONST, val)
		e.eat(vld.IsInt)
	} else if next.T == tokenizer.TokenTypeStringConstant {
		// stringConstant
		vm.WritePush(e.outFile, mapping.SegmentCONST, len(next.S))
		vm.WriteCall(e.outFile, "String", "new", 1)
		runes := []rune(next.S)
		for _, r := range runes {
			vm.WritePush(e.outFile, mapping.SegmentCONST, int(r))
			vm.WriteCall(e.outFile, "String", "appendChar", 2)
		}
		e.eat(vld.IsAny())
	} else if next.T == tokenizer.TokenTypeKeyword && next.S == "true" {
		// keywordConstant true
		vm.WritePush(e.outFile, mapping.SegmentCONST, 1)
		vm.WriteArithmetic(e.outFile, mapping.ArithmCmdNEG)
		e.eat(vld.Identity("true"))
	} else if next.T == tokenizer.TokenTypeKeyword && vld.OneOf("false", "null")(next.S) {
		// keywordConstant false | null
		vm.WritePush(e.outFile, mapping.SegmentCONST, 0)
		e.eat(vld.OneOf("false", "null"))
	} else if next.T == tokenizer.TokenTypeKeyword && next.S == "this" {
		// keywordConstant this
		vm.WritePush(e.outFile, mapping.SegmentPOINT, 0)
		e.eat(vld.Identity("this"))
	} else if next.T == tokenizer.TokenTypeSymbol && next.S == "(" {
		// '('expression')'
		e.eat(vld.Identity("("))
		e.compileExpression()
		e.eat(vld.Identity(")"))
	} else if next.T == tokenizer.TokenTypeSymbol && next.S == "-" {
		// unaryOp term -
		e.eat(vld.Identity("-"))
		e.compileTerm()
		vm.WriteArithmetic(e.outFile, mapping.ArithmCmdNEG)
	} else if next.T == tokenizer.TokenTypeSymbol && next.S == "~" {
		// unaryOp term ~
		e.eat(vld.Identity("~"))
		e.compileTerm()
		vm.WriteArithmetic(e.outFile, mapping.ArithmCmdNOT)
	} else if afterNext.T == tokenizer.TokenTypeSymbol && afterNext.S == "." {
		// class subroutine call
		// Grammar: (className|varName) '.' subroutineName '(' expressionList ')'
		identifier := e.eatType(tokenizer.TokenTypeIdentifier)
		e.eat(vld.Identity(".")) // '.'
		funcName := e.eatType(tokenizer.TokenTypeIdentifier)
		nArgs := 0
		className := identifier
		identifierType := e.table.TypeOf(identifier)
		if identifierType != "" {
			segment := e.table.KindOf(identifier)
			index := e.table.IndexOf(identifier)
			vm.WritePush(e.outFile, segment, index)
			nArgs++
			className = identifierType
		}
		e.eat(vld.Identity("("))           // '('
		nArgs += e.compileExpressionList() // expressionList
		e.eat(vld.Identity(")"))           // ')'
		vm.WriteCall(e.outFile, className, funcName, nArgs)
	} else {
		// other ?
		identifier := e.eatType(tokenizer.TokenTypeIdentifier) // varName
		index := e.table.IndexOf(identifier)
		segment := e.table.KindOf(identifier)
		vm.WritePush(e.outFile, segment, index)

		// '['expression']' ?
		next := e.currentToken()
		if next.T == tokenizer.TokenTypeSymbol && next.S == "[" {
			e.eat(vld.Identity("["))
			e.compileExpression()
			e.eat(vld.Identity("]"))
			vm.WriteArithmetic(e.outFile, mapping.ArithmCmdADD)
			vm.WritePop(e.outFile, mapping.SegmentPOINT, 1)
			vm.WritePush(e.outFile, mapping.SegmentTHAT, 0)
		}
	}
}

// validate next token value and return its value
func (e *Engine) eat(rules ...vld.Rule) string {
	nextToken, err := e.tokens.Lookup(0)
	if err != nil {
		panic("No more tokens")
	}
	for _, r := range rules {
		if r(nextToken.S) {
			e.tokens.Next()
			return nextToken.S
		}
	}
	panic("Rule not met for: " + nextToken.S)
}

// validate next token type and return its value
func (e *Engine) eatType(tType string) string {
	nextToken, err := e.tokens.Lookup(0)
	if err != nil {
		panic("No more tokens")
	}
	if nextToken.T == tType {
		e.tokens.Next()
		return nextToken.S
	}
	panic("Wrong token type: " + nextToken.S + " " + nextToken.T + ", expected: " + tType)
}

// get current token
func (e *Engine) currentToken() tokenizer.Token {
	token, err := e.tokens.Lookup(0)
	if err != nil {
		panic("No token at this position")
	}
	return token
}

// get next token
func (e *Engine) nextToken() tokenizer.Token {
	token, err := e.tokens.Lookup(1)
	if err != nil {
		panic("No token at this position")
	}
	return token
}
