package tests

import (
	"boomerang/node"
	"boomerang/tokens"
	"testing"
)

func AssertTokensEqual(t *testing.T, expectedTokens []tokens.Token, actualTokens []tokens.Token) {
	if len(expectedTokens) != len(actualTokens) {
		t.Fatalf("Expected length: %d, Actual length: %d\n", len(expectedTokens), len(actualTokens))
	}

	for i := range expectedTokens {
		expected := expectedTokens[i]
		actual := actualTokens[i]

		AssertTokenEqual(t, expected, actual)
	}
}

func AssertTokenEqual(t *testing.T, expected tokens.Token, actual tokens.Token) {
	if expected.Literal != actual.Literal {
		t.Fatalf("Expected Literal: %s, Actual Literal: %s\n", expected.Literal, actual.Literal)
	}

	if expected.Type != actual.Type {
		t.Fatalf("Expected Type: %s, Actual Type: %s\n", expected.Type, actual.Type)
	}
}

func AssertNodesEqual(t *testing.T, expectedNodes []node.Node, actualNodes []node.Node) {
	if len(expectedNodes) != len(actualNodes) {
		t.Fatalf("Expected length: %d, Actual length: %d\n", len(expectedNodes), len(actualNodes))
	}
	for i := range expectedNodes {
		expected := expectedNodes[i]
		actual := actualNodes[i]
		AssertNodeEqual(t, expected, actual)
	}
}

func AssertNodeEqual(t *testing.T, expected node.Node, actual node.Node) {
	if expected.Type != actual.Type {
		t.Fatalf("Expected type: %s, Actual type: %s\n", expected.Type, actual.Type)
	}

	if expected.Value != actual.Value {
		t.Fatalf("Expected value: %s, Actual value: %s\n", expected.Value, actual.Value)
	}

	keys := make([]string, 0, len(expected.Params))
	for k := range expected.Params {
		keys = append(keys, k)
	}

	for _, key := range keys {
		expectedParamNode := expected.GetParam(key)
		actualParamNode := actual.GetParam(key)
		AssertNodeEqual(t, expectedParamNode, actualParamNode)
	}
}
