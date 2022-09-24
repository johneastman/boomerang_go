package tests

import (
	"boomerang/node"
	"boomerang/tokens"
	"testing"
)

func TestNode_CreateNumber(t *testing.T) {
	actualNode := node.CreateNumber("1")
	expectedNode := node.Node{
		Type:  node.NUMBER,
		Value: "1",
	}

	AssertNodeEqual(t, expectedNode, actualNode)
}

func TestNode_CreateBoolean(t *testing.T) {

	booleanLiterals := []string{
		"true",
		"false",
	}

	for _, booleanLiteral := range booleanLiterals {
		actualNode := node.CreateBoolean(booleanLiteral)
		expectedNode := node.Node{
			Type:  node.BOOLEAN,
			Value: booleanLiteral,
		}

		AssertNodeEqual(t, expectedNode, actualNode)
	}
}

func TestNode_CreateString(t *testing.T) {
	actualNode := node.CreateString("hello, world!", []node.Node{})
	expectedNode := node.Node{
		Type:   node.STRING,
		Value:  "hello, world!",
		Params: []node.Node{},
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
	actualNode := node.CreateIdentifier("variable")
	expectedNode := node.Node{
		Type:  node.IDENTIFIER,
		Value: "variable",
	}

	AssertNodeEqual(t, expectedNode, actualNode)
}

func TestNode_CreatePrintStatement(t *testing.T) {
	token := tokens.MINUS_TOKEN
	actualNode := node.CreatePrintStatement([]node.Node{
		node.CreateNumber("1"),
		node.CreateIdentifier("variable"),
		node.CreateBinaryExpression(
			node.CreateNumber("3"),
			token,
			node.CreateNumber("4"),
		),
	})
	expectedNode := node.Node{
		Type: node.PRINT_STMT,
		Params: []node.Node{
			{Type: node.NUMBER, Value: "1"},
			{Type: node.IDENTIFIER, Value: "variable"},
			{Type: node.BIN_EXPR, Params: []node.Node{
				{Type: node.NUMBER, Value: "3"},
				{Type: token.Type, Value: token.Literal},
				{Type: node.NUMBER, Value: "4"},
			}},
		},
	}

	AssertNodeEqual(t, expectedNode, actualNode)
}

func TestNode_CreateUnaryExpression(t *testing.T) {
	token := tokens.MINUS_TOKEN
	actualNode := node.CreateUnaryExpression(
		token,
		node.CreateNumber("36"),
	)
	expectedNode := node.Node{
		Type: node.UNARY_EXPR,
		Params: []node.Node{
			{Type: token.Type, Value: token.Literal},
			{Type: node.NUMBER, Value: "36"},
		},
	}

	AssertNodeEqual(t, expectedNode, actualNode)
}

func TestNode_CreateBinaryExpression(t *testing.T) {
	token := tokens.MINUS_TOKEN
	actualNode := node.CreateBinaryExpression(
		node.CreateNumber("44"),
		token,
		node.CreateNumber("36"),
	)
	expectedNode := node.Node{
		Type: node.BIN_EXPR,
		Params: []node.Node{
			{Type: node.NUMBER, Value: "44"},
			{Type: token.Type, Value: token.Literal},
			{Type: node.NUMBER, Value: "36"},
		},
	}

	AssertNodeEqual(t, expectedNode, actualNode)
}

func TestNode_CreateAssignmentStatement(t *testing.T) {
	actualNode := node.CreateAssignmentStatement(
		"my_number",
		node.CreateNumber("789"),
	)
	expectedNode := node.Node{
		Type: node.ASSIGN_STMT,
		Params: []node.Node{
			{Type: node.IDENTIFIER, Value: "my_number"},
			{Type: node.NUMBER, Value: "789"},
		},
	}

	AssertNodeEqual(t, expectedNode, actualNode)
}

func TestNode_CreateFunction(t *testing.T) {

	multiplyToken := tokens.ASTERISK_TOKEN
	divideToken := tokens.FORWARD_SLASH_TOKEN

	actualNode := node.CreateFunction(
		[]node.Node{
			node.CreateIdentifier("x"),
			node.CreateIdentifier("y"),
			node.CreateIdentifier("z"),
		},
		[]node.Node{
			node.CreateBinaryExpression(
				node.CreateIdentifier("x"),
				multiplyToken,
				node.CreateBinaryExpression(
					node.CreateIdentifier("y"),
					divideToken,
					node.CreateIdentifier("z"),
				),
			),
		},
	)
	expectedNode := node.Node{
		Type: node.FUNCTION,
		Params: []node.Node{
			{Type: node.PARAMETER, Params: []node.Node{
				{Type: node.IDENTIFIER, Value: "x"},
				{Type: node.IDENTIFIER, Value: "y"},
				{Type: node.IDENTIFIER, Value: "z"},
			}},
			{Type: node.STMTS, Params: []node.Node{
				{Type: node.BIN_EXPR, Params: []node.Node{
					{Type: node.IDENTIFIER, Value: "x"},
					{Type: multiplyToken.Type, Value: multiplyToken.Literal},
					{Type: node.BIN_EXPR, Params: []node.Node{
						{Type: node.IDENTIFIER, Value: "y"},
						{Type: divideToken.Type, Value: divideToken.Literal},
						{Type: node.IDENTIFIER, Value: "z"},
					}},
				}},
			}},
		},
	}

	AssertNodeEqual(t, expectedNode, actualNode)
}

func TestNode_CreateFunctionCall(t *testing.T) {
	token := tokens.MINUS_TOKEN
	actualNode := node.CreateFunctionCall(
		node.CreateIdentifier("functionCall"),
		[]node.Node{
			node.CreateNumber("1"),
			node.CreateIdentifier("variable"),
			node.CreateBinaryExpression(
				node.CreateNumber("3"),
				token,
				node.CreateNumber("4"),
			),
		})

	expectedNode := node.Node{
		Type: node.FUNCTION_CALL,
		Params: []node.Node{
			{Type: node.CALL_PARAMS, Params: []node.Node{
				{Type: node.NUMBER, Value: "1"},
				{Type: node.IDENTIFIER, Value: "variable"},
				{Type: node.BIN_EXPR, Params: []node.Node{
					{Type: node.NUMBER, Value: "3"},
					{Type: token.Type, Value: token.Literal},
					{Type: node.NUMBER, Value: "4"},
				}},
			}},
			{Type: node.IDENTIFIER, Value: "functionCall"},
		},
	}

	AssertNodeEqual(t, expectedNode, actualNode)
}

func TestNode_CreateParameters(t *testing.T) {
	token := tokens.ASTERISK_TOKEN
	actualNode := node.CreateParameters([]node.Node{
		node.CreateNumber("1"),
		node.CreateIdentifier("variable"),
		node.CreateBinaryExpression(
			node.CreateNumber("10"),
			token,
			node.CreateNumber("20"),
		),
	})

	expectedNode := node.Node{
		Type: node.PARAMETER,
		Params: []node.Node{
			{Type: node.NUMBER, Value: "1"},
			{Type: node.IDENTIFIER, Value: "variable"},
			{Type: node.BIN_EXPR, Params: []node.Node{
				{Type: node.NUMBER, Value: "10"},
				{Type: token.Type, Value: token.Literal},
				{Type: node.NUMBER, Value: "20"},
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
			Node:   node.CreateNumber("1"),
			String: "1",
		},
		{
			Node:   node.CreateNumber("2"),
			String: "2",
		},
		{
			Node: node.CreateParameters([]node.Node{
				node.CreateNumber("1"),
				node.CreateNumber("2"),
				node.CreateNumber("3"),
			}),
			String: "(1, 2, 3)",
		},
		{
			Node: node.CreateParameters([]node.Node{
				node.CreateNumber("1"),
				node.CreateNumber("2"),
				node.CreateNumber("3"),
				node.CreateParameters([]node.Node{
					node.CreateNumber("5"),
					node.CreateNumber("7"),
					node.CreateNumber("6"),
				}),
				node.CreateNumber("4"),
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
	actualNode := node.CreateReturnStatement(node.CreateBinaryExpression(
		node.CreateNumber("1"),
		tokens.FORWARD_SLASH_TOKEN,
		node.CreateNumber("2"),
	))
	expectedNode := node.Node{Type: node.RETURN, Params: []node.Node{
		{
			Type: node.BIN_EXPR,
			Params: []node.Node{
				{Type: node.NUMBER, Value: "1"},
				{Type: tokens.FORWARD_SLASH_TOKEN.Type, Value: tokens.FORWARD_SLASH_TOKEN.Literal},
				{Type: node.NUMBER, Value: "2"},
			}},
	}}

	AssertNodeEqual(t, expectedNode, actualNode)
}

func TestNode_BuiltinFunction(t *testing.T) {
	actualNode := node.CreateBuiltinFunction(node.BUILTIN_LEN)
	expectedNode := node.Node{Type: node.BUILTIN_FUNC, Value: node.BUILTIN_LEN}
	AssertNodeEqual(t, expectedNode, actualNode)
}
