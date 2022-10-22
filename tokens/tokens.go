package tokens

import (
	"fmt"
	"regexp"
)

type TokenMetaData struct {
	Literal     string
	Type        string
	IsRegexChar bool // Any characters in Boomerang that are the same as regex characters need to be escaped.
	IsKeyword   bool
}

func (tmd *TokenMetaData) RegexPattern() string {
	if tmd.IsRegexChar {
		return regexp.QuoteMeta(tmd.Literal)
	}
	return tmd.Literal
}

type Token struct {
	Literal    string
	Type       string
	LineNumber int
}

func (t *Token) ErrorDisplay() string {
	// How tokens should be displayed in error messages
	return fmt.Sprintf("%s (%#v)", t.Type, t.Literal)
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
	FUNCTION             = "FUNCTION"
	NUMBER               = "NUMBER"
	STRING               = "STRING"
	BOOLEAN              = "BOOLEAN"
	IDENTIFIER           = "IDENTIFIER"
	EOF                  = "EOF"
	SEND                 = "SEND"
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
	LT                   = "LESS_THAN"
	WHEN                 = "WHEN"
	IS                   = "IS"
	OR                   = "OR"
	AND                  = "AND"
	FOR                  = "FOR"
	IN                   = "IN"
	WHILE                = "WHILE"
	BREAK                = "BREAK"
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
	SEND_TOKEN                 = getToken(SEND)
	OPEN_BRACKET_TOKEN         = getToken(OPEN_BRACKET)
	CLOSED_BRACKET_TOKEN       = getToken(CLOSED_BRACKET)
	AT_TOKEN                   = getToken(AT)
	INLINE_COMMENT_TOKEN       = getToken(INLINE_COMMENT)
	BLOCK_COMMENT_TOKEN        = getToken(BLOCK_COMMENT)
	EQ_TOKEN                   = getToken(EQ)
	LT_TOKEN                   = getToken(LT)

	// Keywords
	FUNCTION_TOKEN = getToken(FUNCTION)
	TRUE_TOKEN     = Token{Type: BOOLEAN, Literal: "true"}
	FALSE_TOKEN    = Token{Type: BOOLEAN, Literal: "false"}
	ELSE_TOKEN     = getToken(ELSE)
	WHEN_TOKEN     = getToken(WHEN)
	IS_TOKEN       = getToken(IS)
	NOT_TOKEN      = getToken(NOT)
	OR_TOKEN       = getToken(OR)
	AND_TOKEN      = getToken(AND)
	FOR_TOKEN      = getToken(FOR)
	IN_TOKEN       = getToken(IN)
	WHILE_TOKEN    = getToken(WHILE)
	BREAK_TOKEN    = getToken(BREAK)

	// Data Types
	NUMBER_TOKEN  = getToken(NUMBER)
	STRING_TOKEN  = getToken(STRING)
	BOOLEAN_TOKEN = getToken(BOOLEAN)

	// Misc
	IDENTIFIER_TOKEN = getToken(IDENTIFIER)
	EOF_TOKEN        = Token{Type: EOF, Literal: ""} // end of file
)

var tokenData = []TokenMetaData{
	// Data types/misc
	{Type: NUMBER, Literal: "[0-9]*[.]?[0-9]+"},
	{Type: STRING, Literal: "\"(.*?)\""},
	{Type: BOOLEAN, Literal: "(true|false)"},

	// Keywords. Need to be defined before "IDENTIFIER" in this list so they are not misclassified
	{Type: WHEN, Literal: "when", IsKeyword: true},
	{Type: IS, Literal: "is", IsKeyword: true},
	{Type: NOT, Literal: "not", IsKeyword: true},
	{Type: OR, Literal: "or", IsKeyword: true},
	{Type: AND, Literal: "and", IsKeyword: true},
	{Type: ELSE, Literal: "else", IsKeyword: true},
	{Type: FUNCTION, Literal: "func", IsKeyword: true},
	{Type: FOR, Literal: "for", IsKeyword: true},
	{Type: IN, Literal: "in", IsKeyword: true},
	{Type: BOOLEAN, Literal: "true", IsKeyword: true},
	{Type: BOOLEAN, Literal: "false", IsKeyword: true},
	{Type: WHILE, Literal: "while", IsKeyword: true},
	{Type: BREAK, Literal: "break", IsKeyword: true},

	// Identifier
	{Type: IDENTIFIER, Literal: "[a-zA-Z]+[a-zA-Z0-9_]*"},

	/*
		Symbols

		NOTE: If a token's literal value contains a substring of another token, that tokens must be declared before
		the substring tokens in this list. For example "==" must come before "=" and "##" must come before "#". Otherwise,
		the tokenizer would match "==" as two "=" tokens.
	*/
	{Type: PLUS, Literal: "+", IsRegexChar: true},
	{Type: MINUS, Literal: "-"},
	{Type: ASTERISK, Literal: "*", IsRegexChar: true},
	{Type: FORWARD_SLASH, Literal: "/", IsRegexChar: true},
	{Type: SEMICOLON, Literal: ";"},
	{Type: OPEN_PAREN, Literal: "(", IsRegexChar: true},
	{Type: CLOSED_PAREN, Literal: ")", IsRegexChar: true},
	{Type: SEND, Literal: "<-"},
	{Type: EQ, Literal: "=="},
	{Type: LT, Literal: "<"},
	{Type: ASSIGN, Literal: "="},
	{Type: COMMA, Literal: ","},
	{Type: OPEN_CURLY_BRACKET, Literal: "{", IsRegexChar: true},
	{Type: CLOSED_CURLY_BRACKET, Literal: "}", IsRegexChar: true},
	{Type: OPEN_BRACKET, Literal: "[", IsRegexChar: true},
	{Type: CLOSED_BRACKET, Literal: "]", IsRegexChar: true},
	{Type: AT, Literal: "@"},
	{Type: BLOCK_COMMENT, Literal: "##"},
	{Type: INLINE_COMMENT, Literal: "#"},
}

var keywords = getKeywords()

func getKeywords() map[string]string {

	keywordTokens := map[string]string{}

	for _, tokenDatum := range tokenData {
		if tokenDatum.IsKeyword {
			keywordTokens[tokenDatum.Literal] = tokenDatum.Type
		}
	}
	return keywordTokens
}

func getToken(name string) Token {
	for _, token := range tokenData {
		if token.Type == name {
			return Token{Type: token.Type, Literal: token.Literal}
		}
	}
	panic(fmt.Sprintf("No token matching name: %s", name))
}

func TokenTypesEqual(t Token, tokenType string) bool {
	return t.Type == tokenType
}
