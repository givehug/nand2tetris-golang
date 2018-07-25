package writer

import (
	"nand2tetris-golang/common/utils"
	"os"
)

// TODO: This module is obsolete,
// use os.File, Close, Create and WriteLine directly.
// Each compiled file should end with empty line for good.

// WriterInterface ...
type WriterInterface interface {
	Close()
	WriteLine(s string)
}

// Writer type
type Writer struct {
	file *os.File
	WriterInterface
}

// New creates new Writer
func New(outFile string) *Writer {
	f, err := os.Create(outFile)
	utils.HandleErr(err)
	return &Writer{
		file: f,
	}
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
