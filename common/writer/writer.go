package writer

import (
	"nand2tetris-golang/common/utils"
	"os"
)

// Writer type
type Writer struct {
	file *os.File
}

// New creates new Writer
func New(outFile string) *Writer {
	f, err := os.Create(outFile)
	utils.HandleErr(err)
	return &Writer{f}
}

// Close closes file
func (w *Writer) Close() {
	w.file.Close()
}

// WriteLine writes line to file
func (w *Writer) WriteLine(s string) {
	fi, err := w.file.Stat()
	utils.HandleErr(err)
	eol := "\n"
	if fi.Size() == 0 {
		eol = ""
	}
	w.file.WriteString(eol + s)
}
