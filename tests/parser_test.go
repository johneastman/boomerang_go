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
			node.CreateNumber(number),
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
			node.CreateBooleanTrue(),
		},
		{
			"false",
			node.CreateBooleanFalse(),
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
				node.CreateNumber("55"),
			},
		},
		{
			InputSource:  "{1 + 1} {1 + 1} {1 + 1}",
			OutputSource: "<0> <1> <2>",
			Params: []node.Node{
				node.CreateBinaryExpression(
					node.CreateNumber("1"),
					tokens.PLUS_TOKEN,
					node.CreateNumber("1"),
				),
				node.CreateBinaryExpression(
					node.CreateNumber("1"),
					tokens.PLUS_TOKEN,
					node.CreateNumber("1"),
				),
				node.CreateBinaryExpression(
					node.CreateNumber("1"),
					tokens.PLUS_TOKEN,
					node.CreateNumber("1"),
				),
			},
		},
	}

	for _, test := range tests {
		actualAST := getAST(fmt.Sprintf("\"%s\";", test.InputSource))
		expectedAST := []node.Node{
			node.CreateString(test.OutputSource, test.Params),
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
				node.CreateNumber("1"),
			},
		},
		{
			Source: "(1,2,)",
			Params: []node.Node{
				node.CreateNumber("1"),
				node.CreateNumber("2"),
			},
		},
		{
			Source: "(1,2,3)",
			Params: []node.Node{
				node.CreateNumber("1"),
				node.CreateNumber("2"),
				node.CreateNumber("3"),
			},
		},
		{
			Source: "(1,2,(4, 5, 6),3)",
			Params: []node.Node{
				node.CreateNumber("1"),
				node.CreateNumber("2"),
				node.CreateList([]node.Node{
					node.CreateNumber("4"),
					node.CreateNumber("5"),
					node.CreateNumber("6"),
				}),
				node.CreateNumber("3"),
			},
		},
	}

	for _, test := range tests {
		actualAST := getAST(fmt.Sprintf("%s;", test.Source))
		expectedAST := []node.Node{
			node.CreateList(test.Params),
		}

		AssertNodesEqual(t, expectedAST, actualAST)
	}
}

func TestParser_NegativeNumber(t *testing.T) {
	actualAST := getAST("-66;")
	expectedAST := []node.Node{
		node.CreateUnaryExpression(
			tokens.MINUS_TOKEN,
			node.CreateNumber("66"),
		),
	}

	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_BinaryExpression(t *testing.T) {
	actualAST := getAST("7 + 3;")
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			node.CreateNumber("7"),
			tokens.PLUS_TOKEN,
			node.CreateNumber("3"),
		),
	}

	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_Parentheses(t *testing.T) {
	actualAST := getAST("7 + (3);")
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			node.CreateNumber("7"),
			tokens.PLUS_TOKEN,
			node.CreateNumber("3"),
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_ParenthesesBinaryExpression(t *testing.T) {
	actualAST := getAST("7 + (5 - 2);")
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			node.CreateNumber("7"),
			tokens.PLUS_TOKEN,
			node.CreateBinaryExpression(
				node.CreateNumber("5"),
				tokens.MINUS_TOKEN,
				node.CreateNumber("2"),
			),
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_VariableAssignment(t *testing.T) {
	actualAST := getAST("variable = 8 / 2;")
	expectedAST := []node.Node{
		node.CreateAssignmentStatement(
			"variable",
			node.CreateBinaryExpression(
				node.CreateNumber("8"),
				tokens.FORWARD_SLASH_TOKEN,
				node.CreateNumber("2"),
			),
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_Identifier(t *testing.T) {
	actualAST := getAST("variable;")
	expectedAST := []node.Node{
		node.CreateIdentifier("variable"),
	}

	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_PrintStatement(t *testing.T) {
	actualAST := getAST("print(1, 2, variable);")
	expectedAST := []node.Node{
		node.CreatePrintStatement(
			[]node.Node{
				node.CreateNumber("1"),
				node.CreateNumber("2"),
				node.CreateIdentifier("variable"),
			},
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_PrintStatementNoArguments(t *testing.T) {
	actualAST := getAST("print();")
	expectedAST := []node.Node{
		node.CreatePrintStatement([]node.Node{}),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_Function(t *testing.T) {
	actualAST := getAST("func(a, b) { a + b; };")
	expectedAST := []node.Node{
		node.CreateFunction(
			[]node.Node{
				node.CreateIdentifier("a"),
				node.CreateIdentifier("b"),
			},
			[]node.Node{
				node.CreateBinaryExpression(
					node.CreateIdentifier("a"),
					tokens.PLUS_TOKEN,
					node.CreateIdentifier("b"),
				),
			},
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_FunctionNoParameters(t *testing.T) {
	actualAST := getAST("func() { 3 + 4; };")
	expectedAST := []node.Node{
		node.CreateFunction(
			[]node.Node{},
			[]node.Node{
				node.CreateBinaryExpression(
					node.CreateNumber("3"),
					tokens.PLUS_TOKEN,
					node.CreateNumber("4"),
				),
			},
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_FunctionNoParametersNoStatements(t *testing.T) {
	actualAST := getAST("func() {};")
	expectedAST := []node.Node{
		node.CreateFunction(
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
			node.CreateIdentifier("divide"),
			tokens.PTR_TOKEN,
			node.CreateList([]node.Node{}),
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_FunctionCallWithFunctionLiteralAndLeftPointer(t *testing.T) {
	actualAST := getAST("func(c, d) { d - c; } <- (10, 2);")

	functionNode := node.CreateFunction(
		[]node.Node{
			node.CreateIdentifier("c"),
			node.CreateIdentifier("d"),
		},
		[]node.Node{
			node.CreateBinaryExpression(
				node.CreateIdentifier("d"),
				tokens.MINUS_TOKEN,
				node.CreateIdentifier("c"),
			),
		},
	)

	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			functionNode,
			tokens.PTR_TOKEN,
			node.CreateList([]node.Node{
				node.CreateNumber("10"),
				node.CreateNumber("2"),
			}),
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_FunctionCallWithIdentifierAndLeftPointer(t *testing.T) {
	actualAST := getAST("multiply <- (10, 3);")
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			node.CreateIdentifier("multiply"),
			tokens.PTR_TOKEN,
			node.CreateList([]node.Node{
				node.CreateNumber("10"),
				node.CreateNumber("3"),
			}),
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_ReturnStatements(t *testing.T) {
	/*
		Because return statements are now allowed in the global scope (which is enforced by the parser), the return statement
		neeeds to be in a function. So to test that the return AST is correct, it is extracted from the function AST and used
		in the tests.
	*/
	ast := getAST("func() { return 1 + 1; };")
	functionAST := ast[0]
	actualAST := functionAST.GetParam(node.STMTS).Params

	expectedAST := []node.Node{
		node.CreateReturnStatement(node.CreateBinaryExpression(
			node.CreateNumber("1"),
			tokens.PLUS_TOKEN,
			node.CreateNumber("1"),
		)),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_ListIndex(t *testing.T) {
	actualAST := getAST("numbers[1];")
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			node.CreateIdentifier("numbers"),
			tokens.OPEN_BRACKET_TOKEN,
			node.CreateNumber("1"),
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_IfStatement(t *testing.T) {
	actualAST := getAST("if true { print(\"true!!!\"); };")
	expectedAST := []node.Node{
		node.CreateIfStatement(
			node.CreateBooleanTrue(),
			[]node.Node{
				node.CreatePrintStatement([]node.Node{
					node.CreateRawString("true!!!"),
				}),
			},
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
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
