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

		actualResults := getEvaluatorResults(ast)
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

		actualResults := getEvaluatorResults(ast)
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

	actualResults := getEvaluatorResults(ast)
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

		actualResults := getEvaluatorResults(ast)
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

		actualResults := getEvaluatorResults(ast)
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

		actualResults := getEvaluatorResults(ast)
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
		// Math expression
		{
			AST: node.CreateBinaryExpression(
				CreateNumber("1"),
				CreateTokenFromToken(tokens.PLUS_TOKEN),
				CreateNumber("1"),
			),
			Result: CreateNumber("2"),
		},

		// List append
		{
			AST: node.CreateBinaryExpression(
				CreateList([]node.Node{
					CreateBooleanTrue(),
					CreateBooleanFalse(),
				}),
				CreateTokenFromToken(tokens.SEND_TOKEN),
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
				CreateTokenFromToken(tokens.SEND_TOKEN),
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

		// "in" operator: 5 in (1, 3, 5)
		{
			AST: node.CreateBinaryExpression(
				CreateNumber("5"),
				CreateTokenFromToken(tokens.IN_TOKEN),
				CreateList([]node.Node{
					CreateNumber("1"),
					CreateNumber("3"),
					CreateNumber("5"),
				}),
			),
			Result: CreateBooleanTrue(),
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

		actualResults := getEvaluatorResults(ast)
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

	actualResults := getEvaluatorResults(ast)
	expectedResults := []node.Node{
		CreateNumber("4"),
	}
	AssertNodesEqual(t, 0, actualResults, expectedResults)
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

	actualResults := getEvaluatorResults(ast)
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

	actualResults := getEvaluatorResults(ast)

	expectedReturnValue := CreateNumber("8")
	expectedResults := []node.Node{
		CreateBlockStatementReturnValue(&expectedReturnValue),
	}

	AssertNodesEqual(t, 0, expectedResults, actualResults)
}

func TestEvaluator_FunctionCallKeywordArgs(t *testing.T) {
	functionNode := CreateFunction(
		[]node.Node{
			CreateIdentifier("a"),
			CreateAssignmentStatement("b", CreateNumber("1")),
		},
		[]node.Node{
			node.CreateBinaryExpression(
				CreateIdentifier("a"),
				tokens.PLUS_TOKEN,
				CreateIdentifier("b"),
			),
		},
	)

	tests := []struct {
		Function      node.Node
		CallParams    []node.Node
		ExpectedValue node.Node
	}{
		{
			Function: functionNode,
			CallParams: []node.Node{
				CreateNumber("10"),
			},
			ExpectedValue: CreateNumber("11"),
		},
		{
			Function: functionNode,
			CallParams: []node.Node{
				CreateNumber("10"),
				CreateNumber("5"),
			},
			ExpectedValue: CreateNumber("15"),
		},
		{
			// A function where the first parameter has a default value
			Function: CreateFunction(
				[]node.Node{
					CreateAssignmentStatement("a", CreateNumber("1")),
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
			CallParams: []node.Node{
				CreateNumber("4"),
				CreateNumber("3"),
			},
			ExpectedValue: CreateNumber("7"),
		},
	}

	for i, test := range tests {
		ast := []node.Node{
			CreateFunctionCall(test.Function, test.CallParams),
		}

		actualResults := getEvaluatorResults(ast)

		expectedReturnValue := test.ExpectedValue
		expectedResults := []node.Node{
			CreateBlockStatementReturnValue(&expectedReturnValue),
		}

		AssertNodesEqual(t, i, expectedResults, actualResults)
	}
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

	actualResults := getEvaluatorResults(ast)

	expectedResults := []node.Node{
		CreateBlockStatementReturnValue(CreateNumber("5").Ptr()),
		CreateBlockStatementReturnValue(CreateNumber("2").Ptr()),
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

	actualResults := getEvaluatorResults(ast)

	expectedReturnValue := CreateNumber("7")
	expectedResults := []node.Node{
		CreateBlockStatementReturnValue(&expectedReturnValue),
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

	actualResults := getEvaluatorResults(ast)
	expectedResults := []node.Node{
		CreateBlockStatementReturnValue(nil),
	}
	AssertNodesEqual(t, 0, expectedResults, actualResults)
}

func TestEvaluator_ListIndex(t *testing.T) {

	tests := []struct {
		Collection    node.Node
		Index         node.Node
		ExpectedValue node.Node
	}{
		{
			Collection: CreateList([]node.Node{
				CreateNumber("1"),
				CreateNumber("2"),
				CreateNumber("3"),
			}),
			Index:         CreateNumber("1"),
			ExpectedValue: CreateNumber("2"),
		},
		{
			Collection:    CreateRawString("hello, world!"),
			Index:         CreateNumber("2"),
			ExpectedValue: CreateRawString("l"),
		},
	}

	for i, test := range tests {
		ast := []node.Node{
			node.CreateBinaryExpression(
				test.Collection,
				tokens.AT_TOKEN,
				test.Index,
			),
		}
		actualResults := getEvaluatorResults(ast)
		expectedResults := []node.Node{
			test.ExpectedValue,
		}
		AssertNodesEqual(t, i, expectedResults, actualResults)
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
				CreateBuiltinFunctionIdentifier("unwrap"),
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

	actualResults := getEvaluatorResults(ast)
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

		// Less than
		{
			BinaryExpressionAST: node.CreateBinaryExpression(
				CreateNumber("5"),
				CreateTokenFromToken(tokens.LT_TOKEN),
				CreateNumber("5"),
			),
			ExpectedResult: CreateBooleanFalse(),
		},
		{
			BinaryExpressionAST: node.CreateBinaryExpression(
				CreateNumber("4"),
				CreateTokenFromToken(tokens.LT_TOKEN),
				CreateNumber("5"),
			),
			ExpectedResult: CreateBooleanTrue(),
		},
		{
			BinaryExpressionAST: node.CreateBinaryExpression(
				CreateNumber("3.14159"),
				CreateTokenFromToken(tokens.LT_TOKEN),
				CreateNumber("36.9"),
			),
			ExpectedResult: CreateBooleanTrue(),
		},
		{
			BinaryExpressionAST: node.CreateBinaryExpression(
				CreateNumber("3.14159"),
				CreateTokenFromToken(tokens.LT_TOKEN),
				CreateNumber("3.14159"),
			),
			ExpectedResult: CreateBooleanFalse(),
		},
	}

	for i, test := range tests {
		ast := []node.Node{test.BinaryExpressionAST}

		actualResults := getEvaluatorResults(ast)
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
						[]node.Node{
							CreateNumber("5"),
						},
					),
					CreateWhenCaseNode(
						CreateNumber("1"),
						[]node.Node{
							CreateNumber("10"),
						},
					),
				},
				[]node.Node{
					CreateNumber("15"),
				},
			),
		}

		actualResults := getEvaluatorResults(ast)
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
						[]node.Node{
							CreateNumber("5"),
						},
					),
					CreateWhenCaseNode(
						node.CreateBinaryExpression(
							CreateIdentifier("number"),
							CreateTokenFromToken(tokens.EQ_TOKEN),
							CreateNumber("1"),
						),
						[]node.Node{
							CreateNumber("10"),
						},
					),
				},
				[]node.Node{
					CreateNumber("15"),
				},
			),
		}

		actualResults := getEvaluatorResults(ast)
		expectedResults := []node.Node{
			CreateList([]node.Node{
				CreateBooleanTrue(),
				CreateNumber(test.ExpectedValue),
			}),
		}
		AssertNodesEqual(t, i, expectedResults, actualResults)
	}
}

func TestEvaluator_VariablesFromOuterScopes(t *testing.T) {
	// "var" is defined outside the function "f", but is still accessible within that function.
	ast := []node.Node{
		CreateAssignmentStatement("var", CreateNumber("10")),
		CreateAssignmentStatement(
			"f",
			CreateFunction(
				[]node.Node{CreateIdentifier("c")},
				[]node.Node{
					node.CreateBinaryExpression(
						CreateIdentifier("c"),
						CreateTokenFromToken(tokens.PLUS_TOKEN),
						CreateIdentifier("var"),
					),
				},
			),
		),
		node.CreateBinaryExpression(
			CreateIdentifier("f"),
			CreateTokenFromToken(tokens.SEND_TOKEN),
			CreateList([]node.Node{
				CreateNumber("2"),
			}),
		),
	}

	actualResults := getEvaluatorResults(ast)
	expectedResults := []node.Node{
		CreateList([]node.Node{
			CreateBooleanTrue(),
			CreateNumber("12"),
		}),
	}
	AssertNodesEqual(t, 0, expectedResults, actualResults)
}

func TestEvaluator_ForLoop(t *testing.T) {

	tests := []struct {
		BlockStatement []node.Node
		ReturnValue    node.Node
	}{
		{
			BlockStatement: []node.Node{
				CreateAssignmentStatement(
					"i",
					CreateIdentifier("e"),
				),
			},
			ReturnValue: CreateList([]node.Node{
				CreateBlockStatementReturnValue(nil),
				CreateBlockStatementReturnValue(nil),
				CreateBlockStatementReturnValue(nil),
				CreateBlockStatementReturnValue(nil),
			}),
		},
		{
			BlockStatement: []node.Node{
				node.CreateBinaryExpression(
					CreateIdentifier("e"),
					CreateTokenFromToken(tokens.ASTERISK_TOKEN),
					CreateIdentifier("e"),
				),
			},
			ReturnValue: CreateList([]node.Node{
				CreateBlockStatementReturnValue(CreateNumber("1").Ptr()),
				CreateBlockStatementReturnValue(CreateNumber("4").Ptr()),
				CreateBlockStatementReturnValue(CreateNumber("9").Ptr()),
				CreateBlockStatementReturnValue(CreateNumber("16").Ptr()),
			}),
		},
	}

	for i, test := range tests {
		ast := []node.Node{
			CreateForLoop(
				CreateIdentifier("e"),
				CreateList([]node.Node{
					CreateNumber("1"),
					CreateNumber("2"),
					CreateNumber("3"),
					CreateNumber("4"),
				}),
				test.BlockStatement,
			),
		}

		actualResults := getEvaluatorResults(ast)
		expectedResults := []node.Node{
			test.ReturnValue,
		}
		AssertNodesEqual(t, i, expectedResults, actualResults)
	}
}

func TestEvaluator_WhileLoop(t *testing.T) {
	ast := []node.Node{
		CreateAssignmentStatement(
			"i",
			CreateNumber("0"),
		),
		CreateWhileLoop(
			node.CreateBinaryExpression(
				CreateIdentifier("i"),
				CreateTokenFromToken(tokens.LT_TOKEN),
				CreateNumber("10"),
			),
			[]node.Node{
				CreateAssignmentStatement(
					"i",
					node.CreateBinaryExpression(
						CreateIdentifier("i"),
						CreateTokenFromToken(tokens.PLUS_TOKEN),
						CreateNumber("1"),
					),
				),
			},
		),
		CreateIdentifier("i"),
	}
	actualResults := getEvaluatorResults(ast)
	expectedResults := []node.Node{
		CreateNumber("10"),
	}
	AssertNodesEqual(t, 0, expectedResults, actualResults)
}

func TestEvaluator_BreakStatement(t *testing.T) {
	ast := []node.Node{
		CreateAssignmentStatement("i", CreateNumber("0")),
		CreateWhileLoop(
			node.CreateBinaryExpression(
				CreateIdentifier("i"),
				CreateTokenFromToken(tokens.LT_TOKEN),
				CreateNumber("10"),
			),
			[]node.Node{
				CreateAssignmentStatement(
					"i",
					node.CreateBinaryExpression(
						CreateIdentifier("i"),
						CreateTokenFromToken(tokens.PLUS_TOKEN),
						CreateNumber("1"),
					),
				),
				CreateBreakStatement(),
			},
		),
		CreateIdentifier("i"),
	}
	actualResults := getEvaluatorResults(ast)
	expectedResults := []node.Node{
		CreateNumber("1"),
	}
	AssertNodesEqual(t, 0, expectedResults, actualResults)
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

	actualError := getEvaluatorError(t, ast)
	expectedError := "error at line 1: invalid unary operator: PLUS (\"+\")"

	AssertErrorEqual(t, 0, expectedError, actualError)
}

func TestEvaluator_InvalidBinaryOperatorError(t *testing.T) {
	ast := []node.Node{
		node.CreateBinaryExpression(
			CreateNumber("1"),
			CreateTokenFromToken(tokens.NOT_TOKEN),
			CreateNumber("1"),
		),
	}

	actualError := getEvaluatorError(t, ast)
	expectedError := "error at line 1: invalid binary operator: NOT (\"not\")"

	AssertErrorEqual(t, 0, expectedError, actualError)
}

func TestEvaluator_FunctionCallErrors(t *testing.T) {

	tests := []struct {
		FunctionCall node.Node
		Error        string
	}{
		{
			FunctionCall: CreateFunctionCall(
				CreateNumber("1"),
				[]node.Node{
					CreateNumber("1"),
					CreateNumber("2"),
				},
			),
			Error: "error at line 1: cannot make function call on type Number (\"1\")",
		},
		{
			FunctionCall: CreateFunctionCall(
				CreateFunction(
					[]node.Node{
						CreateAssignmentStatement("a", CreateNumber("1")),
						CreateIdentifier("b"),
					},
					[]node.Node{},
				),
				[]node.Node{
					CreateNumber("1"),
				},
			),
			Error: "error at line 1: Function paramter \"b\" does not have a value. Either add 1 more parameters to the function call or assign \"b\" a default value in the function definition.",
		},
		{
			FunctionCall: CreateFunctionCall(
				CreateFunction(
					[]node.Node{
						CreateAssignmentStatement("a", CreateNumber("1")),
						CreateAssignmentStatement("b", CreateNumber("2")),
						CreateIdentifier("c"),
					},
					[]node.Node{},
				),
				[]node.Node{
					CreateNumber("1"),
				},
			),
			Error: "error at line 1: Function paramter \"c\" does not have a value. Either add 2 more parameters to the function call or assign \"c\" a default value in the function definition.",
		},
		{
			FunctionCall: CreateFunctionCall(
				CreateFunction(
					[]node.Node{
						CreateAssignmentStatement("a", CreateNumber("1")),
						CreateAssignmentStatement("b", CreateNumber("2")),
						CreateIdentifier("c"),
					},
					[]node.Node{},
				),
				[]node.Node{
					CreateNumber("5"),
					CreateNumber("5"),
				},
			),
			Error: "error at line 1: Function paramter \"c\" does not have a value. Either add 1 more parameters to the function call or assign \"c\" a default value in the function definition.",
		},
	}

	for i, test := range tests {
		ast := []node.Node{test.FunctionCall}

		actualError := getEvaluatorError(t, ast)
		AssertErrorEqual(t, i, test.Error, actualError)
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

	actualError := getEvaluatorError(t, ast)
	expectedError := "error at line 1: expected 1 arguments, got 2"

	AssertErrorEqual(t, 0, expectedError, actualError)
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

	actualError := getEvaluatorError(t, ast)
	expectedError := "error at line 1: list index must be an integer"

	AssertErrorEqual(t, 0, expectedError, actualError)
}

func TestEvaluator_IndexOutOfRangeError(t *testing.T) {
	tests := []struct {
		Sequence node.Node
		Index    node.Node
		Error    string
	}{
		{
			Sequence: CreateList([]node.Node{
				CreateNumber("1"),
			}),
			Index: CreateNumber("3"),
			Error: "error at line 1: index of 3 out of range (0 to 0)",
		},
		{
			Sequence: CreateRawString("test string"),
			Index:    CreateNumber("-1"),
			Error:    "error at line 1: index of -1 out of range (0 to 10)",
		},
	}

	for i, test := range tests {
		ast := []node.Node{
			node.CreateBinaryExpression(
				test.Sequence,
				CreateTokenFromToken(tokens.AT_TOKEN),
				test.Index,
			),
		}

		actualError := getEvaluatorError(t, ast)

		AssertErrorEqual(t, i, test.Error, actualError)
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

	actualError := getEvaluatorError(t, ast)
	expectedError := "error at line 1: invalid types for index: Number (\"3\") and Number (\"3\")"

	AssertErrorEqual(t, 0, expectedError, actualError)
}

func TestEvaluator_AddInvalidTypesError(t *testing.T) {
	ast := []node.Node{
		node.CreateBinaryExpression(
			CreateRawString("hello"),
			CreateTokenFromToken(tokens.PLUS_TOKEN),
			CreateRawString(" world!"),
		),
	}
	actualError := getEvaluatorError(t, ast)
	expectedError := "error at line 1: cannot add types String (\"hello\") and String (\" world!\")"

	AssertErrorEqual(t, 0, expectedError, actualError)
}

func TestEvaluator_SubtractInvalidTypesError(t *testing.T) {
	ast := []node.Node{
		node.CreateBinaryExpression(
			CreateRawString("hello"),
			CreateTokenFromToken(tokens.MINUS_TOKEN),
			CreateRawString(" world!"),
		),
	}
	actualError := getEvaluatorError(t, ast)
	expectedError := "error at line 1: cannot subtract types String (\"hello\") and String (\" world!\")"

	AssertErrorEqual(t, 0, expectedError, actualError)
}

func TestEvaluator_MultiplyInvalidTypesError(t *testing.T) {
	ast := []node.Node{
		node.CreateBinaryExpression(
			CreateRawString("hello"),
			CreateTokenFromToken(tokens.ASTERISK_TOKEN),
			CreateRawString(" world!"),
		),
	}
	actualError := getEvaluatorError(t, ast)
	expectedError := "error at line 1: cannot multiply types String (\"hello\") and String (\" world!\")"

	AssertErrorEqual(t, 0, expectedError, actualError)
}

func TestEvaluator_DivideInvalidTypesError(t *testing.T) {
	ast := []node.Node{
		node.CreateBinaryExpression(
			CreateBooleanTrue(),
			CreateTokenFromToken(tokens.FORWARD_SLASH_TOKEN),
			CreateBooleanFalse(),
		),
	}

	actualError := getEvaluatorError(t, ast)
	expectedError := "error at line 1: cannot divide types Boolean (\"true\") and Boolean (\"false\")"

	AssertErrorEqual(t, 0, expectedError, actualError)
}

func TestEvaluator_DivideZeroError(t *testing.T) {
	ast := []node.Node{
		node.CreateBinaryExpression(
			CreateNumber("1"),
			CreateTokenFromToken(tokens.FORWARD_SLASH_TOKEN),
			CreateNumber("0"),
		),
	}

	actualError := getEvaluatorError(t, ast)
	expectedError := "error at line 1: cannot divide by zero"

	AssertErrorEqual(t, 0, expectedError, actualError)
}

func TestEvaluator_SendInvalidTypesError(t *testing.T) {
	ast := []node.Node{
		node.CreateBinaryExpression(
			CreateBooleanTrue(),
			CreateTokenFromToken(tokens.SEND_TOKEN),
			CreateBooleanFalse(),
		),
	}

	actualError := getEvaluatorError(t, ast)
	expectedError := "error at line 1: cannot use send on types Boolean (\"true\") and Boolean (\"false\")"

	AssertErrorEqual(t, 0, expectedError, actualError)
}

func TestEvaluator_BangError(t *testing.T) {
	ast := []node.Node{
		node.CreateUnaryExpression(
			CreateTokenFromToken(tokens.NOT_TOKEN),
			CreateNumber("1"),
		),
	}
	actualError := getEvaluatorError(t, ast)
	expectedError := "error at line 1: invalid type for bang operator: Number (\"1\")"

	AssertErrorEqual(t, 0, expectedError, actualError)
}

func TestEvaluator_MinusUnayError(t *testing.T) {
	ast := []node.Node{
		node.CreateUnaryExpression(
			CreateTokenFromToken(tokens.MINUS_TOKEN),
			CreateBooleanFalse(),
		),
	}
	actualError := getEvaluatorError(t, ast)
	expectedError := "error at line 1: invalid type for minus operator: Boolean (\"false\")"

	AssertErrorEqual(t, 0, expectedError, actualError)
}

func TestEvaluator_BreakStatementOutsideLoopError(t *testing.T) {

	tests := []node.Node{
		CreateBreakStatement(),

		// "break" in "when" expression
		CreateWhenNode(
			CreateBooleanTrue(),
			[]node.Node{
				CreateWhenCaseNode(
					CreateBooleanTrue(),
					[]node.Node{
						CreateBreakStatement(),
					},
				),
			},
			[]node.Node{},
		),

		// "break" in function
		CreateFunctionCall(
			CreateFunction(
				[]node.Node{},
				[]node.Node{
					CreateBreakStatement(),
				},
			),
			[]node.Node{},
		),
	}

	for i, test := range tests {
		ast := []node.Node{
			test,
		}
		actualError := getEvaluatorError(t, ast)
		expectedError := "error at line 1: break statements not allowed outside loops"

		AssertErrorEqual(t, i, expectedError, actualError)
	}
}

func TestEvaluator_VariableNameSameAsBuiltinsError(t *testing.T) {

	builtinIdentifiers := []string{
		// Functions
		evaluator.BUILTIN_LEN,
		evaluator.BUILTIN_UNWRAP,
		evaluator.BUILTIN_SLICE,
		evaluator.BUILTIN_UNWRAP_ALL,
		evaluator.BUILTIN_RANGE,
		evaluator.BUILTIN_RANDOM,
		evaluator.BUILTIN_PRINT,
		evaluator.BUILTIN_INPUT,

		// Variables
		evaluator.BUILTIN_PI,
	}

	for i, identifier := range builtinIdentifiers {
		ast := []node.Node{
			CreateAssignmentStatement(identifier, CreateNumber("20")),
		}

		actualError := getEvaluatorError(t, ast)
		expectedError := fmt.Sprintf("error at line 1: \"%s\" is a builtin function or variable", identifier)

		AssertErrorEqual(t, i, expectedError, actualError)
	}
}
