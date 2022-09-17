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
			Type:  "Number",
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
			Type: "BinaryExpression",
			Params: map[string]node.Node{
				"left":     {Type: "Number", Value: "7"},
				"right":    {Type: "Number", Value: "3"},
				"operator": {Type: "PLUS", Value: "+"},
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
			Type: "BinaryExpression",
			Params: map[string]node.Node{
				"left":     {Type: "Number", Value: "7"},
				"operator": {Type: "PLUS", Value: "+"},
				"right":    {Type: "Number", Value: "3"},
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
			Type: "BinaryExpression",
			Params: map[string]node.Node{
				"left":     {Type: "Number", Value: "7"},
				"operator": {Type: "PLUS", Value: "+"},
				"right": {Type: "BinaryExpression", Params: map[string]node.Node{
					"left":     {Type: "Number", Value: "5"},
					"right":    {Type: "Number", Value: "3"},
					"operator": {Type: "MINUS", Value: "-"},
				}},
			},
		},
	}
	testing_utils.AssertNodesEqual(expectedAST, actualAST)
}
