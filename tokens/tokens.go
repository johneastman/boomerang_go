package tokens

import (
	"fmt"
)

type Token struct {
	Literal string
	Type    string
}

// Token names
const (
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
	PRINT                = "PRINT"
	FUNCTION             = "FUNCTION"
	NUMBER               = "NUMBER"
	IDENTIFIER           = "IDENTIFIER"
	EOF                  = "EOF"
)

// Tokens
var (
	PLUS_TOKEN                 = getToken(PLUS)
	MINUS_TOKEN                = getToken(MINUS)
	ASTERISK_TOKEN             = getToken(ASTERISK)
	FORWARD_SLASH_TOKEN        = getToken(FORWARD_SLASH)
	SEMICOLON_TOKEN            = getToken(SEMICOLON)
	OPEN_PAREN_TOKEN           = getToken(OPEN_PAREN)
	CLOSED_PAREN_TOKEN         = getToken(CLOSED_PAREN)
	ASSIGN_TOKEN               = getToken(ASSIGN)
	COMMA_TOKEN                = getToken(COMMA)
	OPEN_CURLY_BRACKET_TOKEN   = getToken(OPEN_CURLY_BRACKET)
	CLOSED_CURLY_BRACKET_TOKEN = getToken(CLOSED_CURLY_BRACKET)

	// Keywords
	PRINT_TOKEN    = getToken(PRINT)
	FUNCTION_TOKEN = getToken(FUNCTION)

	// Data Types
	NUMBER_TOKEN = getToken(NUMBER)

	// Misc
	IDENTIFIER_TOKEN = getToken(IDENTIFIER)
	EOF_TOKEN        = getToken(EOF) // end of file
)

var tokenData = map[string]Token{
	PLUS:                 {Type: "PLUS", Literal: "+"},
	MINUS:                {Type: "MINUS", Literal: "-"},
	ASTERISK:             {Type: "ASTERISK", Literal: "*"},
	FORWARD_SLASH:        {Type: "FORWARD_SLASH", Literal: "/"},
	SEMICOLON:            {Type: "SEMICOLON", Literal: ";"},
	OPEN_PAREN:           {Type: "OPEN_PAREN", Literal: "("},
	CLOSED_PAREN:         {Type: "CLOSED_PAREN", Literal: ")"},
	ASSIGN:               {Type: "ASSIGN", Literal: "="},
	COMMA:                {Type: "COMMA", Literal: ","},
	OPEN_CURLY_BRACKET:   {Type: "OPEN_CURLY_BRACKET", Literal: "{"},
	CLOSED_CURLY_BRACKET: {Type: "CLOSED_CURLY_BRACKET", Literal: "}"},
	EOF:                  {Type: "EOF", Literal: ""},
	PRINT:                {Type: "PRINT", Literal: "print"},
	FUNCTION:             {Type: "FUNCTION", Literal: "func"},
	NUMBER:               {Type: "NUMBER", Literal: ""},
	IDENTIFIER:           {Type: "IDENTIFIER", Literal: ""},
}

func getToken(name string) Token {
	if token, ok := tokenData[name]; ok {
		return token
	}
	panic(fmt.Sprintf("No token matching name: %s", name))
}

func GetTokenType(name string) string {
	token := getToken(name)
	return token.Type
}

var keywords = map[string]Token{
	"print": PRINT_TOKEN,
	"func":  FUNCTION_TOKEN,
}

func GetKeywordToken(literal string) Token {
	if token, ok := keywords[literal]; ok {
		return token
	}

	identifierToken := getToken(IDENTIFIER)
	identifierToken.Literal = literal
	return identifierToken
}

var symbols = map[byte]Token{
	'+': PLUS_TOKEN,
	'-': MINUS_TOKEN,
	'*': ASTERISK_TOKEN,
	'/': FORWARD_SLASH_TOKEN,
	';': SEMICOLON_TOKEN,
	'(': OPEN_PAREN_TOKEN,
	')': CLOSED_PAREN_TOKEN,
	'=': ASSIGN_TOKEN,
	',': COMMA_TOKEN,
	'{': OPEN_CURLY_BRACKET_TOKEN,
	'}': CLOSED_CURLY_BRACKET_TOKEN,
}

func getSymbolType(literal byte) Token {
	if tokenType, ok := symbols[literal]; ok {
		return tokenType
	}
	panic(fmt.Sprintf("Invalid symbol: %c", literal))
}
