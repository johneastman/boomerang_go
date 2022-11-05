package tests

import (
	"boomerang/evaluator"
	"boomerang/node"
	"boomerang/parser"
	"boomerang/tokens"
	"fmt"
	"io"
	"os"
	"testing"
)

const TEST_LINE_NUM = 1

/* * * * * * * * * *
 * Test Data Utils *
 * * * * * * * * * */

func getEvaluatorResults(ast []node.Node) []node.Node {
	evaluatorObj := evaluator.NewEvaluator(ast)
	actualResults, err := evaluatorObj.Evaluate()
	if err != nil {
		panic(err.Error())
	}
	return actualResults
}

func getEvaluatorError(t *testing.T, ast []node.Node) string {
	evaluatorObj := evaluator.NewEvaluator(ast)
	_, err := evaluatorObj.Evaluate()

	if err == nil {
		t.Fatal("error is nil")
	}
	return err.Error()
}

func getParserAST(source string) []node.Node {
	t := tokens.NewTokenizer(source)

	p, err := parser.NewParser(t)
	if err != nil {
		panic(err.Error())
	}

	ast, err := p.Parse()
	if err != nil {
		panic(err.Error())
	}
	return *ast
}

func getParserError(t *testing.T, source string) string {
	tokenizer := tokens.NewTokenizer(source)

	p, err := parser.NewParser(tokenizer)
	if err != nil {
		panic(err.Error())
	}

	_, err = p.Parse()
	if err == nil {
		t.Fatalf("Expected error to not be nil")
	}
	return err.Error()
}

func CreateTokenFromToken(token tokens.Token) tokens.Token {
	return tokens.Token{Type: token.Type, Literal: token.Literal, LineNumber: TEST_LINE_NUM}
}

func CreateTokenWithLineNum(token tokens.Token, lineNum int) tokens.Token {
	return tokens.Token{Type: token.Type, Literal: token.Literal, LineNumber: lineNum}
}

func CreateTokenFromValues(tokenType, literal string) tokens.Token {
	token := tokens.Token{Type: tokenType, Literal: literal}
	return CreateTokenWithLineNum(token, TEST_LINE_NUM)
}

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

func CreateMonad(value *node.Node) node.Node {
	return node.CreateMonad(TEST_LINE_NUM, value)
}

func CreateIdentifier(value string) node.Node {
	return node.CreateIdentifier(TEST_LINE_NUM, value)
}

func CreateBuiltinFunctionIdentifier(value string) node.Node {
	return node.CreateBuiltinFunctionIdentifier(TEST_LINE_NUM, value)
}

func CreateBuiltinVariableIdentifier(value string) node.Node {
	return node.CreateBuiltinVariableIdentifier(TEST_LINE_NUM, value)
}

func CreateList(values []node.Node) node.Node {
	return node.CreateList(TEST_LINE_NUM, values)
}

func CreateAssignmentNode(variable node.Node, value node.Node) node.Node {
	return node.CreateAssignmentNode(variable, value)
}

func CreateFunction(parameters []node.Node, statements []node.Node) node.Node {
	return node.CreateFunction(
		TEST_LINE_NUM,
		parameters,
		CreateBlockStatements(statements),
	)
}

func CreateFunctionCall(function node.Node, callParams []node.Node) node.Node {
	return node.CreateFunctionCall(TEST_LINE_NUM, function, callParams)
}

func CreateBlockStatementReturnValue(statement *node.Node) node.Node {
	return node.CreateBlockStatementReturnValue(TEST_LINE_NUM, statement)
}

func CreateBlockStatements(statements []node.Node) node.Node {
	return node.CreateBlockStatements(statements)
}

func CreateWhenNode(whenExpression node.Node, cases []node.Node, defaultStatements []node.Node) node.Node {
	return node.CreateWhenNode(
		TEST_LINE_NUM,
		whenExpression,
		cases,
		CreateBlockStatements(defaultStatements),
	)
}

func CreateWhenCaseNode(expression node.Node, statements []node.Node) node.Node {
	return node.CreateCaseNode(
		TEST_LINE_NUM,
		expression,
		CreateBlockStatements(statements),
	)
}

func CreateForLoop(placeholder node.Node, list node.Node, statements []node.Node) node.Node {
	return node.CreateForLoop(
		TEST_LINE_NUM,
		placeholder,
		list,
		CreateBlockStatements(statements),
	)
}

func CreateWhileLoop(condition node.Node, statements []node.Node) node.Node {
	return node.CreateWhileLoop(
		TEST_LINE_NUM,
		condition,
		CreateBlockStatements(statements),
	)
}

func CreateBreakStatement() node.Node {
	return node.CreateBreakStatement(TEST_LINE_NUM)
}

func CreateContinueStatement() node.Node {
	return node.CreateContinueStatement(TEST_LINE_NUM)
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
		return fmt.Errorf(
			"expected length: %d, actual length: %d",
			len(expectedNodes),
			len(actualNodes),
		)
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
		return fmt.Errorf(
			"expected line number: %d (%s), actual line number: %d (%s)",
			expected.LineNum,
			expected.String(),
			actual.LineNum,
			actual.String(),
		)
	}

	if len(expected.Params) != len(actual.Params) {
		return fmt.Errorf(
			"expected %d (%s) params, got %d (%s)",
			len(expected.Params),
			expected.String(),
			len(actual.Params),
			actual.String(),
		)
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

func AssertExpectedOutput(t *testing.T, testNumber int, expectedOutput string, f func()) {
	testName := fmt.Sprintf("Test #%d", testNumber)

	t.Run(testName, func(t *testing.T) {
		actualOutput, err := assertExpectedOutput(expectedOutput, f)
		if err != nil {
			t.Fatal(err)
		}
		if expectedOutput != *actualOutput {
			t.Fatalf("Expected %#v, got %#v", expectedOutput, *actualOutput)
		}
	})
}

func assertExpectedOutput(expectedOutput string, f func()) (*string, error) {
	rescueStdout := os.Stdout

	defer func() {
		// Reset STDOUT after function runs/if any errors occur
		os.Stdout = rescueStdout
	}()

	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	os.Stdout = w

	// Execute code that should print to console
	f()

	w.Close()
	outputBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	actualOutput := string(outputBytes)
	return &actualOutput, nil
}

func AssertExpectedInput(t *testing.T, testNumber int, expectedOutput string, f func()) {
	testName := fmt.Sprintf("Test #%d", testNumber)

	t.Run(testName, func(t *testing.T) {
		if err := assertExpectedInput(expectedOutput, f); err != nil {
			t.Fatal(err)
		}
	})
}

func assertExpectedInput(inputString string, f func()) error {

	rescueStdin := os.Stdin

	defer func() {
		// Reset STDIN after function runs/if any errors occur
		os.Stdin = rescueStdin
	}()

	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	os.Stdin = r

	// Write to stdin
	input := []byte(inputString)
	_, err = w.Write(input)
	if err != nil {
		return err
	}
	w.Close()

	// Execute code that performs IO operations
	f()

	return nil
}

func AssertErrorEqual(t *testing.T, testNumber int, expected string, actual string) {
	testName := fmt.Sprintf("Test #%d", testNumber)

	t.Run(testName, func(t *testing.T) {
		if err := assertErrorEqual(expected, actual); err != nil {
			t.Fatal(err.Error())
		}
	})
}

func assertErrorEqual(expected string, actual string) error {
	if expected != actual {
		return fmt.Errorf("expected error: %s; actual error: %s", expected, actual)
	}
	return nil
}
