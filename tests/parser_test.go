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
			{
				Type:  node.NUMBER,
				Value: number,
			},
		}

		AssertNodesEqual(t, expectedAST, actualAST)
	}
}

func TestParser_NegativeNumber(t *testing.T) {
	tokenizer := tokens.New("-66;")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
	expectedAST := []node.Node{
		{
			Type: node.UNARY_EXPR,
			Params: []node.Node{
				{Type: tokens.MINUS, Value: "-"},
				{Type: node.NUMBER, Value: "66"},
			},
		},
	}

	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_BinaryExpression(t *testing.T) {
	tokenizer := tokens.New("7 + 3;")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
	expectedAST := []node.Node{
		{
			Type: node.BIN_EXPR,
			Params: []node.Node{
				{Type: node.NUMBER, Value: "7"},
				{Type: tokens.PLUS, Value: "+"},
				{Type: node.NUMBER, Value: "3"},
			},
		},
	}

	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_Parentheses(t *testing.T) {
	tokenizer := tokens.New("7 + (3);")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
	expectedAST := []node.Node{
		{
			Type: node.BIN_EXPR,
			Params: []node.Node{
				{Type: node.NUMBER, Value: "7"},
				{Type: tokens.PLUS, Value: "+"},
				{Type: node.NUMBER, Value: "3"},
			},
		},
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_ParenthesesBinaryExpression(t *testing.T) {
	tokenizer := tokens.New("7 + (5 - 2);")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
	expectedAST := []node.Node{
		{
			Type: node.BIN_EXPR,
			Params: []node.Node{
				{Type: node.NUMBER, Value: "7"},
				{Type: tokens.PLUS, Value: "+"},
				{Type: node.BIN_EXPR, Params: []node.Node{
					{Type: node.NUMBER, Value: "5"},
					{Type: tokens.MINUS, Value: "-"},
					{Type: node.NUMBER, Value: "2"},
				}},
			},
		},
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_VariableAssignment(t *testing.T) {
	tokenizer := tokens.New("variable = 8 / 2;")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
	expectedAST := []node.Node{
		{
			Type: node.ASSIGN_STMT,
			Params: []node.Node{
				{Type: tokens.IDENTIFIER, Value: "variable"},
				{
					Type: node.BIN_EXPR,
					Params: []node.Node{
						{Type: node.NUMBER, Value: "8"},
						{Type: tokens.FORWARD_SLASH, Value: "/"},
						{Type: node.NUMBER, Value: "2"},
					},
				},
			},
		},
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_Identifier(t *testing.T) {
	tokenizer := tokens.New("variable;")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
	expectedAST := []node.Node{
		{
			Type: node.IDENTIFIER, Value: "variable",
		},
	}

	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_PrintStatement(t *testing.T) {
	tokenizer := tokens.New("print(1, 2, variable);")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
	expectedAST := []node.Node{
		{
			Type: node.PRINT_STMT,
			Params: []node.Node{
				{
					Type: node.NUMBER, Value: "1",
				},
				{
					Type: node.NUMBER, Value: "2",
				},
				{
					Type: node.IDENTIFIER, Value: "variable",
				},
			},
		},
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_PrintStatementNoArguments(t *testing.T) {
	tokenizer := tokens.New("print();")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
	expectedAST := []node.Node{
		{
			Type:   node.PRINT_STMT,
			Params: []node.Node{},
		},
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_Function(t *testing.T) {
	tokenizer := tokens.New("func(a, b) { a + b; };")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
	expectedAST := []node.Node{
		CreateFunction(
			[]string{"a", "b"},
			[]node.Node{
				{
					Type: node.BIN_EXPR,
					Params: []node.Node{
						{Type: node.IDENTIFIER, Value: "a"},
						{Type: tokens.PLUS, Value: "+"},
						{Type: node.IDENTIFIER, Value: "b"},
					},
				},
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
		CreateFunction(
			[]string{},
			[]node.Node{
				{
					Type: node.BIN_EXPR,
					Params: []node.Node{
						{Type: node.NUMBER, Value: "3"},
						{Type: tokens.PLUS, Value: "+"},
						{Type: node.NUMBER, Value: "4"},
					},
				},
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
		CreateFunction(
			[]string{},
			[]node.Node{},
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_FunctionCallWithFunctionLiteral(t *testing.T) {
	tokenizer := tokens.New("func(c, d) { d - c; }(10, 2);")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()

	functionNode := CreateFunction(
		[]string{"c", "d"},
		[]node.Node{
			{
				Type: node.BIN_EXPR,
				Params: []node.Node{
					{Type: node.IDENTIFIER, Value: "d"},
					{Type: tokens.MINUS, Value: "-"},
					{Type: node.IDENTIFIER, Value: "c"},
				},
			},
		},
	)

	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			functionNode,
			tokens.Token{Type: tokens.OPEN_PAREN, Literal: "("},
			node.Node{Type: node.PARAMETER, Params: []node.Node{
				{Type: node.NUMBER, Value: "10"},
				{Type: node.NUMBER, Value: "2"},
			}},
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParser_FunctionCallWithIdentifier(t *testing.T) {
	tokenizer := tokens.New("multiply(10, 3);")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			node.Node{Type: node.IDENTIFIER, Value: "multiply"},
			tokens.Token{Type: tokens.OPEN_PAREN, Literal: "("},
			node.Node{Type: node.PARAMETER, Params: []node.Node{
				{Type: node.NUMBER, Value: "10"},
				{Type: node.NUMBER, Value: "3"},
			}},
		),
	}
	AssertNodesEqual(t, expectedAST, actualAST)
}
