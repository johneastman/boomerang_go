package tests

import (
	"boomerang/node"
	"boomerang/parser"
	"boomerang/tokens"
	"testing"
)

func TestParserNumber(t *testing.T) {
	tokenizer := tokens.New("10;")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
	expectedAST := []node.Node{
		{
			Type:  node.NUMBER,
			Value: "10",
		},
	}

	AssertNodesEqual(t, expectedAST, actualAST)
}

func TestParserNegativeNumber(t *testing.T) {
	tokenizer := tokens.New("-66;")
	parserObj := parser.New(tokenizer)

	actualAST := parserObj.Parse()
	expectedAST := []node.Node{
		{
			Type: node.UNARY_EXPR,
			Params: map[string]node.Node{
				node.EXPR:     {Type: node.NUMBER, Value: "66"},
				node.OPERATOR: {Type: tokens.MINUS, Value: "-"},
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
			Params: map[string]node.Node{
				node.BIN_EXPR_LEFT:  {Type: node.NUMBER, Value: "7"},
				node.BIN_EXPR_RIGHT: {Type: node.NUMBER, Value: "3"},
				node.OPERATOR:       {Type: tokens.PLUS, Value: "+"},
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
			Params: map[string]node.Node{
				node.BIN_EXPR_LEFT:  {Type: node.NUMBER, Value: "7"},
				node.OPERATOR:       {Type: tokens.PLUS, Value: "+"},
				node.BIN_EXPR_RIGHT: {Type: node.NUMBER, Value: "3"},
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
			Params: map[string]node.Node{
				node.BIN_EXPR_LEFT: {Type: node.NUMBER, Value: "7"},
				node.OPERATOR:      {Type: tokens.PLUS, Value: "+"},
				node.BIN_EXPR_RIGHT: {Type: node.BIN_EXPR, Params: map[string]node.Node{
					node.BIN_EXPR_LEFT:  {Type: node.NUMBER, Value: "5"},
					node.BIN_EXPR_RIGHT: {Type: node.NUMBER, Value: "2"},
					node.OPERATOR:       {Type: tokens.MINUS, Value: "-"},
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
			Params: map[string]node.Node{
				node.ASSIGN_STMT_IDENTIFIER: {Type: tokens.IDENTIFIER, Value: "variable"},
				node.EXPR: {
					Type: node.BIN_EXPR,
					Params: map[string]node.Node{
						node.BIN_EXPR_LEFT:  {Type: node.NUMBER, Value: "8"},
						node.BIN_EXPR_RIGHT: {Type: node.NUMBER, Value: "2"},
						node.OPERATOR:       {Type: tokens.FORWARD_SLASH, Value: "/"},
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
