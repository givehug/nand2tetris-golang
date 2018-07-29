package vmwriter

import (
	"os"
	"strconv"
)

func writeLine(f *os.File, s string) {
	f.WriteString(s + "\n")
}

// WritePush ...
// segment: one of mapping push pop segments
func WritePush(f *os.File, segment string, index int) {
	writeLine(f, "push "+segment+" "+strconv.Itoa(index))
}

// WritePop ...
// segment: one of mapping push pop segments
func WritePop(f *os.File, segment string, index int) {
	writeLine(f, "pop "+segment+" "+strconv.Itoa(index))
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
	writeLine(f, "call "+className+"."+subName+" "+strconv.Itoa(nArgs))
}

// WriteFunction ...
func WriteFunction(f *os.File, className, subName string, nLocals int) {
	writeLine(f, "function "+className+"."+subName+" "+strconv.Itoa(nLocals))
}

// WriteReturn ...
func WriteReturn(f *os.File) {
	writeLine(f, "return")
}

// WriteComment ...
func WriteComment(f *os.File, s string) {
	writeLine(f, "")
	writeLine(f, "// "+s)
}
