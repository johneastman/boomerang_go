package tokens

import (
	"fmt"
)

type TokenMetaData struct {
	Literal string
	Type    string
	IsRegex bool
}

type Token struct {
	Literal    string
	Type       string
	LineNumber int
}

// Token types/labels
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
	STRING               = "STRING"
	BOOLEAN              = "BOOLEAN"
	IDENTIFIER           = "IDENTIFIER"
	EOF                  = "EOF"
	PTR                  = "POINTER"
	TRUE                 = "TRUE"
	FALSE                = "FALSE"
	OPEN_BRACKET         = "OPEN_BRACKET"
	CLOSED_BRACKET       = "CLOSED_BRACKET"
	ELSE                 = "ELSE"
	AT                   = "AT"
	INLINE_COMMENT       = "INLINE_COMMENT"
	BLOCK_COMMENT        = "BLOCK_COMMENT"
	NOT                  = "NOT"
	EQ                   = "EQUAL"
	WHEN                 = "WHEN"
	IS                   = "IS"
)

// Tokens
var (
	// Symbols
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
	PTR_TOKEN                  = getToken(PTR)
	OPEN_BRACKET_TOKEN         = getToken(OPEN_BRACKET)
	CLOSED_BRACKET_TOKEN       = getToken(CLOSED_BRACKET)
	AT_TOKEN                   = getToken(AT)
	INLINE_COMMENT_TOKEN       = getToken(INLINE_COMMENT)
	BLOCK_COMMENT_TOKEN        = getToken(BLOCK_COMMENT)
	NOT_TOKEN                  = getToken(NOT)
	EQ_TOKEN                   = getToken(EQ)

	// Keywords
	PRINT_TOKEN    = getToken(PRINT)
	FUNCTION_TOKEN = getToken(FUNCTION)
	TRUE_TOKEN     = Token{Type: BOOLEAN, Literal: "true"}
	FALSE_TOKEN    = Token{Type: BOOLEAN, Literal: "false"}
	ELSE_TOKEN     = getToken(ELSE)
	WHEN_TOKEN     = getToken(WHEN)
	IS_TOKEN       = getToken(IS)

	// Data Types
	NUMBER_TOKEN  = getToken(NUMBER)
	STRING_TOKEN  = getToken(STRING)
	BOOLEAN_TOKEN = getToken(BOOLEAN)

	// Misc
	IDENTIFIER_TOKEN = getToken(IDENTIFIER)
	EOF_TOKEN        = getToken(EOF) // end of file
)

var tokenData = []TokenMetaData{
	// Data types/misc
	{Type: NUMBER, Literal: "[0-9]*[.]?[0-9]+", IsRegex: true},
	{Type: STRING, Literal: "\"(.*)\"", IsRegex: true},
	{Type: BOOLEAN, Literal: "(true|false)", IsRegex: true},

	// Keywords. Need to be defined before "IDENTIFIER" in this list so they are not misclassified
	{Type: WHEN, Literal: "when", IsRegex: true},
	{Type: IS, Literal: "is", IsRegex: true},
	{Type: NOT, Literal: "not", IsRegex: true},
	{Type: ELSE, Literal: "else", IsRegex: true},
	{Type: PRINT, Literal: "print", IsRegex: true},
	{Type: FUNCTION, Literal: "func", IsRegex: true},

	// Identifier
	{Type: IDENTIFIER, Literal: "[a-zA-Z]+[a-zA-Z0-9_]*", IsRegex: true},

	// Symbols
	{Type: PLUS, Literal: "+"},
	{Type: MINUS, Literal: "-"},
	{Type: ASTERISK, Literal: "*"},
	{Type: FORWARD_SLASH, Literal: "/"},
	{Type: SEMICOLON, Literal: ";"},
	{Type: OPEN_PAREN, Literal: "("},
	{Type: CLOSED_PAREN, Literal: ")"},
	{Type: ASSIGN, Literal: "="},
	{Type: COMMA, Literal: ","},
	{Type: OPEN_CURLY_BRACKET, Literal: "{"},
	{Type: CLOSED_CURLY_BRACKET, Literal: "}"},
	{Type: PTR, Literal: "<-"},
	{Type: OPEN_BRACKET, Literal: "["},
	{Type: CLOSED_BRACKET, Literal: "]"},
	{Type: AT, Literal: "@"},
	{Type: INLINE_COMMENT, Literal: "#"},
	{Type: BLOCK_COMMENT, Literal: "##"},
	{Type: EQ, Literal: "=="},
	{Type: EOF, Literal: ""},
}

func getToken(name string) Token {
	for _, token := range tokenData {
		if token.Type == name {
			return Token{Type: token.Type, Literal: token.Literal}
		}
	}
	panic(fmt.Sprintf("No token matching name: %s", name))
}

func GetTokenType(name string) string {
	token := getToken(name)
	return token.Type
}

func GetKeywordToken(literal string) Token {
	for _, token := range tokenData {
		if token.Literal == literal {
			return Token{Type: token.Type, Literal: token.Literal}
		}
	}

	identifierToken := getToken(IDENTIFIER)
	identifierToken.Literal = literal
	return identifierToken
}

func TokenTypesEqual(first Token, second Token) bool {
	return first.Type == second.Type
}

func (t *Token) ErrorDisplay() string {
	// How tokens should be displayed in error messages
	return fmt.Sprintf("%s (%#v)", t.Type, t.Literal)
}
