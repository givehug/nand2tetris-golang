package assembler

import (
	"nand2tetris-golang/assembler/codewriter"
	"nand2tetris-golang/assembler/parser"
	"nand2tetris-golang/assembler/symboltable"
	"nand2tetris-golang/common/utils"
	"os"
	"strconv"
)

func main() {
	// get desired source (.asm) and out (.hack) files from cli
	Compile(os.Args[1], os.Args[2])
}

// Compile compiles src (.asm) file into binary out (.hack) file
func Compile(src string, out string) {
	var lines []string        // lines slice
	st := symboltable.New()   // new symbol table
	p := parser.New(src)      // new parser
	cw := codewriter.New(out) // out file (.asm)
	defer p.Close()
	defer cw.Close()

	// First run:
	hasMore := true
	for hasMore {
		c, ok := p.Parse()
		hasMore = ok
		if ok {
			if parser.CommandType(c) != parser.CmdTypeL {
				// collect instructions
				lines = append(lines, c)
			} else {
				// remember labels
				label, _, _ := parser.CommandArgs(c)
				st[label] = len(lines)
			}
		}
	}

	// Second run:
	n := 16
	for _, l := range lines {
		cmdType := parser.CommandType(l)
		arg1, arg2, arg3 := parser.CommandArgs(l)

		switch cmdType {
		// process A instruction
		case parser.CmdTypeA:
			if parser.IsVariable(arg1) {
				// if variable, use from table
				_, found := st[arg1]
				if !found { // if not declared yet, put to table
					st[arg1] = n
					n++
				}
				cw.WriteA(st[arg1])
			} else {
				// if value, use directly
				aInt, err := strconv.Atoi(arg1)
				utils.HandleErr(err)
				cw.WriteA(aInt)
			}
		// process C instruction
		case parser.CmdTypeC:
			cw.WriteC(arg1, arg2, arg3)
		}
	}

	utils.LogDone(src, out)
}
