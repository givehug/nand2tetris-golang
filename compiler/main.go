package main

import (
	"fmt"
	"nand2tetris-golang/compiler/analyzer"
	"nand2tetris-golang/compiler/tokenizer"
	"os"
)

func main() {
	file := os.Args[1]

	tokens := tokenizer.GetTokens(file)    // tokens list
	tree := analyzer.CompileClass(&tokens) // parse tree

	// Print tokens XML
	// fmt.Println(tokenizer.ToXML(tokens))

	// Print parsed tree XML
	fmt.Println(analyzer.ToXML(tree, 0))
}
