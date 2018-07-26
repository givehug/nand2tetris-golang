package main

import (
	"nand2tetris-golang/compiler/engine"
	"os"
	"path/filepath"
	"strings"
)

const (
	jackFileExt = ".jack"
	vmFileExt   = ".vm"
)

func main() {
	// handle directory or jack file path as single cli arg
	path := os.Args[1]
	// compile .vm file for each .jack file in dir
	var files []string
	if strings.HasSuffix(path, jackFileExt) {
		files = []string{path}
	} else {
		files, _ = filepath.Glob(path + "/*" + jackFileExt)
	}
	for _, f := range files {
		engine.CompileClass(f, strings.Replace(f, jackFileExt, vmFileExt, 1))
	}
}
