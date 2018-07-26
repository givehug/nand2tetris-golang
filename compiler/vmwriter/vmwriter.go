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
// segment: one of mapping push pop segments
func WritePush(f *os.File, segment string, index int) {
	writeLine(f, "push "+segment+" "+string(index))
}

// WritePop ...
// segment: one of mapping push pop segments
func WritePop(f *os.File, segment string, index int) {
	writeLine(f, "push "+segment+" "+string(index))
}

// WriteArithmetic ...
// command: one of mapping arithmetic commands
func WriteArithmetic(f *os.File, command string) {
	writeLine(f, command)
}

// WriteLabel ...
func WriteLabel(f *os.File, label string) {
	writeLine(f, "label "+label)
}

// WriteGoto ...
func WriteGoto(f *os.File, label string) {
	writeLine(f, "goto "+label)
}

// WriteIf ...
func WriteIf(f *os.File, label string) {
	writeLine(f, "if-goto "+label)
}

// WriteCall ...
func WriteCall(f *os.File, className, subName string, nArgs int) {
	writeLine(f, "call "+className+"."+subName+" "+string(nArgs))
}

// WriteFunction ...
func WriteFunction(f *os.File, className, subName string, nLocals int) {
	writeLine(f, "function "+className+"."+subName+" "+string(nLocals))
}

// WriteReturn ...
func WriteReturn(f *os.File) {
	writeLine(f, "return")
}
