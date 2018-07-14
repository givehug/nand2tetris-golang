package analyzer

import (
	"nand2tetris-golang/common/parsetree"
	"nand2tetris-golang/common/utils"
	"testing"
)

func TestAnalyzerToXML(t *testing.T) {
	root := parsetree.New("root", "")
	child2 := parsetree.New("child2", "")
	child2.AddChildren(
		parsetree.New("child21", "val21"),
		parsetree.New("child22", "val22"),
	)
	root.AddChildren(
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
