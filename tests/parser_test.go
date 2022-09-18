package tests

import (
	"boomerang/node"
	"boomerang/parser"
	"boomerang/tokens"
	"fmt"
	"testing"
)

func TestParserNumbers(t *testing.T) {
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

func TestParserNegativeNumber(t *testing.T) {
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

func TestParserBinaryExpression(t *testing.T) {
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

func TestParserParentheses(t *testing.T) {
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

func TestParserParenthesesBinaryExpression(t *testing.T) {
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

func TestParserVariableAssignment(t *testing.T) {
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

func TestParserIdentifier(t *testing.T) {
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

func TestPrintStatement(t *testing.T) {
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
