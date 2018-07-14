package validators

import (
	"strings"
)

// Rule type
type Rule func(s string) bool

// Validate returns string if valid, panics if not
func Validate(s string, r Rule) string {
	if r(s) == false {
		panic("rule not met: " + s + " | ") // todo return error
		// utils.GetFunctionName(r) // todo cannot compile this
	}
	return s
}

// Rules:

// IsInt ...
func IsInt(s string) bool {
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

// IsIdentifier ...
func IsIdentifier(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i, r := range s {
		if i == 0 && (r >= '0' && r <= '9') {
			return false
		}
		if !IsNonFirstCharOfIdentifier(r) {
			return false
		}
	}
	return true
}

// IsNonFirstCharOfIdentifier ...
func IsNonFirstCharOfIdentifier(r rune) bool {
	if (r < 'a' || r > 'z') &&
		(r < 'A' || r > 'Z') &&
		(r < '0' || r > '9') &&
		r != '_' {
		return false
	}
	return true
}

// IsString ...
func IsString(s string) bool {
	return len(s) > 1 && strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")
}

// IsBlockComment ...
func IsBlockComment(s string) bool {
	return len(s) > 3 && strings.HasPrefix(s, "/*") && strings.HasSuffix(s, "*/")
}

// Identity returns Rule that tests identity
func Identity(a string) Rule {
	return func(b string) bool {
		return a == b
	}
}
