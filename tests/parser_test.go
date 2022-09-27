package tests

import (
	"boomerang/node"
	"boomerang/parser"
	"boomerang/tokens"
	"fmt"
	"testing"
)

func TestParser_Numbers(t *testing.T) {
	numbers := []string{
		"10",
		"1001",
		"5.5",
		"3.14159",
	}

	for _, number := range numbers {
		actualAST := getAST(fmt.Sprintf("%s;", number))
		expectedAST := []node.Node{
			CreateNumber(number),
		}

		AssertNodesEqual(t, expectedAST, actualAST)
	}
}

func TestParser_Booleans(t *testing.T) {
	tests := []struct {
		Source       string
		ExpectedNode node.Node
	}{
		{
			"true",
			CreateBooleanTrue(),
		},
		{
			"false",
			CreateBooleanFalse(),
		},
	}

	for _, test := range tests {
		actualAST := getAST(fmt.Sprintf("%s;", test.Source))
		expectedAST := []node.Node{
			test.ExpectedNode,
		}

		AssertNodesEqual(t, expectedAST, actualAST)
	}
}

func TestParser_Strings(t *testing.T) {

	plusToken := tokens.PLUS_TOKEN
	plusToken.LineNumber = 1

	tests := []struct {
		InputSource  string
		OutputSource string
		Params       []node.Node
	}{
		{
			InputSource:  "hello, world!",
			OutputSource: "hello, world!",
			Params:       []node.Node{},
		},
		{
			InputSource:  "My age is {55}",
			OutputSource: "My age is <0>",
			Params: []node.Node{
				CreateNumber("55"),
			},
		},
		{
			InputSource:  "{1 + 1} {1 + 1} {1 + 1}",
			OutputSource: "<0> <1> <2>",
			Params: []node.Node{
				node.CreateBinaryExpression(
					CreateNumber("1"),
					plusToken,
					CreateNumber("1"),
				),
				node.CreateBinaryExpression(
					CreateNumber("1"),
					plusToken,
					CreateNumber("1"),
				),
				node.CreateBinaryExpression(
					CreateNumber("1"),
					plusToken,
					CreateNumber("1"),
				),
			},
		},
	}

	for _, test := range tests {
		actualAST := getAST(fmt.Sprintf("\"%s\";", test.InputSource))
		expectedAST := []node.Node{
			node.CreateString(1, test.OutputSource, test.Params),
		}
		AssertNodesEqual(t, expectedAST, actualAST)
	}
}

func TestParser_TestParameters(t *testing.T) {

	tests := []struct {
		Source string
		Params []node.Node
	}{
		{
			Source: "()",
			Params: []node.Node{},
		},
		{
			Source: "(1,)",
			Params: []node.Node{
				CreateNumber("1"),
			},
		},
		{
			Source: "(1,2,)",
			Params: []node.Node{
				CreateNumber("1"),
				CreateNumber("2"),
			},
		},
		{
			Source: "(1,2,3)",
			Params: []node.Node{
				CreateNumber("1"),
				CreateNumber("2"),
				CreateNumber("3"),
			},
		},
		{
			Source: "(1,2,(4, 5, 6),3)",
			Params: []node.Node{
				CreateNumber("1"),
				CreateNumber("2"),
				CreateList([]node.Node{
					CreateNumber("4"),
					CreateNumber("5"),
					CreateNumber("6"),
				}),
				CreateNumber("3"),
			},
		},
	}

	for _, test := range tests {
		actualAST := getAST(fmt.Sprintf("%s;", test.Source))
		expectedAST := []node.Node{
			CreateList(test.Params),
		}

		AssertNodesEqual(t, expectedAST, actualAST)
	}
}

func TestParser_NegativeNumber(t *testing.T) {
	actualAST := getAST("-66;")
	expectedAST := []node.Node{
		node.CreateUnaryExpression(
			CreateTokenFromToken(tokens.MINUS_TOKEN),
			CreateNumber("66"),
		),
	}

	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_BinaryExpression(t *testing.T) {
	actualAST := getAST("7 + 3;")
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			CreateNumber("7"),
			CreateTokenFromToken(tokens.PLUS_TOKEN),
			CreateNumber("3"),
		),
	}

	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_Parentheses(t *testing.T) {
	actualAST := getAST("7 + (3);")
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			CreateNumber("7"),
			CreateTokenFromToken(tokens.PLUS_TOKEN),
			CreateNumber("3"),
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_ParenthesesBinaryExpression(t *testing.T) {

	actualAST := getAST("7 + (5 - 2);")
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			CreateNumber("7"),
			CreateTokenFromToken(tokens.PLUS_TOKEN),
			node.CreateBinaryExpression(
				CreateNumber("5"),
				CreateTokenFromToken(tokens.MINUS_TOKEN),
				CreateNumber("2"),
			),
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_VariableAssignment(t *testing.T) {
	actualAST := getAST("variable = 8 / 2;")
	expectedAST := []node.Node{
		CreateAssignmentStatement(
			"variable",
			node.CreateBinaryExpression(
				CreateNumber("8"),
				CreateTokenFromToken(tokens.FORWARD_SLASH_TOKEN),
				CreateNumber("2"),
			),
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_Identifier(t *testing.T) {
	actualAST := getAST("variable;")
	expectedAST := []node.Node{
		CreateIdentifier("variable"),
	}

	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_PrintStatement(t *testing.T) {
	actualAST := getAST("print(1, 2, variable);")
	expectedAST := []node.Node{
		CreatePrintStatement(
			[]node.Node{
				CreateNumber("1"),
				CreateNumber("2"),
				CreateIdentifier("variable"),
			},
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_PrintStatementNoArguments(t *testing.T) {
	actualAST := getAST("print();")
	expectedAST := []node.Node{
		CreatePrintStatement([]node.Node{}),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_Function(t *testing.T) {
	actualAST := getAST("func(a, b) { a + b; };")
	expectedAST := []node.Node{
		CreateFunction(
			[]node.Node{
				CreateIdentifier("a"),
				CreateIdentifier("b"),
			},
			[]node.Node{
				node.CreateBinaryExpression(
					CreateIdentifier("a"),
					CreateTokenFromToken(tokens.PLUS_TOKEN),
					CreateIdentifier("b"),
				),
			},
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_FunctionNoParameters(t *testing.T) {
	actualAST := getAST("func() { 3 + 4; };")
	expectedAST := []node.Node{
		CreateFunction(
			[]node.Node{},
			[]node.Node{
				node.CreateBinaryExpression(
					CreateNumber("3"),
					CreateTokenFromToken(tokens.PLUS_TOKEN),
					CreateNumber("4"),
				),
			},
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_FunctionNoParametersNoStatements(t *testing.T) {
	actualAST := getAST("func() {};")
	expectedAST := []node.Node{
		CreateFunction(
			[]node.Node{},
			[]node.Node{},
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_FunctionCallWithNoParameters(t *testing.T) {
	actualAST := getAST("divide <- ();")
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			CreateIdentifier("divide"),
			CreateTokenFromToken(tokens.PTR_TOKEN),
			CreateList([]node.Node{}),
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_FunctionCallWithFunctionLiteralAndLeftPointer(t *testing.T) {
	actualAST := getAST("func(c, d) { d - c; } <- (10, 2);")

	functionNode := CreateFunction(
		[]node.Node{
			CreateIdentifier("c"),
			CreateIdentifier("d"),
		},
		[]node.Node{
			node.CreateBinaryExpression(
				CreateIdentifier("d"),
				CreateTokenFromToken(tokens.MINUS_TOKEN),
				CreateIdentifier("c"),
			),
		},
	)

	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			functionNode,
			CreateTokenFromToken(tokens.PTR_TOKEN),
			CreateList([]node.Node{
				CreateNumber("10"),
				CreateNumber("2"),
			}),
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_FunctionCallWithIdentifierAndLeftPointer(t *testing.T) {
	actualAST := getAST("multiply <- (10, 3);")
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			CreateIdentifier("multiply"),
			CreateTokenFromToken(tokens.PTR_TOKEN),
			CreateList([]node.Node{
				CreateNumber("10"),
				CreateNumber("3"),
			}),
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_ReturnStatements(t *testing.T) {
	actualAST := getAST("return 1 + 1;")
	expectedAST := []node.Node{
		CreateReturnStatement(node.CreateBinaryExpression(
			CreateNumber("1"),
			CreateTokenFromToken(tokens.PLUS_TOKEN),
			CreateNumber("1"),
		)),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_ListIndex(t *testing.T) {
	actualAST := getAST("numbers[1];")
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			CreateIdentifier("numbers"),
			CreateTokenFromToken(tokens.OPEN_BRACKET_TOKEN),
			CreateNumber("1"),
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_IfStatement(t *testing.T) {
	actualAST := getAST("if true { print(\"true!!!\"); };")
	expectedAST := []node.Node{
		CreateIfStatement(
			CreateBooleanTrue(),
			[]node.Node{
				CreatePrintStatement([]node.Node{
					CreateRawString("true!!!"),
				}),
			},
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_UnexpectedTokenError(t *testing.T) {
	tokenizer := tokens.New("1")
	p, err := parser.New(tokenizer)
	if err != nil {
		panic(err.Error())
	}

	_, err = p.Parse()

	actualError := err.Error()
	expectedError := "error at line 1: expected token type SEMICOLON (\";\"), got EOF (\"\")"

	if expectedError != actualError {
		t.Fatalf("Expected error: %#v, Actual Error: %#v", expectedError, actualError)
	}
}

func getAST(source string) []node.Node {
	t := tokens.New(source)

	p, err := parser.New(t)
	if err != nil {
		panic(err.Error())
	}

	ast, err := p.Parse()
	if err != nil {
		panic(err.Error())
	}
	return *ast
}
