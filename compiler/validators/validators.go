package validators

import (
	"strings"
)

// Rule type
type Rule func(s string) bool

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

// OneOf returns Rule that compares string to one of the options
func OneOf(options ...string) Rule {
	return func(s string) bool {
		for _, o := range options {
			if o == s {
				return true
			}
		}
		return false
	}
}

// IsAny returns Rule which always returns true
func IsAny() Rule {
	return func(s string) bool {
		return true
	}
}
