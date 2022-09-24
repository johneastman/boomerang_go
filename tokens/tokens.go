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
	STRING               = "STRING"
	BOOLEAN              = "BOOLEAN"
	IDENTIFIER           = "IDENTIFIER"
	EOF                  = "EOF"
	LEFT_PTR             = "LEFT_POINTER"
	RIGHT_PTR            = "RIGHT_POINTER"
	DOUBLE_QUOTE         = "DOUBLE_QUOTE"
	RETURN               = "RETURN"
	TRUE                 = "TRUE"
	FALSE                = "FALSE"
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
	LEFT_PTR_TOKEN             = getToken(LEFT_PTR)
	RIGHT_PTR_TOKEN            = getToken(RIGHT_PTR)
	DOUBLE_QUOTE_TOKEN         = getToken(DOUBLE_QUOTE)

	// Keywords
	PRINT_TOKEN    = getToken(PRINT)
	FUNCTION_TOKEN = getToken(FUNCTION)
	RETURN_TOKEN   = getToken(RETURN)
	TRUE_TOKEN     = getToken(TRUE)
	FALSE_TOKEN    = getToken(FALSE)

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
	LEFT_PTR:             {Type: "LEFT_POINTER", Literal: "<-"},
	RIGHT_PTR:            {Type: "RIGHT_POINTER", Literal: "->"},
	NUMBER:               {Type: "NUMBER", Literal: ""},
	IDENTIFIER:           {Type: "IDENTIFIER", Literal: ""},
	DOUBLE_QUOTE:         {Type: "DOUBLE_QUOTE", Literal: "\""},
	STRING:               {Type: "STRING", Literal: ""},
	BOOLEAN:              {Type: "BOOLEAN", Literal: ""},
	RETURN:               {Type: "RETURN", Literal: "return"},
	TRUE:                 {Type: "BOOLEAN", Literal: "true"},
	FALSE:                {Type: "BOOLEAN", Literal: "false"},
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
