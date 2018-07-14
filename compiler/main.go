package main

import (
	"fmt"
	"nand2tetris-golang/compiler/analizer"
	"nand2tetris-golang/compiler/tokenizer"
	"os"
)

func main() {
	path := os.Args[1]

	// tokens
	tokens := tokenizer.GetTokens(path)
	fmt.Println(tokenizer.ToXML(tokens))

	fmt.Println()

	// parse tree
	a := analizer.New(&tokens)
	tree := a.CompileClass()
	fmt.Println(analizer.ToXML(tree, 0))
}
