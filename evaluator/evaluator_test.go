package evaluator

import (
	"my_lang/parser"
	"testing"
)

func TestNumber(t *testing.T) {
	ast := []parser.Node{
		{Type: "Number", Value: "5"},
	}
	evaluator := New(ast)

	actualResults := evaluator.Evaluate()
	expectedResults := []parser.Node{
		{Type: "Number", Value: "5"},
	}

	assertNodesEqual(t, expectedResults, actualResults)
}

func TestBinaryExpression(t *testing.T) {
	ast := []parser.Node{
		{
			Type: "BinaryExpression",
			Params: map[string]parser.Node{
				"left":     {Type: "Number", Value: "1"},
				"right":    {Type: "Number", Value: "1"},
				"operator": {Type: "PLUS", Value: "+"},
			},
		},
	}
	evaluator := New(ast)

	actualResults := evaluator.Evaluate()
	expectedResults := []parser.Node{
		{Type: "Number", Value: "2"},
	}
	assertNodesEqual(t, actualResults, expectedResults)
}

func assertNodesEqual(t *testing.T, expectedNodes []parser.Node, actualNodes []parser.Node) bool {
	if len(expectedNodes) != len(actualNodes) {
		t.Fatalf("Expected length: %d, Actual length: %d", len(expectedNodes), len(actualNodes))
	}
	for i := range expectedNodes {
		expected := expectedNodes[i]
		actual := actualNodes[i]
		if !assertNodeEqual(t, expected, actual) {
			return false
		}
	}
	return true
}

func assertNodeEqual(t *testing.T, expected parser.Node, actual parser.Node) bool {
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
