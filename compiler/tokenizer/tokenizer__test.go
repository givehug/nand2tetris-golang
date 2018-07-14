package tokenizer

import (
	"io/ioutil"
	"nand2tetris-golang/common/utils"
	"testing"
)

func TestTokenizerHelpers(t *testing.T) {
	t.Run("isInt", func(t *testing.T) {
		// true
		t1 := isInt("1")
		if t1 == false {
			t.Error("is int 1")
		}
		t2 := isInt("155")
		if t2 == false {
			t.Error("is int 155")
		}
		// false
		f1 := isInt("")
		if f1 == true {
			t.Error("is int empty string")
		}
		f2 := isInt("1.2")
		if f2 == true {
			t.Error("is float")
		}
		f3 := isInt(" 1 ")
		if f3 == true {
			t.Error("includes whitespace")
		}
		f4 := isInt("1a")
		if f4 == true {
			t.Error("includes letters")
		}
	})
	t.Run("isIdentifier", func(t *testing.T) {
		// true
		t1 := isIdentifier("_")
		if t1 == false {
			t.Error("is _")
		}
		t2 := isIdentifier("_1")
		if t2 == false {
			t.Error("is _1")
		}
		t3 := isIdentifier("_a1")
		if t3 == false {
			t.Error("is _a1")
		}
		t4 := isIdentifier("A1_a1")
		if t4 == false {
			t.Error("is A1_a1")
		}
		// false
		f1 := isIdentifier("1")
		if f1 == true {
			t.Error("is 1")
		}
		f2 := isIdentifier("_*")
		if f2 == true {
			t.Error("is _*")
		}
		f3 := isIdentifier("")
		if f3 == true {
			t.Error("is whitespace")
		}
	})
	t.Run("isNonFirstCharOfIdentifier", func(t *testing.T) {
		// true
		t1 := isNonFirstCharOfIdentifier('_')
		if t1 == false {
			t.Error("is _")
		}
		t2 := isNonFirstCharOfIdentifier('1')
		if t2 == false {
			t.Error("is 1")
		}
		t3 := isNonFirstCharOfIdentifier('a')
		if t3 == false {
			t.Error("is _a1")
		}
		// false
		f1 := isNonFirstCharOfIdentifier('*')
		if f1 == true {
			t.Error("is *")
		}
		f2 := isNonFirstCharOfIdentifier('$')
		if f2 == true {
			t.Error("is $")
		}
	})
	t.Run("isString", func(t *testing.T) {
		// true
		t1 := isString("\"\"")
		if t1 == false {
			t.Error("empty string")
		}
		t2 := isString("\" hello \"")
		if t2 == false {
			t.Error("non empty string")
		}
		// false
		f1 := isString("")
		if f1 == true {
			t.Error("epmty")
		}
		f2 := isString("\"")
		if f2 == true {
			t.Error("1 double quote")
		}
		f3 := isString("a\"\"")
		if f3 == true {
			t.Error("wrong prefix")
		}
		f4 := isString("\"\"a")
		if f4 == true {
			t.Error("wrong suffix")
		}
	})
	t.Run("isBlockComment", func(t *testing.T) {
		// true
		t1 := isBlockComment("/**/")
		if t1 == false {
			t.Error("t1")
		}
		t2 := isBlockComment("/** hello */")
		if t2 == false {
			t.Error("t2")
		}
		t3 := isBlockComment("/**\n * hello \n */")
		if t3 == false {
			t.Error("t3")
		}
		// false
		f1 := isBlockComment("/**")
		if f1 == true {
			t.Error("f1")
		}
		f2 := isBlockComment("**/")
		if f2 == true {
			t.Error("f2")
		}
		f3 := isBlockComment("/*/")
		if f3 == true {
			t.Error("f3")
		}
		f4 := isBlockComment("//")
		if f4 == true {
			t.Error("f4")
		}
	})
}

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
		// fmt.Println(xml)
		if !utils.CompareStrings(xml, string(comp)) {
			t.Error("bad xml " + f)
		}
	}
}
