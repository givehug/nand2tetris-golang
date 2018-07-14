package main

import (
	"fmt"
	"nand2tetris-golang/compiler/tokenizer"
	"os"
)

func main() {
	path := os.Args[1]
	fmt.Println(tokenizer.GetTokens(path))
}
