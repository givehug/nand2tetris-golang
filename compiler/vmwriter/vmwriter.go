package vmwriter

import (
	"nand2tetris-golang/common/utils"
	"os"
)

func writeLine(f *os.File, s string) {
	fi, err := f.Stat()
	utils.HandleErr(err)
	eol := "\n"
	if fi.Size() == 0 {
		eol = ""
	}
	f.WriteString(eol + s)
}

// WritePush ...
// segment: CONST | ARG | LOCAL | STATIC | THI | THAT | POINTER | TEMP
func WritePush(f *os.File, segment string, index int) {
	writeLine(f, "todo")
}

// WritePop ...
// segment: ARG | LOCAL | STATIC | THI | THAT | POINTER | TEMP
func WritePop(f *os.File, segment string, index int) {
	writeLine(f, "todo")
}

// WriteArithmetic ...
// command: ADD, SUB, NEG, EQ, GT, LT, AND, OR, NOT
func WriteArithmetic(f *os.File, command string) {
	writeLine(f, "todo")
}

// WriteLabel ...
func WriteLabel(f *os.File, label string) {
	writeLine(f, "todo")
}

// WriteGoto ...
func WriteGoto(f *os.File, label string) {
	writeLine(f, "todo")
}

// WriteIf ...
func WriteIf(f *os.File, label string) {
	writeLine(f, "todo")
}

// WriteCall ...
func WriteCall(f *os.File, name string, nArgs int) {
	writeLine(f, "todo")
}

// WriteFunction ...
func WriteFunction(f *os.File, name string, nLocals int) {
	writeLine(f, "todo")
}

// WriteReturn ...
func WriteReturn(f *os.File) {
	writeLine(f, "todo")
}
