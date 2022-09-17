package parser

import (
	"my_lang/tokens"
	"testing"
)

func TestNumber(t *testing.T) {
	tokenizer := tokens.New("10;")
	parser := New(tokenizer)

	actualAst := parser.Parse()
	expectedAST := []Node{
		{
			Type:  "Number",
			Value: "10",
		},
	}

	if len(actualAst) != len(expectedAST) {
		t.Fatalf("Expected number of statements: %d, Actual number of statements: %d", len(expectedAST), len(actualAst))
	}

	for i := range expectedAST {
		expected := expectedAST[i]
		actual := actualAst[i]

		if !assertNodeEqual(t, expected, actual) {
			t.Fatalf("Expected Node does not equal Actual node")
		}
	}
}

func TestBinaryExpression(t *testing.T) {
	tokenizer := tokens.New("7 + 3;")
	parser := New(tokenizer)

	actualAST := parser.Parse()
	expectedAST := []Node{
		{
			Type: "BinaryExpression",
			Params: map[string]Node{
				"left":     {Type: "Number", Value: "7"},
				"right":    {Type: "Number", Value: "3"},
				"operator": {Type: "PLUS", Value: "+"},
			},
		},
	}
	if len(actualAST) != len(expectedAST) {
		t.Fatalf("Expected number of statements: %d, Actual number of statements: %d", len(expectedAST), len(actualAST))
	}

	for i := range expectedAST {
		expected := expectedAST[i]
		actual := actualAST[i]

		if !assertNodeEqual(t, expected, actual) {
			t.Fatalf("Expected Node does not equal Actual node")
		}
	}
}

func TestParentheses(t *testing.T) {
	tokenizer := tokens.New("7 + (3);")
	parser := New(tokenizer)

	actualAST := parser.Parse()
	expectedAST := []Node{
		{
			Type: "BinaryExpression",
			Params: map[string]Node{
				"left":     {Type: "Number", Value: "7"},
				"operator": {Type: "PLUS", Value: "+"},
				"right":    {Type: "Number", Value: "3"},
			},
		},
	}
	if len(actualAST) != len(expectedAST) {
		t.Fatalf("Expected number of statements: %d, Actual number of statements: %d", len(expectedAST), len(actualAST))
	}

	for i := range expectedAST {
		expected := expectedAST[i]
		actual := actualAST[i]

		if !assertNodeEqual(t, expected, actual) {
			t.Fatalf("Expected Node does not equal Actual node")
		}
	}
}

func TestParenthesesBinaryExpression(t *testing.T) {
	tokenizer := tokens.New("7 + (5 - 2);")
	parser := New(tokenizer)

	actualAST := parser.Parse()
	expectedAST := []Node{
		{
			Type: "BinaryExpression",
			Params: map[string]Node{
				"left":     {Type: "Number", Value: "7"},
				"operator": {Type: "PLUS", Value: "+"},
				"right": {Type: "BinaryExpression", Params: map[string]Node{
					"left":     {Type: "Number", Value: "5"},
					"right":    {Type: "Number", Value: "3"},
					"operator": {Type: "MINUS", Value: "-"},
				}},
			},
		},
	}
	if len(actualAST) != len(expectedAST) {
		t.Fatalf("Expected number of statements: %d, Actual number of statements: %d", len(expectedAST), len(actualAST))
	}

	for i := range expectedAST {
		expected := expectedAST[i]
		actual := actualAST[i]

		if !assertNodeEqual(t, expected, actual) {
			t.Fatalf("Expected Node does not equal Actual node")
		}
	}
}

func assertNodeEqual(t *testing.T, expected Node, actual Node) bool {
	if expected.Type != actual.Type {
		t.Fatalf("Expected type: %s, Actual type: %s", expected.Type, actual.Type)
		return false
	}

	if expected.Value != actual.Value {
		t.Fatalf("Expected value: %s, Actual value: %s", expected.Value, actual.Value)
		return false
	}

	keys := make([]string, 0, len(expected.Params))
	for k := range expected.Params {
		keys = append(keys, k)
	}

	for _, key := range keys {
		expectedParamNode := expected.GetParam(key)
		actualParamNode := actual.GetParam(key)
		return assertNodeEqual(t, expectedParamNode, actualParamNode)
	}

	return true
}
