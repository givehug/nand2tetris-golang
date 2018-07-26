package st

import (
	"nand2tetris-golang/compiler/mapping"
)

// TableEntry ...
type TableEntry struct {
	name    string
	varType string
	kind    string
	number  int
}

// Table ...
type Table map[string]*TableEntry

// SymbolTableInterface ...
type SymbolTableInterface interface {
	StartSubroutine()
	Define(name, varType string, kind string)
	VarCount(kind int) int
	KindOf(name string) int
	TypeOf(name string) string
	IndexOf(name string) int
}

// SymbolTable ...
type SymbolTable struct {
	SymbolTableInterface
	classTable      Table
	subroutineTable Table
}

// New ...
func New() *SymbolTable {
	return &SymbolTable{
		classTable:      make(Table),
		subroutineTable: make(Table),
	}
}

// Define ...
func (table *SymbolTable) Define(name, varType, kind string) {
	entry := &TableEntry{name, varType, kind, table.VarCount(kind)}
	if kind == mapping.IdentifierTypeVar || kind == mapping.IdentifierTypeArg {
		table.subroutineTable[name] = entry
	} else {
		table.classTable[name] = entry
	}
}

// StartSubroutine ...
func (table *SymbolTable) StartSubroutine() {
	table.subroutineTable = make(Table)
}

// VarCount ...
func (table *SymbolTable) VarCount(kind string) int {
	if kind == mapping.IdentifierTypeVar || kind == mapping.IdentifierTypeArg {
		return varCount(kind, &table.subroutineTable)
	}
	return varCount(kind, &table.classTable)
}

// KindOf ...
func (table *SymbolTable) KindOf(name string) string {
	if val, ok := table.classTable[name]; ok {
		return val.kind
	}
	if val, ok := table.subroutineTable[name]; ok {
		return val.kind
	}
	return ""
}

// TypeOf ...
func (table *SymbolTable) TypeOf(name string) string {
	if val, ok := table.classTable[name]; ok {
		return val.varType
	}
	if val, ok := table.subroutineTable[name]; ok {
		return val.varType
	}
	return ""
}

// IndexOf ...
func (table *SymbolTable) IndexOf(name string) int {
	if val, ok := table.classTable[name]; ok {
		return val.number
	}
	if val, ok := table.subroutineTable[name]; ok {
		return val.number
	}
	return -1
}

func varCount(kind string, t *Table) int {
	varCount := 0
	for _, val := range *t {
		if val.kind == kind {
			varCount++
		}
	}
	return varCount
}
