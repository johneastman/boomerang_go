package tests

import (
	"boomerang/evaluator"
	"boomerang/node"
	"boomerang/tokens"
	"fmt"
	"testing"
)

func TestEvaluator_Numbers(t *testing.T) {

	numbers := []string{
		"5",
		"3.1415928",
		"44.357",
	}

	for i, number := range numbers {
		ast := []node.Node{
			CreateNumber(number),
		}

		actualResults := getResults(ast)
		expectedResults := []node.Node{
			CreateNumber(number),
		}

		AssertNodesEqual(t, i, expectedResults, actualResults)
	}
}

func TestEvaluator_Booleans(t *testing.T) {

	booleans := []string{
		"true",
		"false",
	}

	for i, boolean := range booleans {
		ast := []node.Node{
			CreateBoolean(boolean),
		}

		actualResults := getResults(ast)
		expectedResults := []node.Node{
			CreateBoolean(boolean),
		}

		AssertNodesEqual(t, i, expectedResults, actualResults)
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

	AssertNodesEqual(t, 0, expectedResults, actualResults)
}

func TestEvaluator_Bang(t *testing.T) {

	tests := []struct {
		Input          node.Node
		ExpectedResult node.Node
	}{
		{
			Input:          CreateBooleanTrue(),
			ExpectedResult: CreateBooleanFalse(),
		},
		{
			Input:          CreateBooleanFalse(),
			ExpectedResult: CreateBooleanTrue(),
		},
	}

	for i, test := range tests {
		ast := []node.Node{
			node.CreateUnaryExpression(
				CreateTokenFromToken(tokens.NOT_TOKEN),
				test.Input,
			),
		}

		actualResults := getResults(ast)
		expectedResults := []node.Node{
			test.ExpectedResult,
		}

		AssertNodesEqual(t, i, expectedResults, actualResults)
	}
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
		{
			InputSource:  "My numbers are <0>!",
			OutputSource: "My numbers are (1, 2, 3, 4)!",
			Params: []node.Node{
				CreateList([]node.Node{
					CreateNumber("1"),
					CreateNumber("2"),
					CreateNumber("3"),
					CreateNumber("4"),
				}),
			},
		},
	}

	for i, test := range tests {
		ast := []node.Node{
			CreateString(test.InputSource, test.Params),
		}

		actualResults := getResults(ast)
		expectedResults := []node.Node{
			CreateRawString(test.OutputSource),
		}
		AssertNodesEqual(t, i, expectedResults, actualResults)
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

	for i, test := range tests {
		ast := []node.Node{
			CreateList(test.Parameters),
		}

		actualResults := getResults(ast)
		expectedResults := []node.Node{
			test.ExpectedResult,
		}
		AssertNodesEqual(t, i, expectedResults, actualResults)
	}
}

func TestEvaluator_BinaryExpressions(t *testing.T) {

	tests := []struct {
		AST    node.Node
		Result node.Node
	}{
		{
			AST: node.CreateBinaryExpression(
				CreateNumber("1"),
				CreateTokenFromToken(tokens.PLUS_TOKEN),
				CreateNumber("1"),
			),
			Result: CreateNumber("2"),
		},
		{
			AST: node.CreateBinaryExpression(
				CreateList([]node.Node{
					CreateBooleanTrue(),
					CreateBooleanFalse(),
				}),
				CreateTokenFromToken(tokens.PTR_TOKEN),
				CreateBooleanTrue(),
			),
			Result: CreateList([]node.Node{
				CreateBooleanTrue(),
				CreateBooleanFalse(),
				CreateBooleanTrue(),
			}),
		},
		{
			AST: node.CreateBinaryExpression(
				CreateList([]node.Node{
					CreateBooleanTrue(),
					CreateBooleanFalse(),
				}),
				CreateTokenFromToken(tokens.PTR_TOKEN),
				CreateList([]node.Node{
					CreateBooleanFalse(),
					CreateBooleanTrue(),
				}),
			),
			Result: CreateList([]node.Node{
				CreateBooleanTrue(),
				CreateBooleanFalse(),
				CreateBooleanFalse(),
				CreateBooleanTrue(),
			}),
		},

		// Boolean OR
		{
			AST: node.CreateBinaryExpression(
				CreateBooleanTrue(),
				CreateTokenFromToken(tokens.OR_TOKEN),
				CreateBooleanTrue(),
			),
			Result: CreateBooleanTrue(),
		},
		{
			AST: node.CreateBinaryExpression(
				CreateBooleanTrue(),
				CreateTokenFromToken(tokens.OR_TOKEN),
				CreateBooleanFalse(),
			),
			Result: CreateBooleanTrue(),
		},
		{
			AST: node.CreateBinaryExpression(
				CreateBooleanFalse(),
				CreateTokenFromToken(tokens.OR_TOKEN),
				CreateBooleanTrue(),
			),
			Result: CreateBooleanTrue(),
		},
		{
			AST: node.CreateBinaryExpression(
				CreateBooleanFalse(),
				CreateTokenFromToken(tokens.OR_TOKEN),
				CreateBooleanFalse(),
			),
			Result: CreateBooleanFalse(),
		},

		// Boolean AND
		{
			AST: node.CreateBinaryExpression(
				CreateBooleanTrue(),
				CreateTokenFromToken(tokens.AND_TOKEN),
				CreateBooleanTrue(),
			),
			Result: CreateBooleanTrue(),
		},
		{
			AST: node.CreateBinaryExpression(
				CreateBooleanTrue(),
				CreateTokenFromToken(tokens.AND_TOKEN),
				CreateBooleanFalse(),
			),
			Result: CreateBooleanFalse(),
		},
		{
			AST: node.CreateBinaryExpression(
				CreateBooleanFalse(),
				CreateTokenFromToken(tokens.AND_TOKEN),
				CreateBooleanTrue(),
			),
			Result: CreateBooleanFalse(),
		},
		{
			AST: node.CreateBinaryExpression(
				CreateBooleanFalse(),
				CreateTokenFromToken(tokens.AND_TOKEN),
				CreateBooleanFalse(),
			),
			Result: CreateBooleanFalse(),
		},
	}

	for i, test := range tests {
		ast := []node.Node{
			test.AST,
		}

		actualResults := getResults(ast)
		expectedResults := []node.Node{
			test.Result,
		}

		AssertNodesEqual(t, i, actualResults, expectedResults)
	}
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
	AssertNodesEqual(t, 0, actualResults, expectedResults)
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

	AssertNodesEqual(t, 0, expectedResults, actualResults)
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

	AssertNodesEqual(t, 0, expectedResults, actualResults)
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
	AssertNodesEqual(t, 0, ast, actualResults)
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

	AssertNodesEqual(t, 0, expectedResults, actualResults)
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
			CreateNumber("777"),
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

	AssertNodesEqual(t, 0, expectedResults, actualResults)
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
		// Call the same function multiple times with different parameters to ensure a different result is returned each time
		CreateFunctionCall(
			CreateIdentifier("divide"),
			[]node.Node{
				CreateNumber("10"),
				CreateNumber("2"),
			},
		),
		CreateFunctionCall(
			CreateIdentifier("divide"),
			[]node.Node{
				CreateNumber("6"),
				CreateNumber("3"),
			},
		),
	}

	actualResults := getResults(ast)

	expectedResults := []node.Node{
		CreateFunctionReturnValue(CreateNumber("5").Ptr()),
		CreateFunctionReturnValue(CreateNumber("2").Ptr()),
	}

	AssertNodesEqual(t, 0, expectedResults, actualResults)
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
	AssertNodesEqual(t, 0, expectedResults, actualResults)
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
	AssertNodesEqual(t, 0, expectedResults, actualResults)
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
	AssertNodesEqual(t, 0, expectedResults, actualResults)
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

	for i, test := range tests {
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
		AssertNodesEqual(t, i, expectedResults, actualResults)
	}
}

func TestEvaluator_Slice(t *testing.T) {
	ast := []node.Node{
		CreateFunctionCall(
			CreateIdentifier("slice"),
			[]node.Node{
				CreateList([]node.Node{
					CreateNumber("0"),
					CreateNumber("1"),
					CreateNumber("2"),
					CreateNumber("3"),
					CreateNumber("4"),
					CreateNumber("5"),
				}),
				CreateNumber("1"),
				CreateNumber("4"),
			},
		),
	}

	actualResults := getResults(ast)
	expectedResults := []node.Node{
		CreateList([]node.Node{
			CreateNumber("1"),
			CreateNumber("2"),
			CreateNumber("3"),
			CreateNumber("4"),
		}),
	}
	AssertNodesEqual(t, 0, expectedResults, actualResults)
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
			tokens.AT_TOKEN,
			CreateNumber("1"),
		),
	}
	actualResults := getResults(ast)
	expectedResults := []node.Node{
		CreateNumber("2"),
	}
	AssertNodesEqual(t, 0, expectedResults, actualResults)
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
				node.CreateBinaryExpression(
					CreateIdentifier("a"),
					CreateTokenFromToken(tokens.PLUS_TOKEN),
					CreateIdentifier("b"),
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
	AssertNodesEqual(t, 0, expectedResults, actualResults)
}

func TestEvaluator_CompareOperators(t *testing.T) {

	tests := []struct {
		BinaryExpressionAST node.Node
		ExpectedResult      node.Node
	}{
		{
			BinaryExpressionAST: node.CreateBinaryExpression(
				CreateNumber("7"),
				CreateTokenFromToken(tokens.EQ_TOKEN),
				CreateNumber("7"),
			),
			ExpectedResult: CreateBooleanTrue(),
		},
		{
			BinaryExpressionAST: node.CreateBinaryExpression(
				CreateNumber("146"),
				CreateTokenFromToken(tokens.EQ_TOKEN),
				CreateNumber("66"),
			),
			ExpectedResult: CreateBooleanFalse(),
		},
		{
			BinaryExpressionAST: node.CreateBinaryExpression(
				CreateBooleanTrue(),
				CreateTokenFromToken(tokens.EQ_TOKEN),
				CreateRawString("true"),
			),
			ExpectedResult: CreateBooleanFalse(),
		},
	}

	for i, test := range tests {
		ast := []node.Node{test.BinaryExpressionAST}

		actualResults := getResults(ast)
		expectedResults := []node.Node{
			test.ExpectedResult,
		}
		AssertNodesEqual(t, i, expectedResults, actualResults)
	}
}

func TestEvaluator_WhenExpression(t *testing.T) {

	tests := []struct {
		WhenCondition string
		ExpectedValue string
	}{
		{"0", "5"},
		{"1", "10"},
		{"2", "15"},
	}

	for i, test := range tests {
		ast := []node.Node{
			CreateWhenNode(
				CreateNumber(test.WhenCondition),
				[]node.Node{
					CreateWhenCaseNode(
						CreateNumber("0"),
						CreateBlockStatements([]node.Node{
							CreateNumber("5"),
						}),
					),
					CreateWhenCaseNode(
						CreateNumber("1"),
						CreateBlockStatements([]node.Node{
							CreateNumber("10"),
						}),
					),
				},
				CreateBlockStatements([]node.Node{
					CreateNumber("15"),
				}),
			),
		}

		actualResults := getResults(ast)
		expectedResults := []node.Node{
			CreateList([]node.Node{
				CreateBooleanTrue(),
				CreateNumber(test.ExpectedValue),
			}),
		}
		AssertNodesEqual(t, i, expectedResults, actualResults)
	}
}

func TestEvaluator_WhenExpressionIfStatement(t *testing.T) {

	tests := []struct {
		WhenCondition string
		ExpectedValue string
	}{
		{"true", "5"},
		{"false", "10"},
	}

	for i, test := range tests {
		ast := []node.Node{
			CreateAssignmentStatement("number", CreateNumber("0")),
			CreateWhenNode(
				CreateBoolean(test.WhenCondition),
				[]node.Node{
					CreateWhenCaseNode(
						node.CreateBinaryExpression(
							CreateIdentifier("number"),
							CreateTokenFromToken(tokens.EQ_TOKEN),
							CreateNumber("0"),
						),
						CreateBlockStatements([]node.Node{
							CreateNumber("5"),
						}),
					),
					CreateWhenCaseNode(
						node.CreateBinaryExpression(
							CreateIdentifier("number"),
							CreateTokenFromToken(tokens.EQ_TOKEN),
							CreateNumber("1"),
						),
						CreateBlockStatements([]node.Node{
							CreateNumber("10"),
						}),
					),
				},
				CreateBlockStatements([]node.Node{
					CreateNumber("15"),
				}),
			),
		}

		actualResults := getResults(ast)
		expectedResults := []node.Node{
			CreateList([]node.Node{
				CreateBooleanTrue(),
				CreateNumber(test.ExpectedValue),
			}),
		}
		AssertNodesEqual(t, i, expectedResults, actualResults)
	}
}

/* * * * * * * *
 * ERROR TESTS *
 * * * * * * * */

func TestEvaluator_InvalidUnaryOperatorError(t *testing.T) {
	// an if-statement in the global scope containing a return statement should throw an error
	ast := []node.Node{
		node.CreateUnaryExpression(
			CreateTokenFromToken(tokens.PLUS_TOKEN),
			CreateNumber("1"),
		),
	}

	actualError := getError(t, ast)
	expectedError := "error at line 1: invalid unary operator: PLUS (\"+\")"

	if expectedError != actualError {
		t.Fatalf("Expected error: %s, Actual Error: %s", expectedError, actualError)
	}
}

func TestEvaluator_InvalidBinaryOperatorError(t *testing.T) {
	ast := []node.Node{
		node.CreateBinaryExpression(
			CreateNumber("1"),
			CreateTokenFromToken(tokens.NOT_TOKEN),
			CreateNumber("1"),
		),
	}

	actualError := getError(t, ast)
	expectedError := "error at line 1: invalid binary operator: NOT (\"not\")"

	if expectedError != actualError {
		t.Fatalf("Expected error: %s, Actual Error: %s", expectedError, actualError)
	}
}

func TestEvaluator_FunctionCallError(t *testing.T) {
	// an if-statement in the global scope containing a return statement should throw an error
	ast := []node.Node{
		CreateFunctionCall(
			CreateNumber("1"),
			[]node.Node{
				CreateNumber("1"),
				CreateNumber("2"),
			},
		),
	}

	actualError := getError(t, ast)
	expectedError := "error at line 1: cannot make function call on type Number (\"1\")"

	if expectedError != actualError {
		t.Fatalf("Expected error: %s, Actual Error: %s", expectedError, actualError)
	}
}

func TestEvaluator_FunctionCallWrongNumberOfArgumentsError(t *testing.T) {
	// an if-statement in the global scope containing a return statement should throw an error
	ast := []node.Node{
		CreateFunctionCall(
			CreateFunction([]node.Node{
				CreateIdentifier("a"),
			},
				[]node.Node{},
			),
			[]node.Node{
				CreateNumber("1"),
				CreateNumber("2"),
			},
		),
	}

	actualError := getError(t, ast)
	expectedError := "error at line 1: expected 1 arguments, got 2"

	if expectedError != actualError {
		t.Fatalf("Expected error: %s, Actual Error: %s", expectedError, actualError)
	}
}

func TestEvaluator_IndexValueNotIntegerError(t *testing.T) {
	// an if-statement in the global scope containing a return statement should throw an error
	ast := []node.Node{
		node.CreateBinaryExpression(
			CreateList([]node.Node{
				CreateNumber("1"),
			}),
			CreateTokenFromToken(tokens.AT_TOKEN),
			CreateNumber("3.4"),
		),
	}

	actualError := getError(t, ast)
	expectedError := "error at line 1: list index must be an integer"

	if expectedError != actualError {
		t.Fatalf("Expected error: %s, Actual Error: %s", expectedError, actualError)
	}
}

func TestEvaluator_IndexOutOfRangeError(t *testing.T) {
	// an if-statement in the global scope containing a return statement should throw an error
	ast := []node.Node{
		node.CreateBinaryExpression(
			CreateList([]node.Node{
				CreateNumber("1"),
			}),
			CreateTokenFromToken(tokens.AT_TOKEN),
			CreateNumber("3"),
		),
	}

	actualError := getError(t, ast)
	expectedError := "error at line 1: index 3 out of range. Length of list: 1"

	if expectedError != actualError {
		t.Fatalf("Expected error: %s, Actual Error: %s", expectedError, actualError)
	}
}

func TestEvaluator_IndexInvalidTypeError(t *testing.T) {
	// an if-statement in the global scope containing a return statement should throw an error
	ast := []node.Node{
		node.CreateBinaryExpression(
			CreateNumber("3"),
			CreateTokenFromToken(tokens.AT_TOKEN),
			CreateNumber("3"),
		),
	}

	actualError := getError(t, ast)
	expectedError := "error at line 1: invalid types for index: Number (\"3\") and Number (\"3\")"

	if expectedError != actualError {
		t.Fatalf("Expected error: %s, Actual Error: %s", expectedError, actualError)
	}
}

func TestEvaluator_AddInvalidTypesError(t *testing.T) {
	ast := []node.Node{
		node.CreateBinaryExpression(
			CreateRawString("hello"),
			CreateTokenFromToken(tokens.PLUS_TOKEN),
			CreateRawString(" world!"),
		),
	}
	actualError := getError(t, ast)
	expectedError := "error at line 1: cannot add types String (\"hello\") and String (\" world!\")"

	if expectedError != actualError {
		t.Fatalf("Expected error: %s, Actual Error: %s", expectedError, actualError)
	}
}

func TestEvaluator_SubtractInvalidTypesError(t *testing.T) {
	ast := []node.Node{
		node.CreateBinaryExpression(
			CreateRawString("hello"),
			CreateTokenFromToken(tokens.MINUS_TOKEN),
			CreateRawString(" world!"),
		),
	}
	actualError := getError(t, ast)
	expectedError := "error at line 1: cannot subtract types String (\"hello\") and String (\" world!\")"

	if expectedError != actualError {
		t.Fatalf("Expected error: %s, Actual Error: %s", expectedError, actualError)
	}
}

func TestEvaluator_MultiplyInvalidTypesError(t *testing.T) {
	ast := []node.Node{
		node.CreateBinaryExpression(
			CreateRawString("hello"),
			CreateTokenFromToken(tokens.ASTERISK_TOKEN),
			CreateRawString(" world!"),
		),
	}
	actualError := getError(t, ast)
	expectedError := "error at line 1: cannot multiply types String (\"hello\") and String (\" world!\")"

	if expectedError != actualError {
		t.Fatalf("Expected error: %s, Actual Error: %s", expectedError, actualError)
	}
}

func TestEvaluator_DivideInvalidTypesError(t *testing.T) {
	ast := []node.Node{
		node.CreateBinaryExpression(
			CreateBooleanTrue(),
			CreateTokenFromToken(tokens.FORWARD_SLASH_TOKEN),
			CreateBooleanFalse(),
		),
	}

	actualError := getError(t, ast)
	expectedError := "error at line 1: cannot divide types Boolean (\"true\") and Boolean (\"false\")"

	if expectedError != actualError {
		t.Fatalf("Expected error: %s, Actual Error: %s", expectedError, actualError)
	}
}

func TestEvaluator_DivideZeroError(t *testing.T) {
	ast := []node.Node{
		node.CreateBinaryExpression(
			CreateNumber("1"),
			CreateTokenFromToken(tokens.FORWARD_SLASH_TOKEN),
			CreateNumber("0"),
		),
	}

	actualError := getError(t, ast)
	expectedError := "error at line 1: cannot divide by zero"

	if expectedError != actualError {
		t.Fatalf("Expected error: %s, Actual Error: %s", expectedError, actualError)
	}
}

func TestEvaluator_PointerInvalidTypesError(t *testing.T) {
	ast := []node.Node{
		node.CreateBinaryExpression(
			CreateBooleanTrue(),
			CreateTokenFromToken(tokens.PTR_TOKEN),
			CreateBooleanFalse(),
		),
	}

	actualError := getError(t, ast)
	expectedError := "error at line 1: cannot use pointer on types Boolean (\"true\") and Boolean (\"false\")"

	if expectedError != actualError {
		t.Fatalf("Expected error: %s, Actual Error: %s", expectedError, actualError)
	}
}

func TestEvaluator_BangError(t *testing.T) {
	ast := []node.Node{
		node.CreateUnaryExpression(
			CreateTokenFromToken(tokens.NOT_TOKEN),
			CreateNumber("1"),
		),
	}
	actualError := getError(t, ast)
	expectedError := "error at line 1: invalid type for bang operator: Number (\"1\")"

	if expectedError != actualError {
		t.Fatalf("Expected error: %s, Actual Error: %s", expectedError, actualError)
	}
}

func TestEvaluator_MinusUnayError(t *testing.T) {
	ast := []node.Node{
		node.CreateUnaryExpression(
			CreateTokenFromToken(tokens.MINUS_TOKEN),
			CreateBooleanFalse(),
		),
	}
	actualError := getError(t, ast)
	expectedError := "error at line 1: invalid type for minus operator: Boolean (\"false\")"

	if expectedError != actualError {
		t.Fatalf("Expected error: %s, Actual Error: %s", expectedError, actualError)
	}
}

func TestEvaluator_BuiltinSliceIndexErrors(t *testing.T) {

	tests := []struct {
		StartIndex node.Node
		EndIndex   node.Node
		Error      string
	}{
		{
			StartIndex: CreateNumber("-1"),
			EndIndex:   CreateNumber("4"),
			Error:      "error at line 1: index of -1 out of range (0 to 5)",
		},
		{
			StartIndex: CreateNumber("0"),
			EndIndex:   CreateNumber("6"),
			Error:      "error at line 1: index of 6 out of range (0 to 5)",
		},
		{
			StartIndex: CreateNumber("4.5"),
			EndIndex:   CreateNumber("6"),
			Error:      "error at line 1: list index must be an integer",
		},
		{
			StartIndex: CreateNumber("4"),
			EndIndex:   CreateNumber("6.6"),
			Error:      "error at line 1: list index must be an integer",
		},
		{
			StartIndex: CreateRawString("hello!"),
			EndIndex:   CreateNumber("6.6"),
			Error:      "error at line 1: expected Number, got String (\"hello!\")",
		},
		{
			StartIndex: CreateNumber("4"),
			EndIndex:   CreateBooleanTrue(),
			Error:      "error at line 1: expected Number, got Boolean (\"true\")",
		},
		{
			StartIndex: CreateNumber("4"),
			EndIndex:   CreateNumber("2"),
			Error:      "error at line 1: start index cannot be greater than end index",
		},
	}

	list := CreateList([]node.Node{
		CreateNumber("0"),
		CreateNumber("1"),
		CreateNumber("2"),
		CreateNumber("3"),
		CreateNumber("4"),
		CreateNumber("5"),
	})

	for i, test := range tests {

		ast := []node.Node{
			CreateFunctionCall(
				CreateIdentifier("slice"),
				[]node.Node{
					list,
					test.StartIndex,
					test.EndIndex,
				},
			),
		}

		actualError := getError(t, ast)
		expectedError := test.Error

		if expectedError != actualError {
			t.Fatalf("Test #%d - Expected error: %s, Actual Error: %s", i, expectedError, actualError)
		}
	}
}

func TestEvaluator_SliceInvalidTypeError(t *testing.T) {
	ast := []node.Node{
		CreateFunctionCall(
			CreateIdentifier("slice"),
			[]node.Node{
				CreateBooleanTrue(),
				CreateNumber("0"),
				CreateNumber("2"),
			},
		),
	}

	actualError := getError(t, ast)
	expectedError := "error at line 1: expected List, got Boolean (\"true\")"

	if expectedError != actualError {
		t.Fatalf("Expected error: %s, Actual Error: %s", expectedError, actualError)
	}
}

func TestEvaluator_SliceInvalidNumberOfArgumentsError(t *testing.T) {

	tests := []struct {
		Args []node.Node
	}{
		{
			Args: []node.Node{},
		},
		{
			Args: []node.Node{
				CreateList([]node.Node{}),
			},
		},
		{
			Args: []node.Node{
				CreateList([]node.Node{}),
				CreateNumber("3"),
			},
		},
		{
			Args: []node.Node{
				CreateList([]node.Node{}),
				CreateNumber("3"),
				CreateNumber("3"),
				CreateNumber("3"),
			},
		},
	}

	for i, test := range tests {

		ast := []node.Node{
			CreateFunctionCall(
				CreateIdentifier("slice"),
				test.Args,
			),
		}

		actualError := getError(t, ast)
		expectedError := fmt.Sprintf("error at line 1: incorrect number of arguments. expected 3, got %d", len(test.Args))

		if expectedError != actualError {
			t.Fatalf("Test #%d - Expected error: %s, Actual Error: %s", i, expectedError, actualError)
		}
	}
}

/* * * * * * * * * * * * *
 * EVALUATOR TEST UTILS  *
 * * * * * * * * * * * * */

func getResults(ast []node.Node) []node.Node {
	evaluatorObj := evaluator.NewEvaluator(ast)
	actualResults, err := evaluatorObj.Evaluate()
	if err != nil {
		panic(err.Error())
	}
	return *actualResults
}

func getError(t *testing.T, ast []node.Node) string {
	evaluatorObj := evaluator.NewEvaluator(ast)
	_, err := evaluatorObj.Evaluate()

	if err == nil {
		t.Fatal("error is nil")
	}
	return err.Error()
}
