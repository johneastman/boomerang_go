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
			node.CreateNumber(number),
		}
		evaluatorObj := evaluator.New(ast)

		actualResults := evaluatorObj.Evaluate()
		expectedResults := []node.Node{
			node.CreateNumber(number),
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
				node.CreateNumber("66"),
			},
		},
	}
	evaluatorObj := evaluator.New(ast)

	actualResults := evaluatorObj.Evaluate()
	expectedResults := []node.Node{
		node.CreateNumber("-66"),
	}

	AssertNodesEqual(t, expectedResults, actualResults)
}

func TestEvaluator_Strings(t *testing.T) {

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
			InputSource:  "the time is <0>:<1>",
			OutputSource: "the time is 12:45",
			Params: []node.Node{
				node.CreateNumber("12"),
				node.CreateNumber("45"),
			},
		},
		{
			InputSource:  "the result is <0>",
			OutputSource: "the result is 13",
			Params: []node.Node{
				node.CreateBinaryExpression(
					node.CreateNumber("7"),
					tokens.PLUS_TOKEN,
					node.CreateNumber("6"),
				),
			},
		},
		{
			InputSource:  "Hello, my name is <0>, and I am <1> years old!",
			OutputSource: "Hello, my name is John, and I am 5 years old!",
			Params: []node.Node{
				node.CreateString("John", []node.Node{}),
				node.CreateBinaryExpression(
					node.CreateNumber("3"),
					tokens.PLUS_TOKEN,
					node.CreateNumber("2"),
				),
			},
		},
	}

	for _, test := range tests {
		ast := []node.Node{
			node.CreateString(test.InputSource, test.Params),
		}
		evaluatorObj := evaluator.New(ast)

		actualResults := evaluatorObj.Evaluate()
		expectedResults := []node.Node{
			node.CreateString(test.OutputSource, []node.Node{}),
		}
		AssertNodesEqual(t, expectedResults, actualResults)
	}
}

func TestEvaluator_Parameters(t *testing.T) {

	tests := []struct {
		Parameters     []node.Node
		ExpectedResult node.Node
	}{
		{
			Parameters:     []node.Node{},
			ExpectedResult: node.CreateParameters([]node.Node{}),
		},
		{
			Parameters: []node.Node{
				node.CreateNumber("55"),
			},
			ExpectedResult: node.CreateParameters([]node.Node{node.CreateNumber("55")}),
		},
		{
			Parameters: []node.Node{
				node.CreateNumber("34"),
				node.CreateBinaryExpression(
					node.CreateNumber("40"),
					tokens.ASTERISK_TOKEN,
					node.CreateNumber("3"),
				),
				node.CreateNumber("66"),
			},
			ExpectedResult: node.CreateParameters([]node.Node{
				node.CreateNumber("34"),
				node.CreateNumber("120"),
				node.CreateNumber("66"),
			}),
		},
		{
			Parameters: []node.Node{
				node.CreateNumber("66"),
				node.CreateNumber("4"),
				node.CreateNumber("30"),
			},
			ExpectedResult: node.CreateParameters([]node.Node{
				node.CreateNumber("66"),
				node.CreateNumber("4"),
				node.CreateNumber("30"),
			}),
		},
		{
			Parameters: []node.Node{
				node.CreateNumber("5"),
				node.CreateParameters([]node.Node{
					node.CreateNumber("78"),
				}),
				node.CreateNumber("60"),
			},
			ExpectedResult: node.CreateParameters([]node.Node{
				node.CreateNumber("5"),
				node.CreateParameters([]node.Node{
					node.CreateNumber("78"),
				}),
				node.CreateNumber("60"),
			}),
		},
	}

	for _, test := range tests {
		ast := []node.Node{
			node.CreateParameters(test.Parameters),
		}
		evaluatorObj := evaluator.New(ast)

		actualResults := evaluatorObj.Evaluate()
		expectedResults := []node.Node{
			test.ExpectedResult,
		}
		AssertNodesEqual(t, expectedResults, actualResults)
	}
}

func TestEvaluator_BinaryExpression(t *testing.T) {
	ast := []node.Node{
		node.CreateBinaryExpression(
			node.CreateNumber("1"),
			tokens.PLUS_TOKEN,
			node.CreateNumber("1"),
		),
	}
	evaluatorObj := evaluator.New(ast)

	actualResults := evaluatorObj.Evaluate()
	expectedResults := []node.Node{
		node.CreateNumber("2"),
	}

	AssertNodesEqual(t, actualResults, expectedResults)
}

func TestEvaluator_Variable(t *testing.T) {
	// Source: variable = 8 / 2; variable;
	ast := []node.Node{
		{
			Type: node.ASSIGN_STMT,
			Params: []node.Node{
				node.CreateIdentifier("variable"),
				node.CreateBinaryExpression(
					node.CreateNumber("8"),
					tokens.FORWARD_SLASH_TOKEN,
					node.CreateNumber("2"),
				),
			},
		},
		node.CreateIdentifier("variable"),
	}

	evaluatorObj := evaluator.New(ast)

	actualResults := evaluatorObj.Evaluate()
	expectedResults := []node.Node{
		node.CreateNumber("4"),
	}
	AssertNodesEqual(t, actualResults, expectedResults)
}

func TestEvaluator_PrintStatement(t *testing.T) {
	ast := []node.Node{
		node.CreatePrintStatement(
			[]node.Node{
				node.CreateNumber("1"),
				node.CreateNumber("2"),
				node.CreateNumber("3"),
			},
		),
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
		node.CreatePrintStatement([]node.Node{}),
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
	evaluatorObj := evaluator.New(ast)

	actualResults := evaluatorObj.Evaluate()
	AssertNodesEqual(t, ast, actualResults)
}

func TestEvaluator_FunctionCallWithFunctionLiteral(t *testing.T) {
	functionNode := node.CreateFunction(
		[]node.Node{
			node.CreateIdentifier("c"),
			node.CreateIdentifier("d"),
		},
		[]node.Node{
			node.CreateBinaryExpression(
				node.CreateIdentifier("c"),
				tokens.MINUS_TOKEN,
				node.CreateIdentifier("d"),
			),
		},
	)

	ast := []node.Node{
		node.CreateFunctionCall(functionNode, []node.Node{
			node.CreateNumber("10"),
			node.CreateNumber("2"),
		}),
	}

	evaluatorObj := evaluator.New(ast)

	actualResults := evaluatorObj.Evaluate()
	expectedResults := []node.Node{
		node.CreateNumber("8"),
	}

	AssertNodesEqual(t, expectedResults, actualResults)
}

func TestEvaluator_TestFunctionCallWithIdentifier(t *testing.T) {

	ast := []node.Node{
		node.CreateAssignmentStatement(
			"divide",
			node.CreateFunction(
				[]node.Node{
					node.CreateIdentifier("a"),
					node.CreateIdentifier("b"),
				},
				[]node.Node{
					node.CreateBinaryExpression(
						node.CreateIdentifier("a"),
						tokens.FORWARD_SLASH_TOKEN,
						node.CreateIdentifier("b"),
					),
				},
			),
		),
		node.CreateFunctionCall(
			node.CreateIdentifier("divide"),
			[]node.Node{
				node.CreateNumber("10"),
				node.CreateNumber("2"),
			},
		),
	}

	evaluatorObj := evaluator.New(ast)

	actualResults := evaluatorObj.Evaluate()
	expectedResults := []node.Node{
		node.CreateNumber("5"),
	}

	AssertNodesEqual(t, expectedResults, actualResults)
}

func TestEvaluator_FunctionCallWithNoParameters(t *testing.T) {
	function := node.CreateFunction(
		[]node.Node{},
		[]node.Node{
			node.CreateBinaryExpression(
				node.CreateNumber("3"),
				tokens.PLUS_TOKEN,
				node.CreateNumber("4"),
			),
		},
	)

	ast := []node.Node{
		node.CreateFunctionCall(function, []node.Node{}),
	}

	evaluatorObj := evaluator.New(ast)

	actualResults := evaluatorObj.Evaluate()
	expectedResults := []node.Node{
		node.CreateNumber("7"),
	}
	AssertNodesEqual(t, expectedResults, actualResults)
}
