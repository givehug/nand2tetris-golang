package tokenizer

import (
	"io/ioutil"
	"nand2tetris-golang/common/utils"
	"testing"
)

func TestTokenizerToXML(t *testing.T) {
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
		xml := ToXML(GetTokens(f + ".jack"))
		comp, _ := ioutil.ReadFile(f + "T.xml")
		if !utils.CompareStrings(xml, string(comp)) {
			t.Error("bad xml " + f)
		}
	}
}
