package tokens

import (
	"fmt"
)

type Token struct {
	Literal string
	Type    string
}

// Token types
const (
	// Symbols
	PLUS                 = "PLUS"
	MINUS                = "MINUS"
	ASTERISK             = "ASTERISK"
	FORWARD_SLASH        = "FORWARD_SLASH"
	SEMICOLON            = "SEMICOLON"
	OPEN_PAREN           = "OPEN_PAREN"
	CLOSED_PAREN         = "CLOSED_PAREN"
	ASSIGN               = "ASSIGN"
	COMMA                = "COMMA"
	OPEN_CURLY_BRACKET   = "OPEN_CURLY_BRACKET"
	CLOSED_CURLY_BRACKET = "CLOSED_CURLY_BRACKET"

	// Keywords
	PRINT    = "PRINT"
	FUNCTION = "FUNCTION"

	// Data Types
	NUMBER = "NUMBER"

	// Misc
	IDENTIFIER = "IDENTIFIER"
	EOF        = "EOF" // end of file
)

var keywords = map[string]string{
	"print": PRINT,
	"func":  FUNCTION,
}

var symbols = map[byte]string{
	'+': PLUS,
	'-': MINUS,
	'*': ASTERISK,
	'/': FORWARD_SLASH,
	';': SEMICOLON,
	'(': OPEN_PAREN,
	')': CLOSED_PAREN,
	'=': ASSIGN,
	',': COMMA,
	'{': OPEN_CURLY_BRACKET,
	'}': CLOSED_CURLY_BRACKET,
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
