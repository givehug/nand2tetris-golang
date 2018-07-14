package main

import (
	"fmt"
	"nand2tetris-golang/compiler/tokenizer"
	"os"
)

func main() {
	path := os.Args[1]
	xml := tokenizer.ToXML(tokenizer.GetTokens(path))

	fmt.Println(xml)
}
