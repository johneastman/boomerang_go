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

		actualResults := getResults(ast)
		expectedResults := []node.Node{
			node.CreateNumber(number),
		}

		AssertNodesEqual(t, expectedResults, actualResults)
	}
}

func TestEvaluator_Booleans(t *testing.T) {

	booleans := []string{
		"true",
		"false",
	}

	for _, boolean := range booleans {
		ast := []node.Node{
			node.CreateBoolean(boolean),
		}

		actualResults := getResults(ast)
		expectedResults := []node.Node{
			node.CreateBoolean(boolean),
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

	actualResults := getResults(ast)
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
				node.CreateRawString("John"),
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

		actualResults := getResults(ast)
		expectedResults := []node.Node{
			node.CreateRawString(test.OutputSource),
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
			ExpectedResult: node.CreateList([]node.Node{}),
		},
		{
			Parameters: []node.Node{
				node.CreateNumber("55"),
			},
			ExpectedResult: node.CreateList([]node.Node{node.CreateNumber("55")}),
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
			ExpectedResult: node.CreateList([]node.Node{
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
			ExpectedResult: node.CreateList([]node.Node{
				node.CreateNumber("66"),
				node.CreateNumber("4"),
				node.CreateNumber("30"),
			}),
		},
		{
			Parameters: []node.Node{
				node.CreateNumber("5"),
				node.CreateList([]node.Node{
					node.CreateNumber("78"),
				}),
				node.CreateNumber("60"),
			},
			ExpectedResult: node.CreateList([]node.Node{
				node.CreateNumber("5"),
				node.CreateList([]node.Node{
					node.CreateNumber("78"),
				}),
				node.CreateNumber("60"),
			}),
		},
	}

	for _, test := range tests {
		ast := []node.Node{
			node.CreateList(test.Parameters),
		}

		actualResults := getResults(ast)
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

	actualResults := getResults(ast)
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

	actualResults := getResults(ast)
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

	actualResults := []node.Node{}
	expectedResults := []node.Node{}

	AssertExpectedOutput(t, "1 2 3\n", func() {
		actualResults = getResults(ast)
	})

	AssertNodesEqual(t, expectedResults, actualResults)
}

func TestEvaluator_PrintStatementNoArguments(t *testing.T) {
	ast := []node.Node{
		node.CreatePrintStatement([]node.Node{}),
	}

	actualResults := []node.Node{}
	expectedResults := []node.Node{}

	AssertExpectedOutput(t, "", func() {
		actualResults = getResults(ast)
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

	actualResults := getResults(ast)
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

	actualResults := getResults(ast)

	expectedReturnValue := node.CreateNumber("8")
	expectedResults := []node.Node{
		node.CreateReturnValue(&expectedReturnValue),
	}

	AssertNodesEqual(t, expectedResults, actualResults)
}

func TestEvaluator_FunctionCallReturnStatement(t *testing.T) {
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
			node.CreateReturnStatement(
				node.CreateNumber("777"),
			),
			node.CreateNumber("369"),
		},
	)

	ast := []node.Node{
		node.CreateFunctionCall(functionNode, []node.Node{
			node.CreateNumber("10"),
			node.CreateNumber("2"),
		}),
	}

	actualResults := getResults(ast)

	expectedReturnValue := node.CreateNumber("777")
	expectedResults := []node.Node{
		node.CreateReturnValue(&expectedReturnValue),
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

	actualResults := getResults(ast)

	expectedReturnValue := node.CreateNumber("5")
	expectedResults := []node.Node{
		node.CreateReturnValue(&expectedReturnValue),
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

	actualResults := getResults(ast)

	expectedReturnValue := node.CreateNumber("7")
	expectedResults := []node.Node{
		node.CreateReturnValue(&expectedReturnValue),
	}
	AssertNodesEqual(t, expectedResults, actualResults)
}

func TestEvaluator_FunctionCallNoReturn(t *testing.T) {
	function := node.CreateFunction(
		[]node.Node{},
		[]node.Node{},
	)

	ast := []node.Node{
		node.CreateFunctionCall(function, []node.Node{}),
	}

	actualResults := getResults(ast)
	expectedResults := []node.Node{
		node.CreateReturnValue(nil),
	}
	AssertNodesEqual(t, expectedResults, actualResults)
}

func TestEvaluator_BuiltinLen(t *testing.T) {
	ast := []node.Node{
		node.CreateBinaryExpression(
			node.CreateIdentifier("len"),
			tokens.PTR_TOKEN,
			node.CreateList([]node.Node{
				node.CreateNumber("1"),
				node.CreateNumber("2"),
				node.CreateNumber("3"),
			}),
		),
	}

	actualResults := getResults(ast)
	expectedResults := []node.Node{
		node.CreateNumber("3"),
	}
	AssertNodesEqual(t, expectedResults, actualResults)
}

func TestEvaluator_BuiltinUnwrapReturnValue(t *testing.T) {

	tests := []struct {
		Body                []node.Node
		ExpectedReturnValue node.Node
	}{
		{
			Body: []node.Node{
				node.CreateBinaryExpression(
					node.CreateNumber("13"),
					tokens.PLUS_TOKEN,
					node.CreateNumber("7"),
				),
			},
			ExpectedReturnValue: node.CreateNumber("20"),
		},
		{
			Body:                []node.Node{},
			ExpectedReturnValue: node.CreateRawString("hello, world!"),
		},
	}

	for _, test := range tests {
		functionName := "function"
		functionAssignment := node.CreateAssignmentStatement(
			functionName,
			node.CreateFunction(
				[]node.Node{},
				test.Body,
			),
		)

		resultVariableName := "result"
		functionCallAssignment := node.CreateAssignmentStatement(
			resultVariableName,
			node.CreateFunctionCall(
				node.CreateIdentifier(functionName), []node.Node{},
			),
		)

		unwrapFunctionCall := node.CreateBinaryExpression(
			node.CreateIdentifier("unwrap"),
			tokens.PTR_TOKEN,
			node.CreateList(
				[]node.Node{
					node.CreateIdentifier(resultVariableName),
					node.CreateRawString("hello, world!"),
				},
			),
		)

		ast := []node.Node{
			functionAssignment,
			functionCallAssignment,
			unwrapFunctionCall,
		}

		actualResults := getResults(ast)
		expectedResults := []node.Node{
			test.ExpectedReturnValue,
		}
		AssertNodesEqual(t, expectedResults, actualResults)
	}
}

func TestEvaluator_ListIndex(t *testing.T) {
	ast := []node.Node{
		node.CreateAssignmentStatement(
			"numbers",
			node.CreateList([]node.Node{
				node.CreateNumber("1"),
				node.CreateNumber("2"),
				node.CreateNumber("3"),
			}),
		),
		node.CreateBinaryExpression(
			node.CreateIdentifier("numbers"),
			tokens.OPEN_BRACKET_TOKEN,
			node.CreateNumber("1"),
		),
	}
	actualResults := getResults(ast)
	expectedResults := []node.Node{
		node.CreateNumber("2"),
	}
	AssertNodesEqual(t, expectedResults, actualResults)
}

func getResults(ast []node.Node) []node.Node {
	evaluatorObj := evaluator.New(ast)
	actualResults, _ := evaluatorObj.Evaluate()
	return *actualResults
}

func TestEvaluator_IfStatement(t *testing.T) {

	variableName := "variable"

	tests := []struct {
		Condition     node.Node
		ExpectedValue node.Node
	}{
		{
			node.CreateBooleanTrue(),
			node.CreateNumber("2"),
		},
		{
			node.CreateBooleanFalse(),
			node.CreateNumber("1"),
		},
	}

	for _, test := range tests {
		ast := []node.Node{
			node.CreateAssignmentStatement(variableName, node.CreateNumber("1")),
			node.CreateIfStatement(
				test.Condition,
				[]node.Node{
					node.CreateAssignmentStatement(variableName, node.CreateNumber("2")),
				},
			),
			node.CreateIdentifier(variableName),
		}
		actualResults := getResults(ast)
		expectedResults := []node.Node{
			test.ExpectedValue,
		}
		AssertNodesEqual(t, expectedResults, actualResults)
	}
}

func TestEvaluator_FunctionReturnIfStatement(t *testing.T) {
	/*
		Source:
		func(a, b) {
			if true {
				return a + b;
			}
			return 0;
		}
	*/
	tests := []struct {
		Condition   node.Node
		ReturnValue node.Node
	}{
		{
			Condition: node.CreateBooleanTrue(),
			ReturnValue: node.CreateList([]node.Node{
				node.CreateBooleanTrue(),
				node.CreateNumber("7"),
			}),
		},
		{
			Condition: node.CreateBooleanFalse(),
			ReturnValue: node.CreateList([]node.Node{
				node.CreateBooleanTrue(),
				node.CreateNumber("0"),
			}),
		},
	}

	for _, test := range tests {
		ast := []node.Node{
			node.CreateFunctionCall(
				node.CreateFunction(
					[]node.Node{
						node.CreateIdentifier("a"),
						node.CreateIdentifier("b"),
					},
					[]node.Node{
						node.CreateIfStatement(
							test.Condition,
							[]node.Node{
								node.CreateReturnStatement(
									node.CreateBinaryExpression(
										node.CreateIdentifier("a"),
										tokens.PLUS_TOKEN,
										node.CreateIdentifier("b"),
									),
								),
							},
						),
						node.CreateNumber("0"),
					},
				),
				[]node.Node{
					node.CreateNumber("3"),
					node.CreateNumber("4"),
				},
			),
		}
		actualResults := getResults(ast)
		expectedResults := []node.Node{
			test.ReturnValue,
		}
		AssertNodesEqual(t, expectedResults, actualResults)
	}
}
