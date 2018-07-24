package st

import (
	"fmt"
)

// Identifier type constants
const (
	IdentifierTypeStatic = iota
	IdentifierTypeField
	IdentifierTypeArg
	IdentifierTypeVar
)

// TableEntry ...
type TableEntry struct {
	name    string
	varType string
	kind    int
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
func (table *SymbolTable) Define(name, varType string, kind int) {
	entry := &TableEntry{name, varType, kind, table.VarCount(kind)}
	if kind == IdentifierTypeVar || kind == IdentifierTypeArg {
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
func (table *SymbolTable) VarCount(kind int) int {
	if kind == IdentifierTypeVar || kind == IdentifierTypeArg {
		return varCount(kind, &table.subroutineTable)
	}
	return varCount(kind, &table.classTable)
}

// KindOf ...
func (table *SymbolTable) KindOf(name string) int {
	if val, ok := table.classTable[name]; ok {
		return val.kind
	}
	if val, ok := table.subroutineTable[name]; ok {
		return val.kind
	}
	return -1
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
		fmt.Println(1)
		return val.number
	}
	if val, ok := table.subroutineTable[name]; ok {
		fmt.Println(2, table.subroutineTable)
		return val.number
	}
	return -1
}

func varCount(kind int, t *Table) int {
	varCount := 0
	for _, val := range *t {
		if val.kind == kind {
			varCount++
		}
	}
	return varCount
}
