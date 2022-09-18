package tests

import (
	"boomerang/evaluator"
	"boomerang/node"
	"boomerang/tokens"
	"testing"
)

func TestEvaluatorNumbers(t *testing.T) {

	numbers := []string{
		"5",
		"3.1415928",
		"44.357",
	}

	for _, number := range numbers {
		ast := []node.Node{
			{Type: node.NUMBER, Value: number},
		}
		evaluatorObj := evaluator.New(ast)

		actualResults := evaluatorObj.Evaluate()
		expectedResults := []node.Node{
			{Type: node.NUMBER, Value: number},
		}

		AssertNodesEqual(t, expectedResults, actualResults)
	}
}

func TestEvaluatorNegativeNumber(t *testing.T) {
	ast := []node.Node{
		{
			Type: node.UNARY_EXPR,
			Params: []node.Node{
				{Type: tokens.MINUS, Value: "-"},
				{Type: node.NUMBER, Value: "66"},
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
			Params: []node.Node{
				{Type: node.NUMBER, Value: "1"},
				{Type: tokens.PLUS, Value: "+"},
				{Type: node.NUMBER, Value: "1"},
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

func TestEvaluatorPrintStatement(t *testing.T) {
	ast := []node.Node{
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
					Type: node.NUMBER, Value: "3",
				},
			},
		},
	}

	evaluatorObj := evaluator.New(ast)

	actualResults := []node.Node{}
	expectedResults := []node.Node{}

	AssertExpectedOutput(t, "1 2 3\n", func() {
		actualResults = evaluatorObj.Evaluate()
	})

	AssertNodesEqual(t, expectedResults, actualResults)
}

func TestEvaluatorFunction(t *testing.T) {
	ast := []node.Node{
		{
			Type: node.FUNCTION,
			Params: []node.Node{
				{
					Type: node.PARAMETER,
					Params: []node.Node{
						{Type: node.IDENTIFIER, Value: "a"},
						{Type: node.IDENTIFIER, Value: "b"},
					},
				},
				{
					Type: node.STMTS,
					Params: []node.Node{
						{
							Type: node.BIN_EXPR,
							Params: []node.Node{
								{Type: node.IDENTIFIER, Value: "a"},
								{Type: tokens.PLUS, Value: "+"},
								{Type: node.IDENTIFIER, Value: "b"},
							},
						},
					},
				},
			},
		},
	}
	evaluatorObj := evaluator.New(ast)

	actualResults := evaluatorObj.Evaluate()
	AssertNodesEqual(t, ast, actualResults)
}
