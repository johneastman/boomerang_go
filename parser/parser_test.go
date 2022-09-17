package parser

import (
	"boomerang/node"
	"boomerang/testing_utils"
	"boomerang/tokens"
	"testing"
)

func TestNumber(t *testing.T) {
	tokenizer := tokens.New("10;")
	parser := New(tokenizer)

	actualAST := parser.Parse()
	expectedAST := []node.Node{
		{
			Type:  node.NUMBER,
			Value: "10",
		},
	}

	testing_utils.AssertNodesEqual(expectedAST, actualAST)
}

func TestBinaryExpression(t *testing.T) {
	tokenizer := tokens.New("7 + 3;")
	parser := New(tokenizer)

	actualAST := parser.Parse()
	expectedAST := []node.Node{
		{
			Type: node.BIN_EXPR,
			Params: map[string]node.Node{
				node.BIN_EXPR_LEFT:  {Type: node.NUMBER, Value: "7"},
				node.BIN_EXPR_RIGHT: {Type: node.NUMBER, Value: "3"},
				node.BIN_EXPR_OP:    {Type: tokens.PLUS, Value: "+"},
			},
		},
	}

	testing_utils.AssertNodesEqual(expectedAST, actualAST)
}

func TestParentheses(t *testing.T) {
	tokenizer := tokens.New("7 + (3);")
	parser := New(tokenizer)

	actualAST := parser.Parse()
	expectedAST := []node.Node{
		{
			Type: node.BIN_EXPR,
			Params: map[string]node.Node{
				node.BIN_EXPR_LEFT:  {Type: node.NUMBER, Value: "7"},
				node.BIN_EXPR_OP:    {Type: tokens.PLUS, Value: "+"},
				node.BIN_EXPR_RIGHT: {Type: node.NUMBER, Value: "3"},
			},
		},
	}
	testing_utils.AssertNodesEqual(expectedAST, actualAST)
}

func TestParenthesesBinaryExpression(t *testing.T) {
	tokenizer := tokens.New("7 + (5 - 2);")
	parser := New(tokenizer)

	actualAST := parser.Parse()
	expectedAST := []node.Node{
		{
			Type: node.BIN_EXPR,
			Params: map[string]node.Node{
				node.BIN_EXPR_LEFT: {Type: node.NUMBER, Value: "7"},
				node.BIN_EXPR_OP:   {Type: tokens.PLUS, Value: "+"},
				node.BIN_EXPR_RIGHT: {Type: node.BIN_EXPR, Params: map[string]node.Node{
					node.BIN_EXPR_LEFT:  {Type: node.NUMBER, Value: "5"},
					node.BIN_EXPR_RIGHT: {Type: node.NUMBER, Value: "3"},
					node.BIN_EXPR_OP:    {Type: tokens.MINUS, Value: "-"},
				}},
			},
		},
	}
	testing_utils.AssertNodesEqual(expectedAST, actualAST)
}
