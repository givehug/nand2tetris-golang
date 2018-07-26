package st

import (
	"nand2tetris-golang/compiler/mapping"
	"testing"
)

func TestCompilerSymbolTable(t *testing.T) {
	table := New()
	table.Define("x", "int", mapping.IdentifierTypeVar)
	if table.VarCount(mapping.IdentifierTypeVar) != 1 {
		t.Error("bad VarCount")
	}
	table.Define("name", "Point", mapping.IdentifierTypeStatic)
	if table.VarCount(mapping.IdentifierTypeStatic) != 1 {
		t.Error("bad VarCount classTable")
	}
	if table.IndexOf("x") != 0 {
		t.Error("bad IndexOf x")
	}
	if table.TypeOf("x") != "int" {
		t.Error("bad TypeOf")
	}
	if table.KindOf("x") != mapping.IdentifierTypeVar {
		t.Error("bad KindOf")
	}
	table.Define("y", "int", mapping.IdentifierTypeVar)
	if table.IndexOf("y") != 1 {
		t.Error("bad IndexOf y")
	}
	table.StartSubroutine()
	if ind := table.IndexOf("x"); ind != -1 {
		t.Error("bad IndexOf x after new subroutine started")
	}
}
