package tests

import (
	"boomerang/node"
	"boomerang/tokens"
	"fmt"
	"testing"
)

func TestNode_Equals(t *testing.T) {
	tests := []struct {
		First   node.Node
		Second  node.Node
		IsEqual bool
	}{
		{
			First:   CreateNumber("5"),
			Second:  CreateNumber("5"),
			IsEqual: true,
		},
		{
			First:   CreateNumber("5"),
			Second:  CreateNumber("6"),
			IsEqual: false,
		},
		{
			First:   CreateBooleanTrue(),
			Second:  CreateRawString("true"),
			IsEqual: false,
		},
		{
			First:   CreateList([]node.Node{CreateNumber("1"), CreateNumber("2")}),
			Second:  CreateList([]node.Node{CreateNumber("1"), CreateNumber("2")}),
			IsEqual: true,
		},
		{
			First:   CreateList([]node.Node{CreateNumber("2"), CreateNumber("4"), CreateNumber("6")}),
			Second:  CreateList([]node.Node{CreateNumber("1"), CreateNumber("2")}),
			IsEqual: false,
		},
		{
			First: CreateFunction(
				[]node.Node{
					CreateIdentifier("a"),
					CreateIdentifier("b"),
					CreateIdentifier("c"),
				},
				[]node.Node{
					node.CreateBinaryExpression(
						node.CreateBinaryExpression(
							CreateIdentifier("a"),
							CreateTokenFromToken(tokens.PLUS_TOKEN),
							CreateIdentifier("b"),
						),
						CreateTokenFromToken(tokens.PLUS_TOKEN),
						CreateIdentifier("c"),
					),
				},
			),
			Second: CreateFunction(
				[]node.Node{
					CreateIdentifier("a"),
					CreateIdentifier("b"),
					CreateIdentifier("c"),
				},
				[]node.Node{
					node.CreateBinaryExpression(
						node.CreateBinaryExpression(
							CreateIdentifier("a"),
							CreateTokenFromToken(tokens.PLUS_TOKEN),
							CreateIdentifier("b"),
						),
						CreateTokenFromToken(tokens.PLUS_TOKEN),
						CreateIdentifier("c"),
					),
				},
			),
			IsEqual: true,
		},
	}

	for i, test := range tests {
		isEqual := test.First.Equals(test.Second)
		if isEqual != test.IsEqual {
			t.Fatalf("Test #%d - Nodes not equal. Left: %s, Right: %s", i, &test.First, &test.Second)
		}
	}
}

func TestNode_CreateNumber(t *testing.T) {
	actualNode := node.CreateNumber(20, "1")
	expectedNode := node.Node{
		Type:    node.NUMBER,
		Value:   "1",
		LineNum: 20,
	}

	AssertNodeEqual(t, 0, expectedNode, actualNode)
}

func TestNode_CreateBoolean(t *testing.T) {

	booleanLiterals := []string{
		"true",
		"false",
	}

	for _, booleanLiteral := range booleanLiterals {
		actualNode := node.CreateBoolean(TEST_LINE_NUM, booleanLiteral)
		expectedNode := node.Node{
			Type:    node.BOOLEAN,
			Value:   booleanLiteral,
			LineNum: 1,
		}

		AssertNodeEqual(t, 0, expectedNode, actualNode)
	}
}

func TestNode_CreateBooleanTrue(t *testing.T) {
	actualNode := node.CreateBooleanTrue(1)
	expectedNode := node.Node{
		Type:    node.BOOLEAN,
		Value:   tokens.TRUE_TOKEN.Literal,
		LineNum: 1,
	}

	AssertNodeEqual(t, 0, expectedNode, actualNode)
}

func TestNode_CreateBooleanFalse(t *testing.T) {
	actualNode := node.CreateBooleanFalse(1)
	expectedNode := node.Node{
		Type:    node.BOOLEAN,
		Value:   tokens.FALSE_TOKEN.Literal,
		LineNum: 1,
	}

	AssertNodeEqual(t, 0, expectedNode, actualNode)
}

func TestNode_CreateString(t *testing.T) {
	actualNode := node.CreateRawString(1, "hello, world!")
	expectedNode := node.Node{
		Type:    node.STRING,
		Value:   "hello, world!",
		LineNum: 1,
		Params:  []node.Node{},
	}
	AssertNodeEqual(t, 0, expectedNode, actualNode)
}

func TestNode_CreateTokenNode(t *testing.T) {
	token := tokens.PLUS_TOKEN
	actualNode := node.CreateTokenNode(token)
	expectedNode := node.Node{
		Type:  token.Type,
		Value: token.Literal,
	}

	AssertNodeEqual(t, 0, expectedNode, actualNode)
}

func TestNode_CreateIdentifier(t *testing.T) {
	actualNode := node.CreateIdentifier(1, "variable")
	expectedNode := node.Node{
		Type:    node.IDENTIFIER,
		Value:   "variable",
		LineNum: 1,
	}

	AssertNodeEqual(t, 0, expectedNode, actualNode)
}

func TestNode_CreateUnaryExpression(t *testing.T) {
	token := tokens.MINUS_TOKEN
	token.LineNumber = TEST_LINE_NUM

	actualNode := node.CreateUnaryExpression(
		token,
		node.CreateNumber(TEST_LINE_NUM, "36"),
	)
	expectedNode := node.Node{
		Type:    node.UNARY_EXPR,
		LineNum: TEST_LINE_NUM,
		Params: []node.Node{
			{Type: token.Type, Value: token.Literal, LineNum: TEST_LINE_NUM},
			{Type: node.NUMBER, Value: "36", LineNum: TEST_LINE_NUM},
		},
	}

	AssertNodeEqual(t, 0, expectedNode, actualNode)
}

func TestNode_CreateBinaryExpression(t *testing.T) {
	token := tokens.MINUS_TOKEN
	token.LineNumber = TEST_LINE_NUM

	actualNode := node.CreateBinaryExpression(
		node.CreateNumber(TEST_LINE_NUM, "44"),
		token,
		node.CreateNumber(TEST_LINE_NUM, "36"),
	)
	expectedNode := node.Node{
		Type:    node.BIN_EXPR,
		LineNum: TEST_LINE_NUM,
		Params: []node.Node{
			{Type: node.NUMBER, Value: "44", LineNum: TEST_LINE_NUM},
			{Type: token.Type, Value: token.Literal, LineNum: TEST_LINE_NUM},
			{Type: node.NUMBER, Value: "36", LineNum: TEST_LINE_NUM},
		},
	}

	AssertNodeEqual(t, 0, expectedNode, actualNode)
}

func TestNode_AssignmentNode(t *testing.T) {
	actualNode := node.CreateAssignmentNode(
		node.CreateIdentifier(20, "my_number"),
		node.CreateNumber(20, "789"),
	)
	expectedNode := node.Node{
		Type:    node.ASSIGN_STMT,
		LineNum: 20,
		Params: []node.Node{
			{Type: node.IDENTIFIER, LineNum: 20, Value: "my_number"},
			{Type: node.NUMBER, Value: "789", LineNum: 20},
		},
	}

	AssertNodeEqual(t, 0, expectedNode, actualNode)
}

func TestNode_CreateFunction(t *testing.T) {

	multiplyToken := tokens.ASTERISK_TOKEN
	multiplyToken.LineNumber = TEST_LINE_NUM

	divideToken := tokens.FORWARD_SLASH_TOKEN
	divideToken.LineNumber = TEST_LINE_NUM

	actualNode := node.CreateFunction(
		TEST_LINE_NUM,
		[]node.Node{
			node.CreateIdentifier(TEST_LINE_NUM, "x"),
			node.CreateIdentifier(TEST_LINE_NUM, "y"),
			node.CreateIdentifier(TEST_LINE_NUM, "z"),
		},
		node.CreateBlockStatements(
			[]node.Node{
				node.CreateBinaryExpression(
					node.CreateIdentifier(TEST_LINE_NUM, "x"),
					multiplyToken,
					node.CreateBinaryExpression(
						node.CreateIdentifier(TEST_LINE_NUM, "y"),
						divideToken,
						node.CreateIdentifier(TEST_LINE_NUM, "z"),
					),
				),
			},
		),
	)
	expectedNode := node.Node{
		Type:    node.FUNCTION,
		LineNum: TEST_LINE_NUM,
		Params: []node.Node{
			{Type: node.LIST, LineNum: TEST_LINE_NUM, Params: []node.Node{
				{Type: node.IDENTIFIER, Value: "x", LineNum: TEST_LINE_NUM},
				{Type: node.IDENTIFIER, Value: "y", LineNum: TEST_LINE_NUM},
				{Type: node.IDENTIFIER, Value: "z", LineNum: TEST_LINE_NUM},
			}},
			{Type: node.BLOCK_STATEMENTS, Params: []node.Node{
				{Type: node.BIN_EXPR, LineNum: TEST_LINE_NUM, Params: []node.Node{
					{Type: node.IDENTIFIER, Value: "x", LineNum: TEST_LINE_NUM},
					{Type: multiplyToken.Type, Value: multiplyToken.Literal, LineNum: multiplyToken.LineNumber},
					{Type: node.BIN_EXPR, LineNum: TEST_LINE_NUM, Params: []node.Node{
						{Type: node.IDENTIFIER, Value: "y", LineNum: TEST_LINE_NUM},
						{Type: divideToken.Type, Value: divideToken.Literal, LineNum: divideToken.LineNumber},
						{Type: node.IDENTIFIER, Value: "z", LineNum: TEST_LINE_NUM},
					}},
				}},
			}},
		},
	}

	AssertNodeEqual(t, 0, expectedNode, actualNode)
}

func TestNode_CreateFunctionCall(t *testing.T) {
	token := tokens.MINUS_TOKEN
	token.LineNumber = TEST_LINE_NUM

	actualNode := node.CreateFunctionCall(TEST_LINE_NUM,
		node.CreateIdentifier(TEST_LINE_NUM, "functionCall"),
		[]node.Node{
			node.CreateNumber(TEST_LINE_NUM, "1"),
			node.CreateIdentifier(TEST_LINE_NUM, "variable"),
			node.CreateBinaryExpression(
				node.CreateNumber(TEST_LINE_NUM, "3"),
				token,
				node.CreateNumber(TEST_LINE_NUM, "4"),
			),
		},
	)

	expectedNode := node.Node{
		Type:    node.FUNCTION_CALL,
		LineNum: TEST_LINE_NUM,
		Params: []node.Node{
			{Type: node.CALL_PARAMS, LineNum: TEST_LINE_NUM, Params: []node.Node{
				{Type: node.NUMBER, Value: "1", LineNum: TEST_LINE_NUM},
				{Type: node.IDENTIFIER, Value: "variable", LineNum: TEST_LINE_NUM},
				{Type: node.BIN_EXPR, LineNum: TEST_LINE_NUM, Params: []node.Node{
					{Type: node.NUMBER, Value: "3", LineNum: TEST_LINE_NUM},
					{Type: token.Type, Value: token.Literal, LineNum: token.LineNumber},
					{Type: node.NUMBER, Value: "4", LineNum: TEST_LINE_NUM},
				}},
			}},
			{Type: node.IDENTIFIER, Value: "functionCall", LineNum: TEST_LINE_NUM},
		},
	}

	AssertNodeEqual(t, 0, expectedNode, actualNode)
}

func TestNode_CreateParameters(t *testing.T) {
	token := tokens.ASTERISK_TOKEN
	token.LineNumber = TEST_LINE_NUM

	actualNode := node.CreateList(TEST_LINE_NUM, []node.Node{
		node.CreateNumber(TEST_LINE_NUM, "1"),
		node.CreateIdentifier(TEST_LINE_NUM, "variable"),
		node.CreateBinaryExpression(
			node.CreateNumber(TEST_LINE_NUM, "10"),
			token,
			node.CreateNumber(TEST_LINE_NUM, "20"),
		),
	})

	expectedNode := node.Node{
		Type:    node.LIST,
		LineNum: TEST_LINE_NUM,
		Params: []node.Node{
			{Type: node.NUMBER, Value: "1", LineNum: TEST_LINE_NUM},
			{Type: node.IDENTIFIER, Value: "variable", LineNum: TEST_LINE_NUM},
			{Type: node.BIN_EXPR, LineNum: TEST_LINE_NUM, Params: []node.Node{
				{Type: node.NUMBER, Value: "10", LineNum: TEST_LINE_NUM},
				{Type: token.Type, Value: token.Literal, LineNum: token.LineNumber},
				{Type: node.NUMBER, Value: "20", LineNum: TEST_LINE_NUM},
			}},
		},
	}

	AssertNodeEqual(t, 0, expectedNode, actualNode)
}

func TestNode_String(t *testing.T) {

	type StringTest struct {
		Node   node.Node
		String string
	}

	tests := []StringTest{
		{
			Node:   CreateNumber("1"),
			String: "1",
		},
		{
			Node:   CreateNumber("2"),
			String: "2",
		},
		{
			Node: CreateList([]node.Node{
				CreateNumber("1"),
				CreateNumber("2"),
				CreateNumber("3"),
			}),
			String: "(1, 2, 3)",
		},
		{
			Node: CreateList([]node.Node{
				CreateNumber("1"),
				CreateNumber("2"),
				CreateNumber("3"),
				CreateList([]node.Node{
					CreateNumber("5"),
					CreateNumber("7"),
					CreateNumber("6"),
				}),
				CreateNumber("4"),
			}),
			String: "(1, 2, 3, (5, 7, 6), 4)",
		},
		{
			Node: CreateFunction(
				[]node.Node{
					CreateIdentifier("a"),
					CreateIdentifier("b"),
					CreateIdentifier("c"),
				},
				[]node.Node{},
			),
			String: "func(a,b,c){...}",
		},
		{
			Node:   CreateBuiltinFunctionIdentifier("print"),
			String: "<built-in function print>",
		},
	}

	for i, test := range tests {
		testName := fmt.Sprintf("Test #%d", i)
		t.Run(testName, func(t *testing.T) {
			actualNodeString := test.Node.String()
			expectedNodeString := test.String

			if expectedNodeString != actualNodeString {
				t.Fatalf("Expected string: %#v, Actual string: %#v", expectedNodeString, actualNodeString)
			}
		})
	}
}

func TestNode_CreateBlockStatementReturnValueNoParams(t *testing.T) {
	actualNode := node.CreateBlockStatementReturnValue(TEST_LINE_NUM, nil)
	expectedNode := node.Node{
		Type:    node.MONAD,
		LineNum: TEST_LINE_NUM,
		Params:  []node.Node{},
	}
	AssertNodeEqual(t, 0, expectedNode, actualNode)
}

func TestNode_CreateBlockStatementReturnValueParams(t *testing.T) {

	actualReturnValue := node.CreateNumber(TEST_LINE_NUM, "5")
	actualNode := node.CreateBlockStatementReturnValue(TEST_LINE_NUM, &actualReturnValue)
	expectedNode := node.Node{
		Type:    node.MONAD,
		LineNum: TEST_LINE_NUM,
		Params: []node.Node{
			{
				Type:    node.NUMBER,
				Value:   "5",
				LineNum: TEST_LINE_NUM,
			},
		},
	}
	AssertNodeEqual(t, 0, expectedNode, actualNode)
}

func TestNode_CreateWhenNode(t *testing.T) {
	actualNode := node.CreateWhenNode(
		TEST_LINE_NUM,
		node.CreateNumber(TEST_LINE_NUM, "3"),
		[]node.Node{
			node.CreateCaseNode(
				TEST_LINE_NUM,
				node.CreateNumber(TEST_LINE_NUM, "1"),
				node.CreateBlockStatements([]node.Node{
					CreateNumber("1"),
				}),
			),
			node.CreateCaseNode(
				TEST_LINE_NUM,
				node.CreateNumber(TEST_LINE_NUM, "2"),
				node.CreateBlockStatements([]node.Node{
					CreateNumber("2"),
				}),
			),
		},
		node.CreateBlockStatements([]node.Node{
			CreateNumber("0"),
		}),
	)

	expectedNode := node.Node{
		Type:    node.WHEN,
		LineNum: TEST_LINE_NUM,
		Params: []node.Node{
			{Type: node.NUMBER, LineNum: TEST_LINE_NUM, Value: "3"},
			{Type: node.WHEN_CASES, LineNum: TEST_LINE_NUM, Params: []node.Node{
				{Type: node.CASE, LineNum: TEST_LINE_NUM, Params: []node.Node{
					{Type: node.NUMBER, LineNum: TEST_LINE_NUM, Value: "1"},
					{Type: node.BLOCK_STATEMENTS, Params: []node.Node{
						{Type: node.NUMBER, LineNum: TEST_LINE_NUM, Value: "1"},
					}},
				}},
				{Type: node.CASE, LineNum: TEST_LINE_NUM, Params: []node.Node{
					{Type: node.NUMBER, LineNum: TEST_LINE_NUM, Value: "2"},
					{Type: node.BLOCK_STATEMENTS, Params: []node.Node{
						{Type: node.NUMBER, LineNum: TEST_LINE_NUM, Value: "2"},
					}},
				}},
			}},
			{Type: node.BLOCK_STATEMENTS, Params: []node.Node{
				{Type: node.NUMBER, LineNum: TEST_LINE_NUM, Value: "0"},
			}},
		},
	}

	AssertNodeEqual(t, 0, expectedNode, actualNode)
}

func TestNode_CreateCaseNode(t *testing.T) {
	actualNode := node.CreateCaseNode(
		TEST_LINE_NUM,
		node.CreateNumber(TEST_LINE_NUM, "1"),
		node.CreateBlockStatements([]node.Node{
			node.CreateUnaryExpression(
				CreateTokenFromToken(tokens.MINUS_TOKEN),
				node.CreateNumber(TEST_LINE_NUM, "1"),
			),
		}),
	)

	expectedNode := node.Node{
		Type:    node.CASE,
		LineNum: TEST_LINE_NUM,
		Params: []node.Node{
			{Type: node.NUMBER, LineNum: TEST_LINE_NUM, Value: "1"},
			{Type: node.BLOCK_STATEMENTS, Params: []node.Node{
				{Type: node.UNARY_EXPR, LineNum: TEST_LINE_NUM, Params: []node.Node{
					{Type: tokens.MINUS, LineNum: TEST_LINE_NUM, Value: tokens.MINUS_TOKEN.Literal},
					{Type: node.NUMBER, LineNum: TEST_LINE_NUM, Value: "1"},
				}},
			}},
		},
	}

	AssertNodeEqual(t, 0, expectedNode, actualNode)
}

func TestNode_ForLoop(t *testing.T) {
	actualNode := node.CreateForLoop(
		TEST_LINE_NUM,
		node.CreateIdentifier(TEST_LINE_NUM, "element"),
		node.CreateList(
			TEST_LINE_NUM,
			[]node.Node{
				node.CreateNumber(TEST_LINE_NUM, "1"),
				node.CreateNumber(TEST_LINE_NUM, "2"),
				node.CreateNumber(TEST_LINE_NUM, "3"),
			},
		),
		node.CreateBlockStatements([]node.Node{
			node.CreateIdentifier(TEST_LINE_NUM, "element"),
		}),
	)

	expectedNode := node.Node{
		Type:    node.FOR_LOOP,
		LineNum: TEST_LINE_NUM,
		Params: []node.Node{
			{Type: node.IDENTIFIER, LineNum: TEST_LINE_NUM, Value: "element"},
			{Type: node.LIST, LineNum: TEST_LINE_NUM, Params: []node.Node{
				{Type: node.NUMBER, LineNum: TEST_LINE_NUM, Value: "1"},
				{Type: node.NUMBER, LineNum: TEST_LINE_NUM, Value: "2"},
				{Type: node.NUMBER, LineNum: TEST_LINE_NUM, Value: "3"},
			}},
			{Type: node.BLOCK_STATEMENTS, Params: []node.Node{
				{Type: node.IDENTIFIER, LineNum: TEST_LINE_NUM, Value: "element"},
			}},
		},
	}
	AssertNodeEqual(t, 0, expectedNode, actualNode)
}

func TestNode_WhileLoop(t *testing.T) {
	actualNode := node.CreateWhileLoop(
		TEST_LINE_NUM,
		node.CreateBooleanTrue(TEST_LINE_NUM),
		node.CreateBlockStatements([]node.Node{
			node.CreateRawString(TEST_LINE_NUM, "hello, world!"),
		}),
	)
	expectedNode := node.Node{
		Type:    node.WHILE_LOOP,
		LineNum: TEST_LINE_NUM,
		Params: []node.Node{
			{Type: node.BOOLEAN, Value: "true", LineNum: TEST_LINE_NUM},
			{Type: node.BLOCK_STATEMENTS, Params: []node.Node{
				{Type: node.STRING, Value: "hello, world!", LineNum: TEST_LINE_NUM},
			}},
		},
	}
	AssertNodeEqual(t, 0, expectedNode, actualNode)
}

func TestNode_BreakStatement(t *testing.T) {
	actualNode := node.CreateBreakStatement(TEST_LINE_NUM)
	expectedNode := node.Node{Type: node.BREAK, LineNum: TEST_LINE_NUM}
	AssertNodeEqual(t, 0, expectedNode, actualNode)
}

func TestNode_ContinueStatement(t *testing.T) {
	actualNode := node.CreateContinueStatement(TEST_LINE_NUM)
	expectedNode := node.Node{Type: node.CONTINUE, LineNum: TEST_LINE_NUM}
	AssertNodeEqual(t, 0, expectedNode, actualNode)
}

func TestNode_MonadWithValue(t *testing.T) {
	actualNode := node.CreateMonad(
		TEST_LINE_NUM,
		node.CreateNumber(TEST_LINE_NUM, "5").Ptr(),
	)
	expectedNode := node.Node{
		Type:    node.MONAD,
		LineNum: TEST_LINE_NUM,
		Params: []node.Node{
			{Type: node.NUMBER, Value: "5", LineNum: TEST_LINE_NUM},
		},
	}
	AssertNodeEqual(t, 0, expectedNode, actualNode)
}

func TestNode_MonadWithoutValue(t *testing.T) {
	actualNode := node.CreateMonad(
		TEST_LINE_NUM,
		nil,
	)
	expectedNode := node.Node{
		Type:    node.MONAD,
		LineNum: TEST_LINE_NUM,
		Params:  []node.Node{},
	}
	AssertNodeEqual(t, 0, expectedNode, actualNode)
}

func TestNode_TestCreateBuiltinFunctionIdentifier(t *testing.T) {
	actualNode := node.CreateBuiltinFunctionIdentifier(TEST_LINE_NUM, "function")
	expectedNode := node.Node{Type: node.BUILTIN_FUNCTION, Value: "function", LineNum: TEST_LINE_NUM}
	AssertNodeEqual(t, 0, expectedNode, actualNode)
}

func TestNode_TestCreateBuiltinVariableIdentifier(t *testing.T) {
	actualNode := node.CreateBuiltinVariableIdentifier(TEST_LINE_NUM, "variable")
	expectedNode := node.Node{Type: node.BUILTIN_VARIABLE, Value: "variable", LineNum: TEST_LINE_NUM}
	AssertNodeEqual(t, 0, expectedNode, actualNode)
}
