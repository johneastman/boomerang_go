package tests

import (
	"boomerang/tokens"
	"fmt"
	"testing"
)

func TestTokenizer_Symbols(t *testing.T) {
	tokenizer := tokens.New("+-*/()=,{}<-[]==;")
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
		CreateTokenFromToken(tokens.EQ_TOKEN),
		CreateTokenFromToken(tokens.SEMICOLON_TOKEN),
	}

	for i, expectedToken := range expectedTokens {
		actualToken, _ := tokenizer.Next()
		AssertTokenEqual(t, i, expectedToken, *actualToken)
	}
}

func TestTokenizer_Keywords(t *testing.T) {
	keywordTokens := []tokens.Token{
		{Type: tokens.PRINT, Literal: "print", LineNumber: TEST_LINE_NUM},
		{Type: tokens.FUNCTION, Literal: "func", LineNumber: TEST_LINE_NUM},
		{Type: tokens.BOOLEAN, Literal: "true", LineNumber: TEST_LINE_NUM},
		{Type: tokens.BOOLEAN, Literal: "false", LineNumber: TEST_LINE_NUM},
		{Type: tokens.NOT, Literal: "not", LineNumber: TEST_LINE_NUM},
		{Type: tokens.WHEN, Literal: "when", LineNumber: TEST_LINE_NUM},
		{Type: tokens.IS, Literal: "is", LineNumber: TEST_LINE_NUM},
	}

	for i, expectedToken := range keywordTokens {
		tokenizer := tokens.New(expectedToken.Literal)
		actualToken, _ := tokenizer.Next()
		AssertTokenEqual(t, i, expectedToken, *actualToken)
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

	for i, source := range numbers {
		tokenizer := tokens.New(source)
		token, _ := tokenizer.Next()

		AssertTokenEqual(t, i, CreateTokenFromValues(tokens.NUMBER, source, 1), *token)
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

	for i, testString := range testStrings {

		source := fmt.Sprintf("\"%s\";", testString)

		tokenizer := tokens.New(source)
		token, _ := tokenizer.Next()

		AssertTokenEqual(t, i, CreateTokenFromValues(tokens.STRING, testString, 1), *token)
	}
}

func TestTokenizer_Identifiers(t *testing.T) {

	variables := []string{
		"variable",
		"varaible1",
		"variable_23",
	}

	for i, variable := range variables {
		tokenizer := tokens.New(variable)
		token, _ := tokenizer.Next()

		AssertTokenEqual(t, i, CreateTokenFromValues(tokens.IDENTIFIER, variable, 1), *token)
	}
}

func TestTokenizer_InlineCommentEOF(t *testing.T) {
	source := "# this is a comment"
	tokenizer := tokens.New(source)

	token, _ := tokenizer.Next()

	AssertTokenEqual(t, 0, CreateTokenFromToken(tokens.EOF_TOKEN), *token)
}

func TestTokenizer_InlineComment(t *testing.T) {
	source := "# this is a comment\n1;"
	tokenizer := tokens.New(source)

	actualToken, _ := tokenizer.Next()
	expectedToken := tokens.Token{Type: tokens.NUMBER, Literal: "1", LineNumber: 2}

	AssertTokenEqual(t, 0, expectedToken, *actualToken)
}

func TestTokenizer_BlockCommentEOF(t *testing.T) {
	source := "##\na = 1;\nb = a + 2\n##"
	tokenizer := tokens.New(source)

	actualToken, _ := tokenizer.Next()
	expectedToken := tokens.Token{Type: tokens.EOF, Literal: tokens.EOF_TOKEN.Literal, LineNumber: 4}

	AssertTokenEqual(t, 0, expectedToken, *actualToken)
}

func TestTokenizer_BlockComment(t *testing.T) {
	source := "##\na = 1;\nb = a + 2\n##1;"
	tokenizer := tokens.New(source)

	actualToken, _ := tokenizer.Next()
	expectedToken := tokens.Token{Type: tokens.NUMBER, Literal: "1", LineNumber: 4}

	AssertTokenEqual(t, 0, expectedToken, *actualToken)
}

func TestTokenizer_BlockCommentError(t *testing.T) {
	source := "##"
	tokenizer := tokens.New(source)

	_, err := tokenizer.Next()
	if err == nil {
		t.Fatal("An error was expected, but no errors occurred")
	}

	actualError := err.Error()
	expectedError := "error at line 1: did not find ending ## while parsing block comment"
	if expectedError != actualError {
		t.Fatalf("Expected error: %#v, Actual error: %#v", expectedError, actualError)
	}
}

func TestTokenizer_Booleans(t *testing.T) {
	/*
		This test is because I originally had the boolean regex "true|false", but after appending "^", the regex
		was interpreted as "^true|false", when it should have been "^(true|false)".
	*/
	expectedTokens := []tokens.Token{
		{Type: tokens.IDENTIFIER, Literal: "b", LineNumber: TEST_LINE_NUM},
		{Type: tokens.ASSIGN, Literal: "=", LineNumber: TEST_LINE_NUM},
		{Type: tokens.BOOLEAN, Literal: "false", LineNumber: TEST_LINE_NUM},
	}

	source := "b = false;"
	tokenizer := tokens.New(source)

	for i, expectedToken := range expectedTokens {
		token, _ := tokenizer.Next()
		AssertTokenEqual(t, i, expectedToken, *token)
	}
}
