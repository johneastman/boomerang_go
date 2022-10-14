package tests

import (
	"boomerang/node"
	"boomerang/tokens"
	"fmt"
	"io"
	"os"
	"testing"
)

const TEST_LINE_NUM = 1

func CreateNumber(value string) node.Node {
	return node.CreateNumber(TEST_LINE_NUM, value)
}

func CreateBoolean(value string) node.Node {
	return node.CreateBoolean(TEST_LINE_NUM, value)
}

func CreateBooleanTrue() node.Node {
	return node.CreateBooleanTrue(TEST_LINE_NUM)
}

func CreateBooleanFalse() node.Node {
	return node.CreateBooleanFalse(TEST_LINE_NUM)
}

func CreateString(value string, params []node.Node) node.Node {
	return node.CreateString(TEST_LINE_NUM, value, params)
}

func CreateRawString(value string) node.Node {
	return node.CreateRawString(TEST_LINE_NUM, value)
}

func CreateIdentifier(value string) node.Node {
	return node.CreateIdentifier(TEST_LINE_NUM, value)
}

func CreateList(values []node.Node) node.Node {
	return node.CreateList(TEST_LINE_NUM, values)
}

func CreatePrintStatement(params []node.Node) node.Node {
	return node.CreatePrintStatement(TEST_LINE_NUM, params)
}

func CreateAssignmentStatement(variableName string, value node.Node) node.Node {
	return node.CreateAssignmentStatement(TEST_LINE_NUM, variableName, value)
}

func CreateFunction(parameters []node.Node, statements []node.Node) node.Node {
	blockStatements := node.CreateBlockStatements(TEST_LINE_NUM, statements)
	return node.CreateFunction(TEST_LINE_NUM, parameters, blockStatements)
}

func CreateFunctionCall(function node.Node, callParams []node.Node) node.Node {
	return node.CreateFunctionCall(TEST_LINE_NUM, function, callParams)
}

func CreateBlockStatementReturnValue(statement *node.Node) node.Node {
	return node.CreateBlockStatementReturnValue(TEST_LINE_NUM, statement)
}

func CreateTokenFromToken(token tokens.Token) tokens.Token {
	return tokens.Token{Type: token.Type, Literal: token.Literal, LineNumber: TEST_LINE_NUM}
}

func CreateBlockStatements(statements []node.Node) node.Node {
	return node.CreateBlockStatements(TEST_LINE_NUM, statements)
}

func CreateWhenNode(whenExpression node.Node, cases []node.Node, defaultStatements node.Node) node.Node {
	return node.CreateWhenNode(TEST_LINE_NUM, whenExpression, cases, defaultStatements)
}

func CreateWhenCaseNode(expression node.Node, statements node.Node) node.Node {
	return node.CreateCaseNode(TEST_LINE_NUM, expression, statements)
}

func CreateTokenFromValues(type_ string, literal string) tokens.Token {
	return tokens.Token{Type: type_, Literal: literal, LineNumber: TEST_LINE_NUM}
}

func CreateForLoop(placeholder node.Node, list node.Node, statements node.Node) node.Node {
	return node.CreateForLoop(TEST_LINE_NUM, placeholder, list, statements)
}

func AssertTokenEqual(t *testing.T, testNumber int, expected tokens.Token, actual tokens.Token) {
	testName := fmt.Sprintf("Test #%d", testNumber)

	t.Run(testName, func(t *testing.T) {
		if err := assertTokenEqual(expected, actual); err != nil {
			t.Fatal(err.Error())
		}
	})
}

func assertTokenEqual(expected tokens.Token, actual tokens.Token) error {
	if expected.Literal != actual.Literal {
		return fmt.Errorf("expected literal: %s, actual literal: %s", expected.Literal, actual.Literal)
	}

	if expected.Type != actual.Type {
		return fmt.Errorf("expected type: %s, actual type: %s", expected.Type, actual.Type)
	}

	if expected.LineNumber != actual.LineNumber {
		return fmt.Errorf("expected line number: %d, actual line number: %d", expected.LineNumber, actual.LineNumber)
	}
	return nil
}

func AssertNodesEqual(t *testing.T, testNumber int, expectedNodes []node.Node, actualNodes []node.Node) {
	testName := fmt.Sprintf("Test #%d", testNumber)

	t.Run(testName, func(t *testing.T) {
		if err := assertNodesEqual(expectedNodes, actualNodes); err != nil {
			t.Fatal(err.Error())
		}
	})
}

func assertNodesEqual(expectedNodes []node.Node, actualNodes []node.Node) error {
	if len(expectedNodes) != len(actualNodes) {
		return fmt.Errorf("expected length: %d, actual length: %d", len(expectedNodes), len(actualNodes))
	}
	for i := range expectedNodes {
		expected := expectedNodes[i]
		actual := actualNodes[i]

		if err := assertNodeEqual(expected, actual); err != nil {
			return err
		}
	}
	return nil
}

func AssertNodeEqual(t *testing.T, testNumber int, expected node.Node, actual node.Node) {
	testName := fmt.Sprintf("Test #%d", testNumber)

	t.Run(testName, func(t *testing.T) {
		if err := assertNodeEqual(expected, actual); err != nil {
			t.Fatal(err.Error())
		}
	})
}

func assertNodeEqual(expected node.Node, actual node.Node) error {
	if expected.Type != actual.Type {
		return fmt.Errorf("expected type: %s, actual type: %s", expected.Type, actual.Type)
	}

	if expected.Value != actual.Value {
		return fmt.Errorf("expected value: %s, actual value: %s", expected.Value, actual.Value)
	}

	if expected.LineNum != actual.LineNum {
		return fmt.Errorf("expected line number: %d, actual line number: %d", expected.LineNum, actual.LineNum)
	}

	if len(expected.Params) != len(actual.Params) {
		return fmt.Errorf("expected %d params, got %d", len(expected.Params), len(actual.Params))
	}

	for i := 0; i < len(expected.Params); i++ {
		expectedParamNode := expected.Params[i]
		actualParamNode := actual.Params[i]
		if err := assertNodeEqual(expectedParamNode, actualParamNode); err != nil {
			return err
		}
	}
	return nil
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
