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
