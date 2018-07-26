package engine

import (
	"fmt"
	pt "nand2tetris-golang/common/parsetree"
	"nand2tetris-golang/compiler/analyzer"
	st "nand2tetris-golang/compiler/symboltable"
	"nand2tetris-golang/compiler/tokenizer"
)

// CompileClass ...
func CompileClass(inFile, outFile string) {
	// fmt.Println(inFile, outFile)
	tokens := tokenizer.GetTokens(inFile)  // get tokens list
	tree := analyzer.CompileClass(&tokens) // create parse tree
	// fmt.Println(tokenizer.ToXML(tokens)) // print tokens XML
	// fmt.Println(analyzer.ToXML(tree, 0)) // print parsed tree XML
	table := st.New() // create symbol table

	for _, leaf := range tree.Leaves() {
		if leaf.Type() == analyzer.RuleTypeClassVarDec {
			compileClassVarDec(leaf, table)
		}
		if leaf.Type() == analyzer.RuleTypeSubroutineDec {
			compileClassSubroutineDec(leaf, table)
		}
	}
}

func compileClassVarDec(tree *pt.ParseTree, table *st.SymbolTable) {
	// add vars to table
	leaves := tree.Leaves()
	kind := leaves[0].Value() // 'static' | 'field'
	varType := leaves[1].Value()
	for _, leaf := range leaves {
		if leaf.Type() == analyzer.RuleTypeIdentifier {
			table.Define(leaf.Value(), varType, kind)
		}
	}
}

func compileClassSubroutineDec(tree *pt.ParseTree, table *st.SymbolTable) {
	// todo continue processing
	fmt.Println("compileClassSubroutineDec")
}

func compileParameterList() {}

func compileSubroutineBody() {}

func compileVarDec() {}

func compileStatements() {}

func compileDo() {}

func compileReturn() {}

func compileWhile() {}

func compileIf() {}

func compileLet() {}

// func compileSubroutineCall() {}

func compileExpression() {}

func compileExpressionList() {}

func compileTerm() {}

// func compileType() {}
