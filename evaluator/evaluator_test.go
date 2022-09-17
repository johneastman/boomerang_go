package evaluator

import (
	"boomerang/node"
	"boomerang/testing_utils"
	"boomerang/tokens"
	"testing"
)

func TestNumber(t *testing.T) {
	ast := []node.Node{
		{Type: node.NUMBER, Value: "5"},
	}
	evaluator := New(ast)

	actualResults := evaluator.Evaluate()
	expectedResults := []node.Node{
		{Type: node.NUMBER, Value: "5"},
	}

	testing_utils.AssertNodesEqual(expectedResults, actualResults)
}

func TestNegativeNumber(t *testing.T) {
	ast := []node.Node{
		{
			Type: node.UNARY_EXPR,
			Params: map[string]node.Node{
				node.UNARY_EXPR_EXPR: {Type: node.NUMBER, Value: "66"},
				node.UNARY_EXPR_OP:   {Type: tokens.MINUS, Value: "-"},
			},
		},
	}
	evaluator := New(ast)

	actualResults := evaluator.Evaluate()
	expectedResults := []node.Node{
		{Type: node.NUMBER, Value: "-66"},
	}

	testing_utils.AssertNodesEqual(expectedResults, actualResults)
}

func TestBinaryExpression(t *testing.T) {
	ast := []node.Node{
		{
			Type: node.BIN_EXPR,
			Params: map[string]node.Node{
				node.BIN_EXPR_LEFT:  {Type: node.NUMBER, Value: "1"},
				node.BIN_EXPR_RIGHT: {Type: node.NUMBER, Value: "1"},
				node.BIN_EXPR_OP:    {Type: tokens.PLUS, Value: "+"},
			},
		},
	}
	evaluator := New(ast)

	actualResults := evaluator.Evaluate()
	expectedResults := []node.Node{
		{Type: node.NUMBER, Value: "2"},
	}
	testing_utils.AssertNodesEqual(actualResults, expectedResults)
}
