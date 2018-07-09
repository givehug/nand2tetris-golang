package codewriter

import (
	"fmt"
	"nand2tetris-golang/common/writer"
	"strconv"
	"strings"
)

// CodeWriter type
type CodeWriter struct {
	*writer.Writer
	count    int
	fileName string
	funcName string
}

// New initializes CodeWriter
func New(filePath string) *CodeWriter {
	return &CodeWriter{writer.New(filePath), 0, "", ""}
}

// WriteArithmetic ...
func (cw *CodeWriter) WriteArithmetic(cmd string) {
	var asm string
	switch cmd {
	case "add":
		asm = "@SP\nAM=M-1\nD=M\nA=A-1\nM=D+M\n"
	case "sub":
		asm = "@SP\nAM=M-1\nD=M\nA=A-1\nM=M-D\n"
	case "neg":
		asm = "@SP\nA=M-1\nM=-M"
	case "and":
		asm = "@SP\nAM=M-1\nD=M\nA=A-1\nM=D&M\n"
	case "or":
		asm = "@SP\nAM=M-1\nD=M\nA=A-1\nM=D|M\n"
	case "not":
		asm = "@SP\nA=M-1\nM=!M\n"
	}
	cw.WriteLine(asm)
}

// WriteFunction ...
func (cw *CodeWriter) WriteFunction(funcName string, nArgs int) {
	cw.funcName = funcName
	s := "(" + funcName + ")\n@SP\nA=M\n"
	for i := 0; i < nArgs; i++ {
		s += "M=0\nA=A+1\n"
	}
	cw.WriteLine(s + "D=A\n@SP\nM=D\n")
}

// WriteCall ...
func (cw *CodeWriter) WriteCall(funcName string, nArgs int) {
	count := cw.nextCount()
	cw.WriteLine("@SP\nD=M\n@R13\nM=D\n" +
		"@ret." + count + "\nD=A\n@SP\nA=M\nM=D\n" +
		tplIncrementSP() +
		tplPointer("LCL") +
		tplIncrementSP() +
		tplPointer("ARG") +
		tplIncrementSP() +
		tplPointer("THIS") +
		tplIncrementSP() +
		tplPointer("THAT") +
		tplIncrementSP() +
		"@R13\nD=M\n@" + strconv.Itoa(nArgs) + "\nD=D-A\n@ARG\nM=D\n" +
		"@SP\nD=M\n@LCL\nM=D\n@" + funcName + "\n" +
		"0;JMP\n(ret." + count + ")\n")
}

// WriteReturn ...
func (cw *CodeWriter) WriteReturn() {
	cw.WriteLine("@LCL\nD=M\n@5\nA=D-A\nD=M\n@R13\nM=D\n" +
		"@SP\nA=M-1\nD=M\n@ARG\nA=M\nM=D\n" +
		"D=A+1\n@SP\nM=D\n" +
		"@LCL\nAM=M-1\nD=M\n@THAT\nM=D\n" +
		"@LCL\nAM=M-1\nD=M\n@THIS\nM=D\n" +
		"@LCL\nAM=M-1\nD=M\n@ARG\nM=D\n" +
		"@LCL\nA=M-1\nD=M\n@LCL\nM=D\n" +
		"@R13\nA=M\n0;JMP\n")
}

// WriteComparator ...
func (cw *CodeWriter) WriteComparator(arg1 string) {
	comp := strings.ToUpper(arg1)
	count := cw.nextCount()
	cw.WriteLine("@SP\nAM=M-1\nD=M\nA=A-1\nD=M-D\n" +
		"@" + comp + ".true." + count + "\nD;J" + comp + "\n" +
		"@SP\nA=M-1\nM=0\n@" + comp + ".after." + count + "\n" +
		"0;JMP\n(" + comp + ".true." + count + ")\n@SP\nA=M-1\n" +
		"M=-1\n(" + comp + ".after." + count + ")\n")
}

// WriteBranching ...
func (cw *CodeWriter) WriteBranching(brType, label string) {
	var asm string
	switch brType {
	case "label":
		asm = "(" + cw.funcName + "$" + label + ")\n"
	case "goto":
		asm = "@" + cw.funcName + "$" + label + "\n0;JMP\n"
	case "if-goto":
		asm = "@SP\nAM=M-1\nD=M\n@" + cw.funcName + "$" + label + "\nD;JNE\n"
	}
	cw.WriteLine(asm)
}

// WritePush ...
func (cw *CodeWriter) WritePush(segment string, val int) {
	var asm string
	valStr := strconv.Itoa(val)
	switch segment {
	case "constant":
		asm = "@" + valStr + "\nD=A\n"
	case "local":
		asm = "@LCL\nD=M\n@" + valStr + "\nA=D+A\nD=M\n"
	case "argument":
		asm = "@ARG\nD=M\n@" + valStr + "\nA=D+A\nD=M\n"
	case "this":
		asm = "@THIS\nD=M\n@" + valStr + "\nA=D+A\nD=M\n"
	case "that":
		asm = "@THAT\nD=M\n@" + valStr + "\nA=D+A\nD=M\n"
	case "pointer":
		if val == 0 {
			asm = "@THIS\nD=M\n"
		} else {
			asm = "@THAT\nD=M\n"
		}
	case "static":
		asm = "@" + cw.fileName + "." + valStr + "\nD=M\n"
	case "temp":
		asm = "@R5\nD=A\n@" + valStr + "\nA=D+A\nD=M\n"
	}
	cw.WriteLine(asm + "@SP\nA=M\nM=D\n@SP\nM=M+1\n")
}

// WritePop ...
func (cw *CodeWriter) WritePop(segment string, val int) { // arg2, arg3
	var asm string
	valStr := strconv.Itoa(val)
	switch segment {
	case "local":
		asm = "@LCL\nD=M\n@" + valStr + "\nD=D+A\n"
	case "argument":
		asm = "@ARG\nD=M\n@" + valStr + "\nD=D+A\n"
	case "this":
		asm = "@THIS\nD=M\n@" + valStr + "\nD=D+A\n"
	case "that":
		asm = "@THAT\nD=M\n@" + valStr + "\nD=D+A\n"
	case "pointer":
		if val == 0 {
			asm = "@THIS\nD=A\n"
		} else {
			asm = "@THAT\nD=A\n"
		}
	case "static":
		asm = "@" + cw.fileName + "." + valStr + "\nD=A\n"
	case "temp":
		asm = "@R5\nD=A\n@" + valStr + "\nD=D+A\n"
	}
	cw.WriteLine(asm + "@R13\nM=D\n@SP\nAM=M-1\nD=M\n@R13\nA=M\nM=D\n")
}

// WriteComment ...
func (cw *CodeWriter) WriteComment(s string) {
	cw.WriteLine("// " + s)
}

// WriteInit ...
func (cw *CodeWriter) WriteInit(callSysInit bool) {
	cw.WriteLine("@256\nD=A\n@SP\nM=D\n")
	// if sys.vm provided call Sys.init 0
	if callSysInit {
		cw.WriteComment("call Sys.init 0")
		cw.WriteCall("Sys.init", 0)
		cw.WriteLine("0;JMP\n")
	}
}

// SetFileName ...
func (cw *CodeWriter) SetFileName(fn string) {
	cw.fileName = fn
}

// Increment counter, return stringified count
func (cw *CodeWriter) nextCount() string {
	cw.count++
	return fmt.Sprintf("%d", cw.count)
}

func tplIncrementSP() string {
	return "@SP\nM=M+1\n"
}

func tplPointer(i string) string {
	return "@" + i + "\nD=M\n@SP\nA=M\nM=D\n"
}
