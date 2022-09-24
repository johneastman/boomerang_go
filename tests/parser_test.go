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
		tokenizer := tokens.New(fmt.Sprintf("%s;", number))
		parserObj := parser.New(tokenizer)

		actualAST := parserObj.Parse()
		expectedAST := []node.Node{
			node.CreateNumber(number),
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
		source := fmt.Sprintf("\"%s\";", test.InputSource)
		tokenizer := tokens.New(source)
		parserObj := parser.New(tokenizer)

		actualAST := parserObj.Parse()
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
				node.CreateParameters([]node.Node{
					node.CreateNumber("4"),
					node.CreateNumber("5"),
					node.CreateNumber("6"),
				}),
				node.CreateNumber("3"),
			},
		},
	}

	for _, test := range tests {
		tokenizer := tokens.New(fmt.Sprintf("%s;", test.Source))
		parserObj := parser.New(tokenizer)

		actualAST := parserObj.Parse()
		expectedAST := []node.Node{
			node.CreateParameters(test.Params),
		}

		AssertNodesEqual(t, expectedAST, actualAST)
	}
}

func TestParser_NegativeNumber(t *testing.T) {
	tokenizer := tokens.New("-66;")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
	expectedAST := []node.Node{
		node.CreateUnaryExpression(
			tokens.MINUS_TOKEN,
			node.CreateNumber("66"),
		),
	}

	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_BinaryExpression(t *testing.T) {
	tokenizer := tokens.New("7 + 3;")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
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
	tokenizer := tokens.New("7 + (3);")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
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
	tokenizer := tokens.New("7 + (5 - 2);")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
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
	tokenizer := tokens.New("variable = 8 / 2;")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
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
	tokenizer := tokens.New("variable;")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
	expectedAST := []node.Node{
		node.CreateIdentifier("variable"),
	}

	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_PrintStatement(t *testing.T) {
	tokenizer := tokens.New("print(1, 2, variable);")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
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
	tokenizer := tokens.New("print();")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
	expectedAST := []node.Node{
		node.CreatePrintStatement([]node.Node{}),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_Function(t *testing.T) {
	tokenizer := tokens.New("func(a, b) { a + b; };")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
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
	tokenizer := tokens.New("func() { 3 + 4; };")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
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
	tokenizer := tokens.New("func() {};")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
	expectedAST := []node.Node{
		node.CreateFunction(
			[]node.Node{},
			[]node.Node{},
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_FunctionCallWithNoParameters(t *testing.T) {
	tokenizer := tokens.New("divide <- ();")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			node.CreateIdentifier("divide"),
			tokens.LEFT_PTR_TOKEN,
			node.CreateParameters([]node.Node{}),
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_FunctionCallWithFunctionLiteralAndLeftPointer(t *testing.T) {
	tokenizer := tokens.New("func(c, d) { d - c; } <- (10, 2);")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()

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
			tokens.LEFT_PTR_TOKEN,
			node.CreateParameters([]node.Node{
				node.CreateNumber("10"),
				node.CreateNumber("2"),
			}),
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_FunctionCallWithFunctionLiteralAndRightPointer(t *testing.T) {
	tokenizer := tokens.New("(10, 2) -> func(c, d) { d - c; };")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()

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
			node.CreateParameters([]node.Node{
				node.CreateNumber("10"),
				node.CreateNumber("2"),
			}),
			tokens.RIGHT_PTR_TOKEN,
			functionNode,
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_FunctionCallWithIdentifierAndLeftPointer(t *testing.T) {
	tokenizer := tokens.New("multiply <- (10, 3);")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			node.CreateIdentifier("multiply"),
			tokens.LEFT_PTR_TOKEN,
			node.CreateParameters([]node.Node{
				node.CreateNumber("10"),
				node.CreateNumber("3"),
			}),
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_FunctionCallWithIdentifierAndRightPointer(t *testing.T) {
	tokenizer := tokens.New("(10, 3) -> multiply;")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			node.CreateParameters([]node.Node{
				node.CreateNumber("10"),
				node.CreateNumber("3"),
			}),
			tokens.RIGHT_PTR_TOKEN,
			node.CreateIdentifier("multiply"),
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_ReturnStatements(t *testing.T) {
	tokenizer := tokens.New("return 1 + 1;")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
	expectedAST := []node.Node{
		node.CreateReturnStatement(node.CreateBinaryExpression(
			node.CreateNumber("1"),
			tokens.PLUS_TOKEN,
			node.CreateNumber("1"),
		)),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}
