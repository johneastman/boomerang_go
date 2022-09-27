package tests

import (
	"boomerang/node"
	"boomerang/tokens"
	"testing"
)

func TestNode_CreateNumber(t *testing.T) {
	actualNode := node.CreateNumber(20, "1")
	expectedNode := node.Node{
		Type:    node.NUMBER,
		Value:   "1",
		LineNum: 20,
	}

	AssertNodeEqual(t, expectedNode, actualNode)
}

func TestNode_CreateBoolean(t *testing.T) {

	booleanLiterals := []string{
		"true",
		"false",
	}

	for _, booleanLiteral := range booleanLiterals {
		actualNode := node.CreateBoolean(booleanLiteral, 1)
		expectedNode := node.Node{
			Type:    node.BOOLEAN,
			Value:   booleanLiteral,
			LineNum: 1,
		}

		AssertNodeEqual(t, expectedNode, actualNode)
	}
}

func TestNode_CreateBooleanTrue(t *testing.T) {
	actualNode := node.CreateBooleanTrue(1)
	expectedNode := node.Node{
		Type:    node.BOOLEAN,
		Value:   tokens.TRUE_TOKEN.Literal,
		LineNum: 1,
	}

	AssertNodeEqual(t, expectedNode, actualNode)
}

func TestNode_CreateBooleanFalse(t *testing.T) {
	actualNode := node.CreateBooleanFalse(1)
	expectedNode := node.Node{
		Type:    node.BOOLEAN,
		Value:   tokens.FALSE_TOKEN.Literal,
		LineNum: 1,
	}

	AssertNodeEqual(t, expectedNode, actualNode)
}

func TestNode_CreateString(t *testing.T) {
	actualNode := node.CreateRawString(1, "hello, world!")
	expectedNode := node.Node{
		Type:    node.STRING,
		Value:   "hello, world!",
		LineNum: 1,
		Params:  []node.Node{},
	}
	AssertNodeEqual(t, expectedNode, actualNode)
}

func TestNode_CreateTokenNode(t *testing.T) {
	token := tokens.PLUS_TOKEN
	actualNode := node.CreateTokenNode(token)
	expectedNode := node.Node{
		Type:  token.Type,
		Value: token.Literal,
	}

	AssertNodeEqual(t, expectedNode, actualNode)
}

func TestNode_CreateIdentifier(t *testing.T) {
	actualNode := node.CreateIdentifier(1, "variable")
	expectedNode := node.Node{
		Type:    node.IDENTIFIER,
		Value:   "variable",
		LineNum: 1,
	}

	AssertNodeEqual(t, expectedNode, actualNode)
}

func TestNode_CreatePrintStatement(t *testing.T) {
	token := tokens.MINUS_TOKEN
	token.LineNumber = TEST_LINE_NUM

	actualNode := CreatePrintStatement([]node.Node{
		node.CreateNumber(TEST_LINE_NUM, "1"),
		node.CreateIdentifier(TEST_LINE_NUM, "variable"),
		node.CreateBinaryExpression(
			node.CreateNumber(TEST_LINE_NUM, "3"),
			token,
			node.CreateNumber(TEST_LINE_NUM, "4"),
		),
	})

	expectedNode := node.Node{
		Type:    node.PRINT_STMT,
		LineNum: TEST_LINE_NUM,
		Params: []node.Node{
			{Type: node.NUMBER, Value: "1", LineNum: TEST_LINE_NUM},
			{Type: node.IDENTIFIER, Value: "variable", LineNum: TEST_LINE_NUM},
			{Type: node.BIN_EXPR, LineNum: TEST_LINE_NUM, Params: []node.Node{
				{Type: node.NUMBER, Value: "3", LineNum: TEST_LINE_NUM},
				{Type: token.Type, Value: token.Literal, LineNum: TEST_LINE_NUM},
				{Type: node.NUMBER, Value: "4", LineNum: TEST_LINE_NUM},
			}},
		},
	}

	AssertNodeEqual(t, expectedNode, actualNode)
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

	AssertNodeEqual(t, expectedNode, actualNode)
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

	AssertNodeEqual(t, expectedNode, actualNode)
}

func TestNode_CreateAssignmentStatement(t *testing.T) {
	actualNode := node.CreateAssignmentStatement(
		"my_number",
		node.CreateNumber(20, "789"),
		10,
	)
	expectedNode := node.Node{
		Type:    node.ASSIGN_STMT,
		LineNum: 10,
		Params: []node.Node{
			{Type: node.IDENTIFIER, Value: "my_number"},
			{Type: node.NUMBER, Value: "789", LineNum: 20},
		},
	}

	AssertNodeEqual(t, expectedNode, actualNode)
}

func TestNode_CreateFunction(t *testing.T) {

	multiplyToken := tokens.ASTERISK_TOKEN
	multiplyToken.LineNumber = TEST_LINE_NUM

	divideToken := tokens.FORWARD_SLASH_TOKEN
	divideToken.LineNumber = TEST_LINE_NUM

	actualNode := node.CreateFunction(
		[]node.Node{
			node.CreateIdentifier(TEST_LINE_NUM, "x"),
			node.CreateIdentifier(TEST_LINE_NUM, "y"),
			node.CreateIdentifier(TEST_LINE_NUM, "z"),
		},
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
		TEST_LINE_NUM,
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
			{Type: node.STMTS, LineNum: TEST_LINE_NUM, Params: []node.Node{
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

	AssertNodeEqual(t, expectedNode, actualNode)
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

	AssertNodeEqual(t, expectedNode, actualNode)
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

	AssertNodeEqual(t, expectedNode, actualNode)
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
	}

	for _, test := range tests {
		actualNodeString := test.Node.String()
		expectedNodeString := test.String

		if expectedNodeString != actualNodeString {
			t.Fatalf("Expected string: %#v, Actual string: %#v", expectedNodeString, actualNodeString)
		}
	}
}

func TestNode_ReturnStatement(t *testing.T) {

	token := tokens.FORWARD_SLASH_TOKEN
	token.LineNumber = TEST_LINE_NUM

	actualNode := node.CreateReturnStatement(TEST_LINE_NUM,
		node.CreateBinaryExpression(
			node.CreateNumber(TEST_LINE_NUM, "1"),
			token,
			node.CreateNumber(TEST_LINE_NUM, "2"),
		),
	)
	expectedNode := node.Node{
		Type:    node.RETURN,
		LineNum: TEST_LINE_NUM,
		Params: []node.Node{
			{
				Type:    node.BIN_EXPR,
				LineNum: TEST_LINE_NUM,
				Params: []node.Node{
					{Type: node.NUMBER, Value: "1", LineNum: TEST_LINE_NUM},
					{Type: token.Type, Value: token.Literal, LineNum: token.LineNumber},
					{Type: node.NUMBER, Value: "2", LineNum: TEST_LINE_NUM},
				}},
		}}

	AssertNodeEqual(t, expectedNode, actualNode)
}

func TestNode_BuiltinFunction(t *testing.T) {
	actualNode := node.CreateBuiltinFunction(node.BUILTIN_LEN, 13)
	expectedNode := node.Node{Type: node.BUILTIN_FUNC, Value: node.BUILTIN_LEN, LineNum: 13}
	AssertNodeEqual(t, expectedNode, actualNode)
}

func TestNode_CreateReturnValueNoParams(t *testing.T) {
	actualNode := node.CreateFunctionReturnValue(TEST_LINE_NUM, nil)
	expectedNode := node.Node{
		Type:    node.LIST,
		LineNum: TEST_LINE_NUM,
		Params: []node.Node{
			{
				Type:    node.BOOLEAN,
				Value:   tokens.FALSE_TOKEN.Literal,
				LineNum: TEST_LINE_NUM,
			},
		},
	}
	AssertNodeEqual(t, expectedNode, actualNode)
}

func TestNode_CreateReturnValueParams(t *testing.T) {

	actualReturnValue := node.CreateNumber(TEST_LINE_NUM, "5")
	actualNode := node.CreateFunctionReturnValue(TEST_LINE_NUM, &actualReturnValue)
	expectedNode := node.Node{
		Type:    node.LIST,
		LineNum: TEST_LINE_NUM,
		Params: []node.Node{
			{
				Type:    node.BOOLEAN,
				Value:   tokens.TRUE_TOKEN.Literal,
				LineNum: TEST_LINE_NUM,
			},
			{
				Type:    node.NUMBER,
				Value:   "5",
				LineNum: TEST_LINE_NUM,
			},
		},
	}
	AssertNodeEqual(t, expectedNode, actualNode)
}

func TestNode_CreateIfStatement(t *testing.T) {
	actualNode := node.CreateIfStatement(TEST_LINE_NUM,
		node.CreateBooleanTrue(TEST_LINE_NUM),
		[]node.Node{
			node.CreatePrintStatement(TEST_LINE_NUM, []node.Node{
				node.CreateRawString(TEST_LINE_NUM, "true!!!"),
			}),
		},
	)
	expectedNode := node.Node{
		Type:    node.IF_STMT,
		LineNum: TEST_LINE_NUM,
		Params: []node.Node{
			{Type: node.BOOLEAN, Value: tokens.TRUE_TOKEN.Literal, LineNum: TEST_LINE_NUM},
			{Type: node.TRUE_BRANCH, LineNum: TEST_LINE_NUM, Params: []node.Node{
				{Type: node.PRINT_STMT, LineNum: TEST_LINE_NUM, Params: []node.Node{
					{Type: node.STRING, Value: "true!!!", LineNum: TEST_LINE_NUM},
				}},
			}},
		},
	}

	AssertNodeEqual(t, expectedNode, actualNode)
}
