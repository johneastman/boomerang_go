package tests

import (
	"boomerang/evaluator"
	"boomerang/node"
	"boomerang/tokens"
	"testing"
)

func TestEvaluator_Numbers(t *testing.T) {

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

func TestEvaluator_NegativeNumber(t *testing.T) {
	ast := []node.Node{
		{
			Type: node.UNARY_EXPR,
			Params: []node.Node{
				node.CreateTokenNode(tokens.MINUS_TOKEN),
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

func TestEvaluator_BinaryExpression(t *testing.T) {
	ast := []node.Node{
		{
			Type: node.BIN_EXPR,
			Params: []node.Node{
				{Type: node.NUMBER, Value: "1"},
				node.CreateTokenNode(tokens.PLUS_TOKEN),
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

func TestEvaluator_Variable(t *testing.T) {
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
						node.CreateTokenNode(tokens.FORWARD_SLASH_TOKEN),
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

func TestEvaluator_PrintStatement(t *testing.T) {
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

func TestEvaluator_PrintStatementNoArguments(t *testing.T) {
	ast := []node.Node{
		{
			Type:   node.PRINT_STMT,
			Params: []node.Node{},
		},
	}

	evaluatorObj := evaluator.New(ast)

	actualResults := []node.Node{}
	expectedResults := []node.Node{}

	AssertExpectedOutput(t, "", func() {
		actualResults = evaluatorObj.Evaluate()
	})

	AssertNodesEqual(t, expectedResults, actualResults)
}

func TestEvaluator_Function(t *testing.T) {
	ast := []node.Node{
		CreateFunction(
			[]string{"a", "b"},
			[]node.Node{
				{
					Type: node.BIN_EXPR,
					Params: []node.Node{
						{Type: node.IDENTIFIER, Value: "a"},
						node.CreateTokenNode(tokens.PLUS_TOKEN),
						{Type: node.IDENTIFIER, Value: "b"},
					},
				},
			},
		),
	}
	evaluatorObj := evaluator.New(ast)

	actualResults := evaluatorObj.Evaluate()
	AssertNodesEqual(t, ast, actualResults)
}

func TestEvaluator_FunctionCallWithFunctionLiteral(t *testing.T) {
	functionNode := CreateFunction(
		[]string{"c", "d"},
		[]node.Node{
			{
				Type: node.BIN_EXPR,
				Params: []node.Node{
					{Type: node.IDENTIFIER, Value: "c"},
					node.CreateTokenNode(tokens.MINUS_TOKEN),
					{Type: node.IDENTIFIER, Value: "d"},
				},
			},
		},
	)

	ast := []node.Node{
		CreateFunctionCall(functionNode, []node.Node{
			{Type: node.NUMBER, Value: "10"},
			{Type: node.NUMBER, Value: "2"},
		}),
	}

	evaluatorObj := evaluator.New(ast)

	actualResults := evaluatorObj.Evaluate()
	expectedResults := []node.Node{
		{Type: node.NUMBER, Value: "8"},
	}
	AssertNodesEqual(t, expectedResults, actualResults)
}

func TestEvaluator_TestFunctionCallWithIdentifier(t *testing.T) {

	ast := []node.Node{
		node.CreateAssignmentStatement(
			"divide",
			CreateFunction(
				[]string{"a", "b"},
				[]node.Node{
					node.CreateBinaryExpression(
						node.Node{Type: node.IDENTIFIER, Value: "a"},
						tokens.FORWARD_SLASH_TOKEN,
						node.Node{Type: node.IDENTIFIER, Value: "b"},
					),
				},
			),
		),
		CreateFunctionCall(
			node.CreateIdentifier("divide"),
			[]node.Node{
				{Type: node.NUMBER, Value: "10"},
				{Type: node.NUMBER, Value: "2"},
			},
		),
	}

	evaluatorObj := evaluator.New(ast)

	actualResults := evaluatorObj.Evaluate()
	expectedResults := []node.Node{
		{Type: node.NUMBER, Value: "5"},
	}
	AssertNodesEqual(t, expectedResults, actualResults)
}

func TestEvaluator_FunctionCallWithNoParameters(t *testing.T) {
	function := CreateFunction(
		[]string{},
		[]node.Node{
			{
				Type: node.BIN_EXPR,
				Params: []node.Node{
					{Type: node.NUMBER, Value: "3"},
					node.CreateTokenNode(tokens.PLUS_TOKEN),
					{Type: node.NUMBER, Value: "4"},
				},
			},
		},
	)

	ast := []node.Node{
		CreateFunctionCall(function, []node.Node{}),
	}

	evaluatorObj := evaluator.New(ast)

	actualResults := evaluatorObj.Evaluate()
	expectedResults := []node.Node{
		{Type: node.NUMBER, Value: "7"},
	}
	AssertNodesEqual(t, expectedResults, actualResults)
}
