package tests

import (
	"boomerang/tokens"
	"fmt"
	"testing"
)

func TestTokenizer_Symbols(t *testing.T) {
	tokenizer := tokens.New("+-*/()=,{}<-[];")
	expectedTokens := []tokens.Token{
		CreateTokenFromToken(tokens.PLUS_TOKEN),
		CreateTokenFromToken(tokens.MINUS_TOKEN),
		CreateTokenFromToken(tokens.ASTERISK_TOKEN),
		CreateTokenFromToken(tokens.FORWARD_SLASH_TOKEN),
		CreateTokenFromToken(tokens.OPEN_PAREN_TOKEN),
		CreateTokenFromToken(tokens.CLOSED_PAREN_TOKEN),
		CreateTokenFromToken(tokens.ASSIGN_TOKEN),
		CreateTokenFromToken(tokens.COMMA_TOKEN),
		CreateTokenFromToken(tokens.OPEN_CURLY_BRACKET_TOKEN),
		CreateTokenFromToken(tokens.CLOSED_CURLY_BRACKET_TOKEN),
		CreateTokenFromToken(tokens.PTR_TOKEN),
		CreateTokenFromToken(tokens.OPEN_BRACKET_TOKEN),
		CreateTokenFromToken(tokens.CLOSED_BRACKET_TOKEN),
		CreateTokenFromToken(tokens.SEMICOLON_TOKEN),
	}

	for _, expectedToken := range expectedTokens {
		actualToken, _ := tokenizer.Next()
		AssertTokenEqual(t, expectedToken, *actualToken)
	}
}

func TestTokenizer_Keywords(t *testing.T) {

	keywordTokens := []tokens.Token{
		CreateTokenFromToken(tokens.PRINT_TOKEN),
		CreateTokenFromToken(tokens.FUNCTION_TOKEN),
		CreateTokenFromToken(tokens.TRUE_TOKEN),
		CreateTokenFromToken(tokens.FALSE_TOKEN),
		CreateTokenFromToken(tokens.IF_TOKEN),
	}

	for _, expectedToken := range keywordTokens {
		tokenizer := tokens.New(expectedToken.Literal)
		actualToken, _ := tokenizer.Next()
		AssertTokenEqual(t, expectedToken, *actualToken)
	}
}

func TestTokenizer_Numbers(t *testing.T) {
	numbers := []string{
		"1",
		"2",
		"3",
		"10",
		"100",
		"500",
		"10000",
		"9999999999",
		"1234567890",
		"9876543210",
		"1.1",
		".1",
		"1234567890.0987654321",
	}

	for _, source := range numbers {
		tokenizer := tokens.New(source)
		token, _ := tokenizer.Next()

		AssertTokenEqual(t, CreateTokenFromValues(tokens.NUMBER_TOKEN.Type, source, 1), *token)
	}
}

func TestTokenizer_Strings(t *testing.T) {

	testStrings := []string{
		"hello, world!",
		"abcdefghijklmnopqrstuvwxyz",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"1",
		"1234567890",
	}

	for _, testString := range testStrings {

		source := fmt.Sprintf("\"%s\";", testString)

		tokenizer := tokens.New(source)
		token, _ := tokenizer.Next()

		AssertTokenEqual(t, CreateTokenFromValues(tokens.STRING_TOKEN.Type, testString, 1), *token)
	}
}

func TestTokenizer_Identifiers(t *testing.T) {

	variables := []string{
		"variable",
		"varaible1",
		"variable_23",
		"_variable_",
	}

	for _, variable := range variables {
		tokenizer := tokens.New(variable)
		token, _ := tokenizer.Next()

		AssertTokenEqual(t, CreateTokenFromValues(tokens.IDENTIFIER_TOKEN.Type, variable, 1), *token)
	}
}

func TestTokenizer_InlineCommentsEOF(t *testing.T) {
	source := "# this is a comment"
	tokenizer := tokens.New(source)

	token, _ := tokenizer.Next()

	AssertTokenEqual(t, CreateTokenFromToken(tokens.EOF_TOKEN), *token)
}

func TestTokenizer_InlineComments(t *testing.T) {
	source := "# this is a comment\n1;"
	tokenizer := tokens.New(source)

	actualToken, _ := tokenizer.Next()
	expectedToken := tokens.Token{Type: tokens.NUMBER, Literal: "1", LineNumber: 2}

	AssertTokenEqual(t, expectedToken, *actualToken)
}
