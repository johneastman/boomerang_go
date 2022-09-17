package parser

import (
	"my_lang/tokens"
	"testing"
)

func TestNumber(t *testing.T) {
	tokenizer := tokens.New("10;")
	parser := New(tokenizer)

	ast := parser.Parse()
	expectedAST := []Statement{
		{Expr: &Number{Value: "10"}},
	}

	if len(ast) != len(expectedAST) {
		t.Fatalf("Expected number of statements: %d, Actual number of statements: %d", len(expectedAST), len(ast))
	}
	// TODO: Add tests for AST node content. First, refacor to use Node class, similar to C++ implementation
}
