package tests

import (
	"boomerang/tokens"
	"testing"
)

func TestSymbols(t *testing.T) {
	tokenizer := tokens.New("+-*/()=,;")
	expectedTokens := []tokens.Token{
		{Literal: "+", Type: tokens.PLUS},
		{Literal: "-", Type: tokens.MINUS},
		{Literal: "*", Type: tokens.ASTERISK},
		{Literal: "/", Type: tokens.FORWARD_SLASH},
		{Literal: "(", Type: tokens.OPEN_PAREN},
		{Literal: ")", Type: tokens.CLOSED_PAREN},
		{Literal: "=", Type: tokens.ASSIGN},
		{Literal: ",", Type: tokens.COMMA},
		{Literal: ";", Type: tokens.SEMICOLON},
	}

	for _, expectedToken := range expectedTokens {
		actualToken := tokenizer.Next()
		AssertTokenEqual(t, expectedToken, actualToken)
	}
}

func TestTokenizerKeywords(t *testing.T) {

	keywordTokens := []tokens.Token{
		{Type: tokens.PRINT, Literal: "print"},
	}

	for _, expectedToken := range keywordTokens {
		tokenizer := tokens.New(expectedToken.Literal)
		actualToken := tokenizer.Next()
		AssertTokenEqual(t, expectedToken, actualToken)
	}
}

func TestNumbers(t *testing.T) {
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
		token := tokenizer.Next()

		AssertTokenEqual(t, tokens.Token{Type: tokens.NUMBER, Literal: source}, token)
	}
}

func TestIdenifiers(t *testing.T) {

	variables := []string{
		"variable",
		"varaible1",
		"variable_23",
		"_variable_",
	}

	for _, variable := range variables {
		tokenizer := tokens.New(variable)
		token := tokenizer.Next()

		AssertTokenEqual(t, tokens.Token{Type: tokens.IDENTIFIER, Literal: variable}, token)
	}
}
