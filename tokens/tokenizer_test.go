package tokens

import "testing"

func TestSymbols(t *testing.T) {
	tokenizer := New("+-*/();")
	expectedTokens := []Token{
		{Literal: "+", Type: PLUS},
		{Literal: "-", Type: MINUS},
		{Literal: "*", Type: ASTERISK},
		{Literal: "/", Type: FORWARD_SLASH},
		{Literal: "(", Type: OPEN_PAREN},
		{Literal: ")", Type: CLOSED_PAREN},
		{Literal: ";", Type: SEMICOLON},
	}

	for _, expectedToken := range expectedTokens {
		actualToken := tokenizer.Next()

		if expectedToken.Literal != actualToken.Literal {
			t.Fatalf("Expected Literal: %s, Actual Literal: %s", expectedToken.Literal, actualToken.Literal)
		}

		if expectedToken.Type != actualToken.Type {
			t.Fatalf("Expected Type: %s, Actual Type: %s", expectedToken.Type, actualToken.Type)
		}
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
	}

	for i, source := range numbers {
		tokenizer := New(source)
		token := tokenizer.Next()

		if token.Literal != source {
			t.Fatalf("[Test #%d] Expected Literal: %s, Actual Literal: %s", i, source, token.Literal)
		}

		if token.Type != NUMBER {
			t.Fatalf("[Test #%d] Expected Type: %s, Actual Type: %s", i, NUMBER, token.Type)
		}
	}
}
