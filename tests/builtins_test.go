package tests

import (
	"boomerang/node"
	"boomerang/tokens"
	"fmt"
	"testing"
)

func TestBuiltin_Len(t *testing.T) {
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

	actualResults := getEvaluatorResults(ast)
	expectedResults := []node.Node{
		CreateNumber("3"),
	}
	AssertNodesEqual(t, 0, expectedResults, actualResults)
}

func TestBuiltin_Unwrap(t *testing.T) {

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

		actualResults := getEvaluatorResults(ast)
		expectedResults := []node.Node{
			test.ExpectedReturnValue,
		}
		AssertNodesEqual(t, i, expectedResults, actualResults)
	}
}

func TestBuiltin_UnwrapAll(t *testing.T) {

	tests := []struct {
		BlockStatementReturnValues node.Node
		ExpectedResult             node.Node
	}{
		{
			BlockStatementReturnValues: CreateList([]node.Node{
				CreateBlockStatementReturnValue(nil),
				CreateBlockStatementReturnValue(nil),
				CreateBlockStatementReturnValue(nil),
			}),
			ExpectedResult: CreateList([]node.Node{
				CreateNumber("-1"),
				CreateNumber("-1"),
				CreateNumber("-1"),
			}),
		},
		{
			BlockStatementReturnValues: CreateList([]node.Node{
				CreateBlockStatementReturnValue(CreateNumber("5").Ptr()),
				CreateBlockStatementReturnValue(CreateNumber("10").Ptr()),
				CreateBlockStatementReturnValue(CreateNumber("15").Ptr()),
			}),
			ExpectedResult: CreateList([]node.Node{
				CreateNumber("5"),
				CreateNumber("10"),
				CreateNumber("15"),
			}),
		},
		{
			BlockStatementReturnValues: CreateList([]node.Node{
				CreateBlockStatementReturnValue(CreateNumber("5").Ptr()),
				CreateBlockStatementReturnValue(nil),
				CreateBlockStatementReturnValue(CreateNumber("15").Ptr()),
			}),
			ExpectedResult: CreateList([]node.Node{
				CreateNumber("5"),
				CreateNumber("-1"),
				CreateNumber("15"),
			}),
		},
	}

	for i, test := range tests {
		ast := []node.Node{
			CreateFunctionCall(
				CreateIdentifier("unwrap_all"),
				[]node.Node{
					test.BlockStatementReturnValues,
					CreateNumber("-1"),
				},
			),
		}
		actualResults := getEvaluatorResults(ast)
		expectedResults := []node.Node{
			test.ExpectedResult,
		}
		AssertNodesEqual(t, i, expectedResults, actualResults)
	}
}

func TestBuiltin_Slice(t *testing.T) {
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

	actualResults := getEvaluatorResults(ast)
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

func TestBuiltin_Range(t *testing.T) {

	tests := []struct {
		StartNumber  string
		EndNumber    string
		ExpectedList node.Node
	}{
		{
			StartNumber: "0",
			EndNumber:   "5",
			ExpectedList: CreateList([]node.Node{
				CreateNumber("0"),
				CreateNumber("1"),
				CreateNumber("2"),
				CreateNumber("3"),
				CreateNumber("4"),
				CreateNumber("5"),
			}),
		},
		{
			StartNumber: "5",
			EndNumber:   "10",
			ExpectedList: CreateList([]node.Node{
				CreateNumber("5"),
				CreateNumber("6"),
				CreateNumber("7"),
				CreateNumber("8"),
				CreateNumber("9"),
				CreateNumber("10"),
			}),
		},
		{
			StartNumber: "5",
			EndNumber:   "5",
			ExpectedList: CreateList([]node.Node{
				CreateNumber("5"),
			}),
		},
		{
			StartNumber:  "5",
			EndNumber:    "4",
			ExpectedList: CreateList([]node.Node{}),
		},
	}

	for i, test := range tests {
		ast := []node.Node{
			CreateFunctionCall(
				CreateIdentifier("range"),
				[]node.Node{
					CreateNumber(test.StartNumber),
					CreateNumber(test.EndNumber),
				},
			),
		}

		actualResults := getEvaluatorResults(ast)
		expectedResults := []node.Node{
			test.ExpectedList,
		}
		AssertNodesEqual(t, i, expectedResults, actualResults)
	}
}

/* * * * * * * *
 * ERROR TESTS *
 * * * * * * * */

func TestBuiltin_RangeErrors(t *testing.T) {
	tests := []struct {
		Arguments []node.Node
		Error     string
	}{
		{
			Arguments: []node.Node{},
			Error:     "error at line 1: incorrect number of arguments. expected 2, got 0",
		},
		{
			Arguments: []node.Node{
				CreateNumber("1"),
			},
			Error: "error at line 1: incorrect number of arguments. expected 2, got 1",
		},
		{
			Arguments: []node.Node{
				CreateNumber("1"),
				CreateNumber("1"),
				CreateNumber("1"),
			},
			Error: "error at line 1: incorrect number of arguments. expected 2, got 3",
		},
		{
			Arguments: []node.Node{
				CreateRawString("hello, world!"),
				CreateNumber("1"),
			},
			Error: "error at line 1: expected Number, got String",
		},
		{
			Arguments: []node.Node{
				CreateNumber("1"),
				CreateList([]node.Node{}),
			},
			Error: "error at line 1: expected Number, got List",
		},
	}

	for i, test := range tests {

		ast := []node.Node{
			CreateFunctionCall(
				CreateIdentifier("range"),
				test.Arguments,
			),
		}

		actualError := getEvaluatorError(t, ast)
		expectedError := test.Error

		AssertErrorEqual(t, i, expectedError, actualError)
	}
}

func TestBuiltin_SliceInvalidNumberOfArgumentsError(t *testing.T) {

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

		actualError := getEvaluatorError(t, ast)
		expectedError := fmt.Sprintf("error at line 1: incorrect number of arguments. expected 3, got %d", len(test.Args))

		AssertErrorEqual(t, i, expectedError, actualError)
	}
}

func TestBuiltin_UnwrapErrors(t *testing.T) {

	tests := []struct {
		Args  []node.Node
		Error string
	}{
		{
			Args: []node.Node{
				CreateList([]node.Node{
					CreateNumber("0"),
				}),
				CreateNumber("-1"),
			},
			Error: "error at line 1: expected Boolean, got Number",
		},
		{
			Args: []node.Node{
				CreateNumber("-1"),
				CreateNumber("-1"),
			},
			Error: "error at line 1: expected List, got Number",
		},
		{
			Args:  []node.Node{},
			Error: "error at line 1: incorrect number of arguments. expected 2, got 0",
		},
		{
			Args: []node.Node{
				CreateNumber("-1"),
			},
			Error: "error at line 1: incorrect number of arguments. expected 2, got 1",
		},
		{
			Args: []node.Node{
				CreateNumber("-1"),
				CreateNumber("-1"),
				CreateNumber("-1"),
			},
			Error: "error at line 1: incorrect number of arguments. expected 2, got 3",
		},
	}

	for i, test := range tests {
		ast := []node.Node{
			CreateFunctionCall(
				CreateIdentifier("unwrap"),
				test.Args,
			),
		}

		actualError := getEvaluatorError(t, ast)

		AssertErrorEqual(t, i, test.Error, actualError)
	}
}

func TestBuiltin_UnwrapAllErrors(t *testing.T) {
	/*
		Not testing other errors associated with "unwrap" because "unwrap_all" calls the method associated
		with "unwrap", so that would be duplicate testing.
	*/

	tests := []struct {
		Args  []node.Node
		Error string
	}{
		{
			Args:  []node.Node{},
			Error: "error at line 1: incorrect number of arguments. expected 2, got 0",
		},
		{
			Args: []node.Node{
				CreateNumber("-1"),
			},
			Error: "error at line 1: incorrect number of arguments. expected 2, got 1",
		},
		{
			Args: []node.Node{
				CreateNumber("-1"),
				CreateNumber("-1"),
				CreateNumber("-1"),
			},
			Error: "error at line 1: incorrect number of arguments. expected 2, got 3",
		},
	}

	for i, test := range tests {
		ast := []node.Node{
			CreateFunctionCall(
				CreateIdentifier("unwrap_all"),
				test.Args,
			),
		}

		actualError := getEvaluatorError(t, ast)

		AssertErrorEqual(t, i, test.Error, actualError)
	}
}

func TestBuiltin_SliceIndexErrors(t *testing.T) {

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

		actualError := getEvaluatorError(t, ast)
		expectedError := test.Error

		AssertErrorEqual(t, i, expectedError, actualError)
	}
}

func TestBuiltin_SliceInvalidTypeError(t *testing.T) {
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

	actualError := getEvaluatorError(t, ast)
	expectedError := "error at line 1: expected List, got Boolean (\"true\")"

	AssertErrorEqual(t, 0, expectedError, actualError)
}