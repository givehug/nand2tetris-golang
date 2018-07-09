package codewriter

import (
	"fmt"
	"nand2tetris-golang/common/writer"
)

// CodeWriter type
type CodeWriter struct {
	*writer.Writer
}

// New creates new CodeWriter
func New(outFile string) *CodeWriter {
	return &CodeWriter{writer.New(outFile)}
}

// WriteA translates a-instruction value to binary code
func (cw *CodeWriter) WriteA(v int) {
	s := fmt.Sprintf("%015b", v) // int to binary string, padd with 0
	hack := "0" + s[len(s)-15:]  // trim and prefix with 0
	cw.WriteLine(hack)
}

// WriteC translates c-instruction (d,c,j) to binary code
func (cw *CodeWriter) WriteC(d, c, j string) {
	if len(d) == 0 {
		d = "null"
	}
	if len(j) == 0 {
		j = "null"
	}
	hack := fmt.Sprintf("111%s%s%s", computations[c], destinations[d], jumps[j])
	cw.WriteLine(hack)
}
