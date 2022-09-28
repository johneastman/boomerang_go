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
			CreateNumber(number),
		}

		actualResults := getResults(ast)
		expectedResults := []node.Node{
			CreateNumber(number),
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
			CreateBoolean(boolean),
		}

		actualResults := getResults(ast)
		expectedResults := []node.Node{
			CreateBoolean(boolean),
		}

		AssertNodesEqual(t, expectedResults, actualResults)
	}
}

func TestEvaluator_NegativeNumber(t *testing.T) {
	ast := []node.Node{
		node.CreateUnaryExpression(
			CreateTokenFromToken(tokens.MINUS_TOKEN),
			CreateNumber("66"),
		),
	}

	actualResults := getResults(ast)
	expectedResults := []node.Node{
		CreateNumber("-66"),
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
				CreateNumber("12"),
				CreateNumber("45"),
			},
		},
		{
			InputSource:  "the result is <0>",
			OutputSource: "the result is 13",
			Params: []node.Node{
				node.CreateBinaryExpression(
					CreateNumber("7"),
					tokens.PLUS_TOKEN,
					CreateNumber("6"),
				),
			},
		},
		{
			InputSource:  "Hello, my name is <0>, and I am <1> years old!",
			OutputSource: "Hello, my name is John, and I am 5 years old!",
			Params: []node.Node{
				CreateRawString("John"),
				node.CreateBinaryExpression(
					CreateNumber("3"),
					tokens.PLUS_TOKEN,
					CreateNumber("2"),
				),
			},
		},
	}

	for _, test := range tests {
		ast := []node.Node{
			CreateString(test.InputSource, test.Params),
		}

		actualResults := getResults(ast)
		expectedResults := []node.Node{
			CreateRawString(test.OutputSource),
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
			ExpectedResult: CreateList([]node.Node{}),
		},
		{
			Parameters: []node.Node{
				CreateNumber("55"),
			},
			ExpectedResult: CreateList([]node.Node{CreateNumber("55")}),
		},
		{
			Parameters: []node.Node{
				CreateNumber("34"),
				node.CreateBinaryExpression(
					CreateNumber("40"),
					tokens.ASTERISK_TOKEN,
					CreateNumber("3"),
				),
				CreateNumber("66"),
			},
			ExpectedResult: CreateList([]node.Node{
				CreateNumber("34"),
				CreateNumber("120"),
				CreateNumber("66"),
			}),
		},
		{
			Parameters: []node.Node{
				CreateNumber("66"),
				CreateNumber("4"),
				CreateNumber("30"),
			},
			ExpectedResult: CreateList([]node.Node{
				CreateNumber("66"),
				CreateNumber("4"),
				CreateNumber("30"),
			}),
		},
		{
			Parameters: []node.Node{
				CreateNumber("5"),
				CreateList([]node.Node{
					CreateNumber("78"),
				}),
				CreateNumber("60"),
			},
			ExpectedResult: CreateList([]node.Node{
				CreateNumber("5"),
				CreateList([]node.Node{
					CreateNumber("78"),
				}),
				CreateNumber("60"),
			}),
		},
	}

	for _, test := range tests {
		ast := []node.Node{
			CreateList(test.Parameters),
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
			CreateNumber("1"),
			tokens.PLUS_TOKEN,
			CreateNumber("1"),
		),
	}

	actualResults := getResults(ast)
	expectedResults := []node.Node{
		CreateNumber("2"),
	}

	AssertNodesEqual(t, actualResults, expectedResults)
}

func TestEvaluator_Variable(t *testing.T) {
	// Source: variable = 8 / 2; variable;
	ast := []node.Node{
		{
			Type: node.ASSIGN_STMT,
			Params: []node.Node{
				CreateIdentifier("variable"),
				node.CreateBinaryExpression(
					CreateNumber("8"),
					tokens.FORWARD_SLASH_TOKEN,
					CreateNumber("2"),
				),
			},
		},
		CreateIdentifier("variable"),
	}

	actualResults := getResults(ast)
	expectedResults := []node.Node{
		CreateNumber("4"),
	}
	AssertNodesEqual(t, actualResults, expectedResults)
}

func TestEvaluator_PrintStatement(t *testing.T) {
	ast := []node.Node{
		CreatePrintStatement(
			[]node.Node{
				CreateNumber("1"),
				CreateNumber("2"),
				CreateNumber("3"),
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
		CreatePrintStatement([]node.Node{}),
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
		CreateFunction(
			[]node.Node{
				CreateIdentifier("a"),
				CreateIdentifier("b"),
			},
			[]node.Node{
				node.CreateBinaryExpression(
					CreateIdentifier("a"),
					tokens.PLUS_TOKEN,
					CreateIdentifier("b"),
				),
			},
		),
	}

	actualResults := getResults(ast)
	AssertNodesEqual(t, ast, actualResults)
}

func TestEvaluator_FunctionCallWithFunctionLiteral(t *testing.T) {
	functionNode := CreateFunction(
		[]node.Node{
			CreateIdentifier("c"),
			CreateIdentifier("d"),
		},
		[]node.Node{
			node.CreateBinaryExpression(
				CreateIdentifier("c"),
				tokens.MINUS_TOKEN,
				CreateIdentifier("d"),
			),
		},
	)

	ast := []node.Node{
		CreateFunctionCall(functionNode, []node.Node{
			CreateNumber("10"),
			CreateNumber("2"),
		}),
	}

	actualResults := getResults(ast)

	expectedReturnValue := CreateNumber("8")
	expectedResults := []node.Node{
		CreateFunctionReturnValue(&expectedReturnValue),
	}

	AssertNodesEqual(t, expectedResults, actualResults)
}

func TestEvaluator_FunctionCallReturnStatement(t *testing.T) {
	functionNode := CreateFunction(
		[]node.Node{
			CreateIdentifier("c"),
			CreateIdentifier("d"),
		},
		[]node.Node{
			node.CreateBinaryExpression(
				CreateIdentifier("c"),
				tokens.MINUS_TOKEN,
				CreateIdentifier("d"),
			),
			CreateReturnStatement(
				CreateNumber("777"),
			),
			CreateNumber("369"),
		},
	)

	ast := []node.Node{
		CreateFunctionCall(functionNode, []node.Node{
			CreateNumber("10"),
			CreateNumber("2"),
		}),
	}

	actualResults := getResults(ast)

	expectedReturnValue := CreateNumber("777")
	expectedResults := []node.Node{
		CreateFunctionReturnValue(&expectedReturnValue),
	}

	AssertNodesEqual(t, expectedResults, actualResults)
}

func TestEvaluator_TestFunctionCallWithIdentifier(t *testing.T) {

	ast := []node.Node{
		CreateAssignmentStatement(
			"divide",
			CreateFunction(
				[]node.Node{
					CreateIdentifier("a"),
					CreateIdentifier("b"),
				},
				[]node.Node{
					node.CreateBinaryExpression(
						CreateIdentifier("a"),
						tokens.FORWARD_SLASH_TOKEN,
						CreateIdentifier("b"),
					),
				},
			),
		),
		CreateFunctionCall(
			CreateIdentifier("divide"),
			[]node.Node{
				CreateNumber("10"),
				CreateNumber("2"),
			},
		),
	}

	actualResults := getResults(ast)

	expectedReturnValue := CreateNumber("5")
	expectedResults := []node.Node{
		CreateFunctionReturnValue(&expectedReturnValue),
	}

	AssertNodesEqual(t, expectedResults, actualResults)
}

func TestEvaluator_FunctionCallWithNoParameters(t *testing.T) {
	function := CreateFunction(
		[]node.Node{},
		[]node.Node{
			node.CreateBinaryExpression(
				CreateNumber("3"),
				tokens.PLUS_TOKEN,
				CreateNumber("4"),
			),
		},
	)

	ast := []node.Node{
		CreateFunctionCall(function, []node.Node{}),
	}

	actualResults := getResults(ast)

	expectedReturnValue := CreateNumber("7")
	expectedResults := []node.Node{
		CreateFunctionReturnValue(&expectedReturnValue),
	}
	AssertNodesEqual(t, expectedResults, actualResults)
}

func TestEvaluator_FunctionCallNoReturn(t *testing.T) {
	function := CreateFunction(
		[]node.Node{},
		[]node.Node{},
	)

	ast := []node.Node{
		CreateFunctionCall(function, []node.Node{}),
	}

	actualResults := getResults(ast)
	expectedResults := []node.Node{
		CreateFunctionReturnValue(nil),
	}
	AssertNodesEqual(t, expectedResults, actualResults)
}

func TestEvaluator_BuiltinLen(t *testing.T) {
	ast := []node.Node{
		CreateFunctionCall(
			CreateIdentifier("len"),
			[]node.Node{
				CreateNumber("1"),
				CreateNumber("2"),
				CreateNumber("3"),
			},
		),
	}

	actualResults := getResults(ast)
	expectedResults := []node.Node{
		CreateNumber("3"),
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
					CreateNumber("13"),
					tokens.PLUS_TOKEN,
					CreateNumber("7"),
				),
			},
			ExpectedReturnValue: CreateNumber("20"),
		},
		{
			Body:                []node.Node{},
			ExpectedReturnValue: CreateRawString("hello, world!"),
		},
	}

	for _, test := range tests {
		functionName := "function"
		functionAssignment := CreateAssignmentStatement(
			functionName,
			CreateFunction(
				[]node.Node{},
				test.Body,
			),
		)

		resultVariableName := "result"
		functionCallAssignment := CreateAssignmentStatement(
			resultVariableName,
			CreateFunctionCall(
				CreateIdentifier(functionName), []node.Node{},
			),
		)

		unwrapFunctionCall := CreateFunctionCall(
			CreateIdentifier("unwrap"),
			[]node.Node{
				CreateIdentifier(resultVariableName),
				CreateRawString("hello, world!"),
			},
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
		CreateAssignmentStatement(
			"numbers",
			CreateList([]node.Node{
				CreateNumber("1"),
				CreateNumber("2"),
				CreateNumber("3"),
			}),
		),
		node.CreateBinaryExpression(
			CreateIdentifier("numbers"),
			tokens.OPEN_BRACKET_TOKEN,
			CreateNumber("1"),
		),
	}
	actualResults := getResults(ast)
	expectedResults := []node.Node{
		CreateNumber("2"),
	}
	AssertNodesEqual(t, expectedResults, actualResults)
}

func TestEvaluator_IfStatement(t *testing.T) {

	variableName := "variable"

	tests := []struct {
		Condition     node.Node
		ExpectedValue node.Node
	}{
		{
			CreateBooleanTrue(),
			CreateNumber("2"),
		},
		{
			CreateBooleanFalse(),
			CreateNumber("1"),
		},
	}

	for _, test := range tests {
		ast := []node.Node{
			CreateAssignmentStatement(variableName, CreateNumber("1")),
			CreateIfStatement(
				test.Condition,
				[]node.Node{
					CreateAssignmentStatement(variableName, CreateNumber("2")),
				},
			),
			CreateIdentifier(variableName),
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
			Condition: CreateBooleanTrue(),
			ReturnValue: CreateList([]node.Node{
				CreateBooleanTrue(),
				CreateNumber("7"),
			}),
		},
		{
			Condition: CreateBooleanFalse(),
			ReturnValue: CreateList([]node.Node{
				CreateBooleanTrue(),
				CreateNumber("0"),
			}),
		},
	}

	for _, test := range tests {
		ast := []node.Node{
			CreateFunctionCall(
				CreateFunction(
					[]node.Node{
						CreateIdentifier("a"),
						CreateIdentifier("b"),
					},
					[]node.Node{
						CreateIfStatement(
							test.Condition,
							[]node.Node{
								CreateReturnStatement(
									node.CreateBinaryExpression(
										CreateIdentifier("a"),
										tokens.PLUS_TOKEN,
										CreateIdentifier("b"),
									),
								),
							},
						),
						CreateNumber("0"),
					},
				),
				[]node.Node{
					CreateNumber("3"),
					CreateNumber("4"),
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

func TestEvaluator_FunctionCallPrecedenceExpression(t *testing.T) {
	/*
		Source:
			add = func(a, b) {
				return a + b;
			};

			sum = add <- (3, 4);
			value = unwrap <- (sum, 0) + 7;
			value;
	*/
	addFunction := CreateAssignmentStatement(
		"add",
		CreateFunction(
			[]node.Node{
				CreateIdentifier("a"),
				CreateIdentifier("b"),
			},
			[]node.Node{
				CreateReturnStatement(
					node.CreateBinaryExpression(
						CreateIdentifier("a"),
						CreateTokenFromToken(tokens.PLUS_TOKEN),
						CreateIdentifier("b"),
					),
				),
			},
		),
	)

	addFunctionReturnValue := CreateAssignmentStatement(
		"sum",
		CreateFunctionCall(
			CreateIdentifier("add"),
			[]node.Node{
				CreateNumber("3"),
				CreateNumber("4"),
			},
		),
	)

	actualValue := CreateAssignmentStatement(
		"value",
		node.CreateBinaryExpression(
			CreateFunctionCall(
				CreateIdentifier("unwrap"),
				[]node.Node{
					CreateIdentifier("sum"),
					CreateNumber("0"),
				},
			),
			CreateTokenFromToken(tokens.PLUS_TOKEN),
			CreateNumber("3"),
		),
	)

	ast := []node.Node{
		addFunction,
		addFunctionReturnValue,
		actualValue,
		CreateIdentifier("value"),
	}

	actualResults := getResults(ast)
	expectedResults := []node.Node{
		CreateNumber("10"),
	}
	AssertNodesEqual(t, expectedResults, actualResults)
}

func getResults(ast []node.Node) []node.Node {
	evaluatorObj := evaluator.New(ast)
	actualResults, err := evaluatorObj.Evaluate()
	if err != nil {
		panic(err.Error())
	}
	return *actualResults
}
