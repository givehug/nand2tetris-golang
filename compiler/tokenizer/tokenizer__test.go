package tokenizer

import (
	"io/ioutil"
	"nand2tetris-golang/common/utils"
	"testing"
)

func TestTokenizerToXML(t *testing.T) {
	files := []string{
		"../test/ArrayTest/Main",
		"../test/ExpressionLessSquare/Main",
		"../test/ExpressionLessSquare/Square",
		"../test/ExpressionLessSquare/SquareGame",
		"../test/Square/Main",
		"../test/Square/Square",
		"../test/Square/SquareGame",
	}

	for _, f := range files {
		xml := ToXML(GetTokens(f + ".jack"))
		comp, _ := ioutil.ReadFile(f + "T.xml")
		if !utils.CompareStrings(xml, string(comp)) {
			t.Error("bad xml " + f)
		}
	}
}
