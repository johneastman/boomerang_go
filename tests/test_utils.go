package tests

import (
	"boomerang/node"
	"boomerang/tokens"
	"io"
	"os"
	"testing"
)

func CreateTokenFromToken(token tokens.Token) tokens.Token {
	token.LineNumber = 1
	return token
}

func CreateTokenFromValues(type_ string, literal string, lineNum int) tokens.Token {
	return tokens.Token{Type: type_, Literal: literal, LineNumber: lineNum}
}

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

	if expected.LineNumber != actual.LineNumber {
		t.Fatalf("Expected Line Number: %d, Actual Line Number: %d", expected.LineNumber, actual.LineNumber)
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
