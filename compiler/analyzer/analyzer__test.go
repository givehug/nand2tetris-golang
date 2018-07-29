package analyzer

import (
	"io/ioutil"
	"nand2tetris-golang/common/parsetree"
	"nand2tetris-golang/common/utils"
	"nand2tetris-golang/compiler/tokenizer"
	"testing"
)

func TestAnalyzerToXML(t *testing.T) {
	root := parsetree.New("root", "")
	child2 := parsetree.New("child2", "")
	child2.AddLeaves(
		parsetree.New("child21", "val21"),
		parsetree.New("child22", "val22"),
	)
	root.AddLeaves(
		parsetree.New("child1", "val1"),
		child2,
		parsetree.New("child3", "val3"),
	)

	if utils.CompareStrings(ToXML(root, 0), `
		<root>
			<child1> val1 </child1>
			<child2>
				<child21> val21 </child21>
				<child22> val22 </child22>
			</child2>
			<child3> val3 </child3>
		</root>
	`) {
		t.Error("bad xml")
	}
}

func TestAnalyzer(t *testing.T) {
	files := []string{
		"../test/analyzer/ArrayTest/Main",
		"../test/analyzer/ExpressionLessSquare/Main",
		"../test/analyzer/ExpressionLessSquare/Square",
		"../test/analyzer/ExpressionLessSquare/SquareGame",
		"../test/analyzer/Square/Main",
		"../test/analyzer/Square/Square",
		"../test/analyzer/Square/SquareGame",
	}

	for _, f := range files {
		tokens := tokenizer.GetTokens(f + ".jack")
		tree := CompileClass(&tokens)
		xml := ToXML(tree, 0)
		comp, _ := ioutil.ReadFile(f + ".xml")
		if !utils.CompareStrings(xml, string(comp)) {
			t.Error("bad xml " + f)
		}
	}
}
