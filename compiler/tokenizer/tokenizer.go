package tokenizer

import (
	"nand2tetris-golang/common/parser"
	"nand2tetris-golang/common/utils"
	"strings"
)

// Token types
const (
	TokenTypeKeyword        = "keyword"
	TokenTypeSymbol         = "symbol"
	TokenTypeIntConstant    = "integerConstant"
	TokenTypeStringConstant = "stringConstant"
	TokenTypeIdentifier     = "identifier"
)

var keywords = map[string]bool{
	"class":       true,
	"constructor": true,
	"function":    true,
	"method":      true,
	"field":       true,
	"static":      true,
	"var":         true,
	"int":         true,
	"char":        true,
	"boolean":     true,
	"void":        true,
	"true":        true,
	"false":       true,
	"null":        true,
	"this":        true,
	"let":         true,
	"do":          true,
	"if":          true,
	"else":        true,
	"while":       true,
	"return":      true,
}

var symbols = map[string]bool{
	"{": true,
	"}": true,
	"(": true,
	")": true,
	"[": true,
	"]": true,
	".": true,
	",": true,
	";": true,
	"+": true,
	"-": true,
	"*": true,
	"/": true,
	"&": true,
	"|": true,
	"<": true,
	">": true,
	"=": true,
	"~": true,
}

// Token ...
type Token struct {
	S string
	T string
}

// GetTokens parses file and returns a list of Tokens
func GetTokens(sourceFile string) []Token {
	p := parser.New(sourceFile)
	tokenList := make([]Token, 0)
	t := "" // current token

	for {
		line, hasMore := p.Parse() // parser returns line by line
		if !hasMore {
			break
		}

		chars := []rune(line)
		charCount := len(chars)

		for i := 0; i <= charCount; i++ { // all chars in line + one more round
			var newChar rune
			if i < charCount {
				newChar = chars[i]
			}

			if _, in := keywords[t]; in {
				// is keyword
				tokenList = append(tokenList, Token{t, TokenTypeKeyword})
				t = ""
			} else if _, in := symbols[t]; in && (t != "/" || newChar != '*') {
				// is symbol (not / followed by * as block cocmment)
				tokenList = append(tokenList, Token{t, TokenTypeSymbol})
				t = ""
			} else if isString(t) {
				// is string, append excluding double quotes
				tokenList = append(tokenList, Token{t[1 : len(t)-1], TokenTypeStringConstant})
				t = ""
			} else if isIdentifier(t) && !isNonFirstCharOfIdentifier(newChar) {
				// is identifier
				tokenList = append(tokenList, Token{t, TokenTypeIdentifier})
				t = ""
			} else if isInt(t) && (newChar < '0' || newChar > '9') {
				// is int
				tokenList = append(tokenList, Token{t, TokenTypeIntConstant})
				t = ""
			} else if isBlockComment(t) {
				// is block comment, remove
				t = ""
			}

			if newChar == 0 {
				// new line chars always skip
			} else if newChar == ' ' && !strings.HasPrefix(t, "\"") {
				// is regular space, skip
			} else {
				// append new char
				t += string(newChar)
			}
		}
	}

	return tokenList
}

// ToXML generates tokens xml
func ToXML(tList []Token) string {
	eol := "\n"
	tokens := ""
	length := len(tList)

	for i, t := range tList {
		nl := eol
		if i == length-1 {
			nl = ""
		}
		tokens += utils.ToXML(t.T, t.S, true) + nl
	}

	return "<tokens>" + eol + tokens + eol + "</tokens>"
}

func isInt(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func isIdentifier(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i, r := range s {
		if i == 0 && (r >= '0' && r <= '9') {
			return false
		}
		if !isNonFirstCharOfIdentifier(r) {
			return false
		}
	}
	return true
}

func isNonFirstCharOfIdentifier(r rune) bool {
	if (r < 'a' || r > 'z') &&
		(r < 'A' || r > 'Z') &&
		(r < '0' || r > '9') &&
		r != '_' {
		return false
	}
	return true
}

func isString(s string) bool {
	return len(s) > 1 && strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")
}

func isBlockComment(s string) bool {
	return len(s) > 3 && strings.HasPrefix(s, "/*") && strings.HasSuffix(s, "*/")
}
