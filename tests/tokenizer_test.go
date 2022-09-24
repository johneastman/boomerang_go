package tests

import (
	"boomerang/tokens"
	"fmt"
	"testing"
)

func TestTokenizer_Symbols(t *testing.T) {
	tokenizer := tokens.New("+-*/()=,{}<-;")
	expectedTokens := []tokens.Token{
		tokens.PLUS_TOKEN,
		tokens.MINUS_TOKEN,
		tokens.ASTERISK_TOKEN,
		tokens.FORWARD_SLASH_TOKEN,
		tokens.OPEN_PAREN_TOKEN,
		tokens.CLOSED_PAREN_TOKEN,
		tokens.ASSIGN_TOKEN,
		tokens.COMMA_TOKEN,
		tokens.OPEN_CURLY_BRACKET_TOKEN,
		tokens.CLOSED_CURLY_BRACKET_TOKEN,
		tokens.PTR_TOKEN,
		tokens.SEMICOLON_TOKEN,
	}

	for _, expectedToken := range expectedTokens {
		actualToken, _ := tokenizer.Next()
		AssertTokenEqual(t, expectedToken, *actualToken)
	}
}

func TestTokenizer_Keywords(t *testing.T) {

	keywordTokens := []tokens.Token{
		tokens.PRINT_TOKEN,
		tokens.FUNCTION_TOKEN,
		tokens.TRUE_TOKEN,
		tokens.FALSE_TOKEN,
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

		AssertTokenEqual(t, tokens.Token{Type: tokens.NUMBER_TOKEN.Type, Literal: source}, *token)
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

		AssertTokenEqual(t, tokens.Token{Type: tokens.STRING_TOKEN.Type, Literal: testString}, *token)
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

		AssertTokenEqual(t, tokens.Token{Type: tokens.IDENTIFIER_TOKEN.Type, Literal: variable}, *token)
	}
}
