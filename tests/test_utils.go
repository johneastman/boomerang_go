package tests

import (
	"boomerang/node"
	"boomerang/tokens"
	"io"
	"os"
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

	if len(expected.Params) != len(actual.Params) {
		t.Fatalf("Expected %d params, got %d", len(expected.Params), len(actual.Params))
	}

	for i := 0; i < len(expected.Params); i++ {
		expectedParamNode := expected.Params[i]
		actualParamNode := actual.Params[i]
		AssertNodeEqual(t, expectedParamNode, actualParamNode)
	}
}

func AssertExpectedOutput(t *testing.T, expectedOutput string, f func()) {
	rescueStdout := os.Stdout

	defer func() {
		// Reset STDOUT after function runs/if any errors occur
		os.Stdout = rescueStdout
	}()

	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute code that should print to console
	f()

	w.Close()
	actualOutput, _ := io.ReadAll(r)

	if expectedOutput != string(actualOutput) {
		t.Fatalf("Expected %#v, got %#v", "1 2 3\n", actualOutput)
	}
}

func CreateFunction(params []string, statements []node.Node) node.Node {

	// Create parameter nodes from the array of strings
	parameterNodes := []node.Node{}
	for _, value := range params {
		paramNode := node.Node{Type: node.IDENTIFIER, Value: value}
		parameterNodes = append(parameterNodes, paramNode)
	}

	return node.Node{
		Type: node.FUNCTION,
		Params: []node.Node{
			{Type: node.PARAMETER, Params: parameterNodes},
			{Type: node.STMTS, Params: statements},
		},
	}
}

func CreateFunctionCall(function node.Node, callParams []node.Node) node.Node {
	return node.Node{
		Type: node.FUNCTION_CALL,
		Params: []node.Node{
			{Type: node.CALL_PARAMS, Params: callParams},
			function,
		},
	}
}
