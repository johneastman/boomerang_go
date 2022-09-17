package tests

import (
	"boomerang/evaluator"
	"boomerang/node"
	"boomerang/tokens"
	"testing"
)

func TestEvaluatorNumber(t *testing.T) {
	ast := []node.Node{
		{Type: node.NUMBER, Value: "5"},
	}
	evaluatorObj := evaluator.New(ast)

	actualResults := evaluatorObj.Evaluate()
	expectedResults := []node.Node{
		{Type: node.NUMBER, Value: "5"},
	}

	AssertNodesEqual(t, expectedResults, actualResults)
}

func TestEvaluatorNegativeNumber(t *testing.T) {
	ast := []node.Node{
		{
			Type: node.UNARY_EXPR,
			Params: map[string]node.Node{
				node.EXPR:     {Type: node.NUMBER, Value: "66"},
				node.OPERATOR: {Type: tokens.MINUS, Value: "-"},
			},
		},
	}
	evaluatorObj := evaluator.New(ast)

	actualResults := evaluatorObj.Evaluate()
	expectedResults := []node.Node{
		{Type: node.NUMBER, Value: "-66"},
	}

	AssertNodesEqual(t, expectedResults, actualResults)
}

func TestEvaluatorBinaryExpression(t *testing.T) {
	ast := []node.Node{
		{
			Type: node.BIN_EXPR,
			Params: map[string]node.Node{
				node.BIN_EXPR_LEFT:  {Type: node.NUMBER, Value: "1"},
				node.BIN_EXPR_RIGHT: {Type: node.NUMBER, Value: "1"},
				node.OPERATOR:       {Type: tokens.PLUS, Value: "+"},
			},
		},
	}
	evaluatorObj := evaluator.New(ast)

	actualResults := evaluatorObj.Evaluate()
	expectedResults := []node.Node{
		{Type: node.NUMBER, Value: "2"},
	}
	AssertNodesEqual(t, actualResults, expectedResults)
}

func TestVariable(t *testing.T) {
	// Source: variable = 8 / 2; variable;
	ast := []node.Node{
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
		{
			Type: node.IDENTIFIER, Value: "variable",
		},
	}

	evaluatorObj := evaluator.New(ast)

	actualResults := evaluatorObj.Evaluate()
	expectedResults := []node.Node{
		{
			Type: node.NUMBER, Value: "4",
		},
	}
	AssertNodesEqual(t, actualResults, expectedResults)
}
