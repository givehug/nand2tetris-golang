package tokenizer

import (
	"errors"
	"nand2tetris-golang/common/parser"
	"nand2tetris-golang/common/utils"
	"nand2tetris-golang/compiler/validators"
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
	S string // token string value
	T string // token type constant
}

// TokensInterface ...
type TokensInterface interface {
	Next() (Token, error)
	Lookup() (Token, error)
}

// Tokens ...
type Tokens struct {
	list  []Token
	index int
}

// GetTokens parses file and returns a list of Tokens
func GetTokens(sourceFile string) *Tokens {
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

			if _, in := keywords[t]; in && !validators.IsIdentifier(string(newChar)) {
				// is keyword
				tokenList = append(tokenList, Token{t, TokenTypeKeyword})
				t = ""
			} else if _, in := symbols[t]; in && (t != "/" || newChar != '*') {
				// is symbol (not / followed by * as block cocmment)
				tokenList = append(tokenList, Token{t, TokenTypeSymbol})
				t = ""
			} else if validators.IsString(t) {
				// is string, append excluding double quotes
				tokenList = append(tokenList, Token{t[1 : len(t)-1], TokenTypeStringConstant})
				t = ""
			} else if validators.IsIdentifier(t) && !validators.IsNonFirstCharOfIdentifier(newChar) {
				// is identifier
				tokenList = append(tokenList, Token{t, TokenTypeIdentifier})
				t = ""
			} else if validators.IsInt(t) && (newChar < '0' || newChar > '9') {
				// is int
				tokenList = append(tokenList, Token{t, TokenTypeIntConstant})
				t = ""
			} else if validators.IsBlockComment(t) {
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

	return &Tokens{tokenList, 0}
}

// ToXML generates tokens xml
func ToXML(tList *Tokens) string {
	eol := "\n"
	tokens := ""

	for {
		t, err := tList.Lookup(0)
		if err != nil {
			break
		}
		tokens += utils.ToXMLTag(t.T, " "+t.S+" ") + eol
		tList.Next()
	}

	return "<tokens>" + eol + tokens + "</tokens>"
}

// Next increments token index, returns next token
func (t *Tokens) Next() (Token, error) {
	t.index++
	if t.index >= len(t.list) {
		return t.list[0], errors.New("no more tokens")
	}
	return t.list[t.index], nil
}

// Lookup returns current token index + i Token
func (t *Tokens) Lookup(i int) (Token, error) {
	if t.index+i >= len(t.list) {
		return t.list[0], errors.New("token list length exceeded")
	}
	return t.list[t.index+i], nil
}
