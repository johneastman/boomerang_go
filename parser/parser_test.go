package parser

import (
	"my_lang/tokens"
	"testing"
)

func TestNumber(t *testing.T) {
	tokenizer := tokens.New("10;")
	parser := New(tokenizer)

	ast := parser.Parse()
	expectedAST := []Node{
		{
			Type: "Statement",
			Params: map[string]Node{
				"Expression": {
					Type:  "Number",
					Value: "10",
				},
			},
		},
	}

	if len(ast) != len(expectedAST) {
		t.Fatalf("Expected number of statements: %d, Actual number of statements: %d", len(expectedAST), len(ast))
	}

	for i, _ := range expectedAST {
		expected := expectedAST[i]
		actual := ast[i]

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
