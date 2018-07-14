package validators

import "testing"

func TestValidators(t *testing.T) {
	t.Run("IsInt", func(t *testing.T) {
		// true
		t1 := IsInt("1")
		if t1 == false {
			t.Error("is int 1")
		}
		t2 := IsInt("155")
		if t2 == false {
			t.Error("is int 155")
		}
		// false
		f1 := IsInt("")
		if f1 == true {
			t.Error("is int empty string")
		}
		f2 := IsInt("1.2")
		if f2 == true {
			t.Error("is float")
		}
		f3 := IsInt(" 1 ")
		if f3 == true {
			t.Error("includes whitespace")
		}
		f4 := IsInt("1a")
		if f4 == true {
			t.Error("includes letters")
		}
	})
	t.Run("IsIdentifier", func(t *testing.T) {
		// true
		t1 := IsIdentifier("_")
		if t1 == false {
			t.Error("is _")
		}
		t2 := IsIdentifier("_1")
		if t2 == false {
			t.Error("is _1")
		}
		t3 := IsIdentifier("_a1")
		if t3 == false {
			t.Error("is _a1")
		}
		t4 := IsIdentifier("A1_a1")
		if t4 == false {
			t.Error("is A1_a1")
		}
		// false
		f1 := IsIdentifier("1")
		if f1 == true {
			t.Error("is 1")
		}
		f2 := IsIdentifier("_*")
		if f2 == true {
			t.Error("is _*")
		}
		f3 := IsIdentifier("")
		if f3 == true {
			t.Error("is whitespace")
		}
	})
	t.Run("IsNonFirstCharOfIdentifier", func(t *testing.T) {
		// true
		t1 := IsNonFirstCharOfIdentifier('_')
		if t1 == false {
			t.Error("is _")
		}
		t2 := IsNonFirstCharOfIdentifier('1')
		if t2 == false {
			t.Error("is 1")
		}
		t3 := IsNonFirstCharOfIdentifier('a')
		if t3 == false {
			t.Error("is _a1")
		}
		// false
		f1 := IsNonFirstCharOfIdentifier('*')
		if f1 == true {
			t.Error("is *")
		}
		f2 := IsNonFirstCharOfIdentifier('$')
		if f2 == true {
			t.Error("is $")
		}
	})
	t.Run("IsString", func(t *testing.T) {
		// true
		t1 := IsString("\"\"")
		if t1 == false {
			t.Error("empty string")
		}
		t2 := IsString("\" hello \"")
		if t2 == false {
			t.Error("non empty string")
		}
		// false
		f1 := IsString("")
		if f1 == true {
			t.Error("epmty")
		}
		f2 := IsString("\"")
		if f2 == true {
			t.Error("1 double quote")
		}
		f3 := IsString("a\"\"")
		if f3 == true {
			t.Error("wrong prefix")
		}
		f4 := IsString("\"\"a")
		if f4 == true {
			t.Error("wrong suffix")
		}
	})
	t.Run("IsBlockComment", func(t *testing.T) {
		// true
		t1 := IsBlockComment("/**/")
		if t1 == false {
			t.Error("t1")
		}
		t2 := IsBlockComment("/** hello */")
		if t2 == false {
			t.Error("t2")
		}
		t3 := IsBlockComment("/**\n * hello \n */")
		if t3 == false {
			t.Error("t3")
		}
		// false
		f1 := IsBlockComment("/**")
		if f1 == true {
			t.Error("f1")
		}
		f2 := IsBlockComment("**/")
		if f2 == true {
			t.Error("f2")
		}
		f3 := IsBlockComment("/*/")
		if f3 == true {
			t.Error("f3")
		}
		f4 := IsBlockComment("//")
		if f4 == true {
			t.Error("f4")
		}
	})
}
