package main

import (
	"nand2tetris-golang/common/utils"
	"nand2tetris-golang/vm_translator/codewriter"
	"nand2tetris-golang/vm_translator/parser"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// get source and out files from cli
	path := os.Args[1]

	// normalize dir, file name, depending which option was provided
	name, dir, isFile := utils.PathInfo(path)
	if !isFile {
		path = dir + "*"
	}

	// get all files in path
	files, err := filepath.Glob(path)
	utils.HandleErr(err)

	// init code writer
	cw := codewriter.New(dir + name + ".asm")
	cw.WriteInit(isSysFileProvided(files))
	defer cw.Close()

	// parse and process all files
	for _, f := range files {
		currFileName, _, _ := utils.PathInfo(f)
		cw.SetFileName(currFileName)
		p := parser.New(f)
		defer p.Close()

		// parse file lines
		hasMore := true
		for hasMore {
			c, ok := p.Parse()
			hasMore = ok
			if ok {
				cw.WriteComment(c)
				arg1, arg2, arg3 := parser.CommandArgs(c)
				switch parser.CommandType(c) {
				case parser.CmdTypeArithmetic:
					cw.WriteArithmetic(arg1)
				case parser.CmdTypePush:
					cw.WritePush(arg2, arg3)
				case parser.CmdTypePop:
					cw.WritePop(arg2, arg3)
				case parser.CmdTypeComparator:
					cw.WriteComparator(arg1)
				case parser.CmdTypeBranching:
					cw.WriteBranching(arg1, arg2)
				case parser.CmdTypeFunction:
					cw.WriteFunction(arg2, arg3)
				case parser.CmdTypeCall:
					cw.WriteCall(arg2, arg3)
				case parser.CmdTypeReturn:
					cw.WriteReturn()
				}
			}
		}
	}

	utils.LogDone(path, dir+name+".asm")
}

func isSysFileProvided(files []string) bool {
	for _, f := range files {
		if strings.HasSuffix(f, "Sys.vm") {
			return true
		}
	}
	return false
}
