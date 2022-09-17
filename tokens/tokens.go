package tokens

import (
	"boomerang/node"
	"fmt"
)

// Token types
const (
	// Symbols
	PLUS          = "PLUS"
	MINUS         = "MINUS"
	ASTERISK      = "ASTERISK"
	FORWARD_SLASH = "FORWARD_SLASH"
	SEMICOLON     = "SEMICOLON"
	OPEN_PAREN    = "OPEN_PAREN"
	CLOSED_PAREN  = "CLOSED_PAREN"

	// Keywords
	LET = "LET"

	// Data Types
	NUMBER = node.NUMBER

	// Misc
	IDENTIFIER = "IDENTIFIER"
	EOF        = "EOF" // end of file
)

var keywords = map[string]string{
	"let": LET,
}

var symbols = map[byte]string{
	'+': PLUS,
	'-': MINUS,
	'*': ASTERISK,
	'/': FORWARD_SLASH,
	';': SEMICOLON,
	'(': OPEN_PAREN,
	')': CLOSED_PAREN,
}

func getSymbolType(literal byte) string {
	if tokenType, ok := symbols[literal]; ok {
		return tokenType
	}
	panic(fmt.Sprintf("Invalid symbol: %c", literal))
}

func getTokenType(literal string) string {
	if tokenType, ok := keywords[literal]; ok {
		return tokenType
	}
	return IDENTIFIER
}
