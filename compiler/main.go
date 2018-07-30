package main

import (
	"fmt"
	"nand2tetris-golang/compiler/engine"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	jackFileExt = ".jack"
	vmFileExt   = ".vm"
)

func main() {
	// handle directory or jack file path as single cli arg
	path := os.Args[1]

	// prep list with all .jack files to compile
	var files []string
	if strings.HasSuffix(path, jackFileExt) {
		files = []string{path}
	} else {
		files, _ = filepath.Glob(path + "/*" + jackFileExt)
	}

	// process each file as separate goroutine
	wg := sync.WaitGroup{}
	defer wg.Wait()
	for _, f := range files {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			engine.Compile(f, strings.Replace(f, jackFileExt, vmFileExt, 1))
		}(f)
	}

	fmt.Println("done")
}
