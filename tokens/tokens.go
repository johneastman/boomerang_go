package tokens

import (
	"fmt"
)

type Token struct {
	Literal    string
	Type       string
	LineNumber int
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
	STRING               = "STRING"
	BOOLEAN              = "BOOLEAN"
	IDENTIFIER           = "IDENTIFIER"
	EOF                  = "EOF"
	PTR                  = "POINTER"
	DOUBLE_QUOTE         = "DOUBLE_QUOTE"
	RETURN               = "RETURN"
	TRUE                 = "TRUE"
	FALSE                = "FALSE"
	OPEN_BRACKET         = "OPEN_BRACKET"
	CLOSED_BRACKET       = "CLOSED_BRACKET"
	IF                   = "IF"
	AT                   = "AT"
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
	DOUBLE_QUOTE_TOKEN         = getToken(DOUBLE_QUOTE)
	OPEN_BRACKET_TOKEN         = getToken(OPEN_BRACKET)
	CLOSED_BRACKET_TOKEN       = getToken(CLOSED_BRACKET)
	AT_TOKEN                   = getToken(AT)

	// Keywords
	PRINT_TOKEN    = getToken(PRINT)
	FUNCTION_TOKEN = getToken(FUNCTION)
	RETURN_TOKEN   = getToken(RETURN)
	TRUE_TOKEN     = getToken(TRUE)
	FALSE_TOKEN    = getToken(FALSE)
	IF_TOKEN       = getToken(IF)

	// Data Types
	NUMBER_TOKEN  = getToken(NUMBER)
	STRING_TOKEN  = getToken(STRING)
	BOOLEAN_TOKEN = getToken(BOOLEAN)

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
	PTR:                  {Type: "LEFT_POINTER", Literal: "<-"},
	NUMBER:               {Type: "NUMBER", Literal: ""},
	IDENTIFIER:           {Type: "IDENTIFIER", Literal: ""},
	DOUBLE_QUOTE:         {Type: "DOUBLE_QUOTE", Literal: "\""},
	STRING:               {Type: "STRING", Literal: ""},
	BOOLEAN:              {Type: "BOOLEAN", Literal: ""},
	RETURN:               {Type: "RETURN", Literal: "return"},
	TRUE:                 {Type: "BOOLEAN", Literal: "true"},
	FALSE:                {Type: "BOOLEAN", Literal: "false"},
	OPEN_BRACKET:         {Type: "OPEN_BRACKET", Literal: "["},
	CLOSED_BRACKET:       {Type: "CLOSED_BRACKET", Literal: "]"},
	IF:                   {Type: "IF", Literal: "if"},
	AT:                   {Type: "AT", Literal: "@"},
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
	"print":  PRINT_TOKEN,
	"func":   FUNCTION_TOKEN,
	"return": RETURN_TOKEN,
	"true":   TRUE_TOKEN,
	"false":  FALSE_TOKEN,
	"if":     IF_TOKEN,
}

func GetKeywordToken(literal string) Token {
	if token, ok := keywords[literal]; ok {
		return token
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
