package main

import (
	"fmt"
	"nand2tetris-golang/compiler/analyzer"
	"nand2tetris-golang/compiler/tokenizer"
	"os"
)

func main() {
	path := os.Args[1]

	// tokens
	tokens := tokenizer.GetTokens(path)
	// fmt.Println(tokenizer.ToXML(tokens))
	// fmt.Println()

	// parse tree
	a := analyzer.New(&tokens)
	tree := a.CompileClass()
	fmt.Println(analyzer.ToXML(tree, 0))
}
