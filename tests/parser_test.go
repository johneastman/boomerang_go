package tests

import (
	"boomerang/node"
	"boomerang/tokens"
	"fmt"
	"testing"
)

func TestParser_Numbers(t *testing.T) {
	numbers := []string{
		"10",
		"1001",
		"5.5",
		"3.14159",
	}

	for i, number := range numbers {
		actualAST := getParserAST(fmt.Sprintf("%s;", number))
		expectedAST := []node.Node{
			CreateNumber(number),
		}

		AssertNodesEqual(t, i, expectedAST, actualAST)
	}
}

func TestParser_Booleans(t *testing.T) {
	tests := []struct {
		Source       string
		ExpectedNode node.Node
	}{
		{
			"true",
			CreateBooleanTrue(),
		},
		{
			"false",
			CreateBooleanFalse(),
		},
	}

	for i, test := range tests {
		actualAST := getParserAST(fmt.Sprintf("%s;", test.Source))
		expectedAST := []node.Node{
			test.ExpectedNode,
		}

		AssertNodesEqual(t, i, expectedAST, actualAST)
	}
}

func TestParser_Strings(t *testing.T) {

	plusToken := tokens.PLUS_TOKEN
	plusToken.LineNumber = 1

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
			InputSource:  "My age is {55}",
			OutputSource: "My age is <0>",
			Params: []node.Node{
				CreateNumber("55"),
			},
		},
		{
			InputSource:  "{1 + 1} {1 + 1} {1 + 1}",
			OutputSource: "<0> <1> <2>",
			Params: []node.Node{
				node.CreateBinaryExpression(
					CreateNumber("1"),
					plusToken,
					CreateNumber("1"),
				),
				node.CreateBinaryExpression(
					CreateNumber("1"),
					plusToken,
					CreateNumber("1"),
				),
				node.CreateBinaryExpression(
					CreateNumber("1"),
					plusToken,
					CreateNumber("1"),
				),
			},
		},
	}

	for i, test := range tests {
		actualAST := getParserAST(fmt.Sprintf("\"%s\";", test.InputSource))
		expectedAST := []node.Node{
			node.CreateString(1, test.OutputSource, test.Params),
		}
		AssertNodesEqual(t, i, expectedAST, actualAST)
	}
}

func TestParser_TestParameters(t *testing.T) {

	tests := []struct {
		Source string
		Params []node.Node
	}{
		// Empty list
		{
			Source: "()",
			Params: []node.Node{},
		},
		// 1-element list
		{
			Source: "(1,)",
			Params: []node.Node{
				CreateNumber("1"),
			},
		},
		// 2-element list
		{
			Source: "(1,2)",
			Params: []node.Node{
				CreateNumber("1"),
				CreateNumber("2"),
			},
		},
		// 3-element list
		{
			Source: "(1,2,3)",
			Params: []node.Node{
				CreateNumber("1"),
				CreateNumber("2"),
				CreateNumber("3"),
			},
		},
		{
			Source: "(1,2,(4, 5, 6),3)",
			Params: []node.Node{
				CreateNumber("1"),
				CreateNumber("2"),
				CreateList([]node.Node{
					CreateNumber("4"),
					CreateNumber("5"),
					CreateNumber("6"),
				}),
				CreateNumber("3"),
			},
		},
	}

	for i, test := range tests {
		actualAST := getParserAST(fmt.Sprintf("%s;", test.Source))
		expectedAST := []node.Node{
			CreateList(test.Params),
		}

		AssertNodesEqual(t, i, expectedAST, actualAST)
	}
}

func TestParser_NegativeNumber(t *testing.T) {
	actualAST := getParserAST("-66;")
	expectedAST := []node.Node{
		node.CreateUnaryExpression(
			CreateTokenFromToken(tokens.MINUS_TOKEN),
			CreateNumber("66"),
		),
	}

	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_Bang(t *testing.T) {
	actualAST := getParserAST("not true;")
	expectedAST := []node.Node{
		node.CreateUnaryExpression(
			CreateTokenFromToken(tokens.NOT_TOKEN),
			CreateBooleanTrue(),
		),
	}

	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_BinaryExpression(t *testing.T) {

	tests := []struct {
		Source      string
		ExpectedAST node.Node
	}{
		{
			"7 + 3",
			node.CreateBinaryExpression(
				CreateNumber("7"),
				CreateTokenFromToken(tokens.PLUS_TOKEN),
				CreateNumber("3"),
			),
		},
		{
			"13 % 4",
			node.CreateBinaryExpression(
				CreateNumber("13"),
				CreateTokenFromToken(tokens.MODULO_TOKEN),
				CreateNumber("4"),
			),
		},
		{
			"14 == 13",
			node.CreateBinaryExpression(
				CreateNumber("14"),
				CreateTokenFromToken(tokens.EQ_TOKEN),
				CreateNumber("13"),
			),
		},
		{
			"true != false",
			node.CreateBinaryExpression(
				CreateBooleanTrue(),
				CreateTokenFromToken(tokens.NE_TOKEN),
				CreateBooleanFalse(),
			),
		},
		{
			"true or false",
			node.CreateBinaryExpression(
				CreateBooleanTrue(),
				CreateTokenFromToken(tokens.OR_TOKEN),
				CreateBooleanFalse(),
			),
		},
		{
			"false and true",
			node.CreateBinaryExpression(
				CreateBooleanFalse(),
				CreateTokenFromToken(tokens.AND_TOKEN),
				CreateBooleanTrue(),
			),
		},
		{
			"5 in (1, 2, 3, 4, 5)",
			node.CreateBinaryExpression(
				CreateNumber("5"),
				CreateTokenFromToken(tokens.IN_TOKEN),
				CreateList([]node.Node{
					CreateNumber("1"),
					CreateNumber("2"),
					CreateNumber("3"),
					CreateNumber("4"),
					CreateNumber("5"),
				}),
			),
		},
	}

	for i, test := range tests {
		source := fmt.Sprintf("%s;", test.Source)
		actualAST := getParserAST(source)
		expectedAST := []node.Node{
			test.ExpectedAST,
		}
		AssertNodesEqual(t, i, expectedAST, actualAST)
	}
}

func TestParser_Parentheses(t *testing.T) {
	actualAST := getParserAST("7 + (3);")
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			CreateNumber("7"),
			CreateTokenFromToken(tokens.PLUS_TOKEN),
			CreateNumber("3"),
		),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_ParenthesesBinaryExpression(t *testing.T) {

	actualAST := getParserAST("7 + (5 - 2);")
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			CreateNumber("7"),
			CreateTokenFromToken(tokens.PLUS_TOKEN),
			node.CreateBinaryExpression(
				CreateNumber("5"),
				CreateTokenFromToken(tokens.MINUS_TOKEN),
				CreateNumber("2"),
			),
		),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_VariableAssignment(t *testing.T) {
	actualAST := getParserAST("variable = 8 / 2;")
	expectedAST := []node.Node{
		CreateAssignmentNode(
			CreateIdentifier("variable"),
			node.CreateBinaryExpression(
				CreateNumber("8"),
				CreateTokenFromToken(tokens.FORWARD_SLASH_TOKEN),
				CreateNumber("2"),
			),
		),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_Identifier(t *testing.T) {
	actualAST := getParserAST("variable;")
	expectedAST := []node.Node{
		CreateIdentifier("variable"),
	}

	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_Function(t *testing.T) {
	actualAST := getParserAST("func(a, b) { a + b; };")
	expectedAST := []node.Node{
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
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_FunctionNoParameters(t *testing.T) {
	actualAST := getParserAST("func() { 3 + 4; };")
	expectedAST := []node.Node{
		CreateFunction(
			[]node.Node{},
			[]node.Node{
				node.CreateBinaryExpression(
					CreateNumber("3"),
					CreateTokenFromToken(tokens.PLUS_TOKEN),
					CreateNumber("4"),
				),
			},
		),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_FunctionNoParametersNoStatements(t *testing.T) {
	actualAST := getParserAST("func() {};")
	expectedAST := []node.Node{
		CreateFunction(
			[]node.Node{},
			[]node.Node{},
		),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_FunctionKeywordArguments(t *testing.T) {
	actualAST := getParserAST("func(a, b=2, c=3) {};")
	expectedAST := []node.Node{
		CreateFunction(
			[]node.Node{
				CreateIdentifier("a"),
				CreateAssignmentNode(CreateIdentifier("b"), CreateNumber("2")),
				CreateAssignmentNode(CreateIdentifier("c"), CreateNumber("3")),
			},
			[]node.Node{},
		),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_FunctionCallWithNoParameters(t *testing.T) {
	actualAST := getParserAST("divide <- ();")
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			CreateIdentifier("divide"),
			CreateTokenFromToken(tokens.SEND_TOKEN),
			CreateList([]node.Node{}),
		),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_FunctionCallWithFunctionLiteralAndLeftSend(t *testing.T) {
	actualAST := getParserAST("func(c, d) { d - c; } <- (10, 2);")

	functionNode := CreateFunction(
		[]node.Node{
			CreateIdentifier("c"),
			CreateIdentifier("d"),
		},
		[]node.Node{
			node.CreateBinaryExpression(
				CreateIdentifier("d"),
				CreateTokenFromToken(tokens.MINUS_TOKEN),
				CreateIdentifier("c"),
			),
		},
	)

	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			functionNode,
			CreateTokenFromToken(tokens.SEND_TOKEN),
			CreateList([]node.Node{
				CreateNumber("10"),
				CreateNumber("2"),
			}),
		),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_FunctionCallWithIdentifierAndLeftSend(t *testing.T) {
	actualAST := getParserAST("multiply <- (10, 3);")
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			CreateIdentifier("multiply"),
			CreateTokenFromToken(tokens.SEND_TOKEN),
			CreateList([]node.Node{
				CreateNumber("10"),
				CreateNumber("3"),
			}),
		),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_ListIndex(t *testing.T) {
	actualAST := getParserAST("numbers @ 1;")
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			CreateIdentifier("numbers"),
			CreateTokenFromToken(tokens.AT_TOKEN),
			CreateNumber("1"),
		),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_FunctionCallPrecedenceExpression(t *testing.T) {
	actualAST := getParserAST("add <- (3, 4) + 3;")
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			node.CreateBinaryExpression(
				CreateIdentifier("add"),
				CreateTokenFromToken(tokens.SEND_TOKEN),
				CreateList([]node.Node{
					CreateNumber("3"),
					CreateNumber("4"),
				}),
			),
			CreateTokenFromToken(tokens.PLUS_TOKEN),
			CreateNumber("3"),
		),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_WhenExpression_NotBoolean(t *testing.T) {
	actualAST := getParserAST("when pos { is 0 { 5; } is 1 { 10; } else { 15; } };")
	expectedAST := []node.Node{
		CreateWhenNode(
			CreateIdentifier("pos"),
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
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_WhenExpression_BooleanTrueAndNoElse(t *testing.T) {
	actualAST := getParserAST("when { true { 5; } false { 10; } };")
	expectedAST := []node.Node{
		CreateWhenNode(
			CreateBooleanTrue(),
			[]node.Node{
				CreateWhenCaseNode(
					CreateBooleanTrue(),
					[]node.Node{
						CreateNumber("5"),
					},
				),
				CreateWhenCaseNode(
					CreateBooleanFalse(),
					[]node.Node{
						CreateNumber("10"),
					},
				),
			},
			[]node.Node{},
		),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_WhenExpression_BooleanFalseAndNoElse(t *testing.T) {
	actualAST := getParserAST("when not { true { 5; } false { 10; } };")
	expectedAST := []node.Node{
		CreateWhenNode(
			CreateBooleanFalse(),
			[]node.Node{
				CreateWhenCaseNode(
					CreateBooleanTrue(),
					[]node.Node{
						CreateNumber("5"),
					},
				),
				CreateWhenCaseNode(
					CreateBooleanFalse(),
					[]node.Node{
						CreateNumber("10"),
					},
				),
			},
			[]node.Node{},
		),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_WhileLoop(t *testing.T) {
	actualAST := getParserAST("while i < 10 { i = i + 1; };")
	expectedAST := []node.Node{
		CreateWhileLoop(
			node.CreateBinaryExpression(
				CreateIdentifier("i"),
				CreateTokenFromToken(tokens.LT_TOKEN),
				CreateNumber("10"),
			),
			[]node.Node{
				CreateAssignmentNode(
					CreateIdentifier("i"),
					node.CreateBinaryExpression(
						CreateIdentifier("i"),
						CreateTokenFromToken(tokens.PLUS_TOKEN),
						CreateNumber("1"),
					),
				),
			},
		),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_ForLoop(t *testing.T) {
	actualAST := getParserAST("for e in list { print <- (e,); };")
	expectedAST := []node.Node{
		CreateForLoop(
			CreateIdentifier("e"),
			CreateIdentifier("list"),
			[]node.Node{
				node.CreateBinaryExpression(
					CreateBuiltinFunctionIdentifier("print"),
					CreateTokenFromToken(tokens.SEND_TOKEN),
					CreateList([]node.Node{
						CreateIdentifier("e"),
					}),
				),
			},
		),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_BreakStatement(t *testing.T) {
	actualAST := getParserAST("break;")
	expectedAST := []node.Node{
		CreateBreakStatement(),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_ContinueStatement(t *testing.T) {
	actualAST := getParserAST("continue;")
	expectedAST := []node.Node{
		CreateContinueStatement(),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

/* * * * * * * *
 * ERROR TESTS *
 * * * * * * * */

func TestParser_UnexpectedTokenError(t *testing.T) {
	actualError := getParserError(t, "1")
	expectedError := "error at line 1: expected token type SEMICOLON (\";\"), got EOF (\"\")"

	AssertErrorEqual(t, 0, expectedError, actualError)
}

func TestParser_InvalidPrefixError(t *testing.T) {
	actualError := getParserError(t, "+;")
	expectedError := "error at line 1: invalid prefix: PLUS (\"+\")"

	AssertErrorEqual(t, 0, expectedError, actualError)
}

func TestParser_InvalidPrefixForGroupedExpressionError(t *testing.T) {
	actualError := getParserError(t, "(1];")
	expectedError := "error at line 1: expected CLOSED_PAREN (\")\") or COMMA (\",\"), got CLOSED_BRACKET (\"]\")"

	AssertErrorEqual(t, 0, expectedError, actualError)
}

func TestParser_WhenExpressionErrors(t *testing.T) {

	tests := []struct {
		Source string
		Error  string
	}{
		{
			Source: "when { is",
			Error:  "error at line 1: \"IS\" not allowed for boolean values",
		},
		{
			Source: "when num { 1",
			Error:  "error at line 1: expected token type IS (\"is\"), got NUMBER (\"1\")",
		},
	}

	for _, test := range tests {
		actualError := getParserError(t, test.Source)

		AssertErrorEqual(t, 0, test.Error, actualError)
	}
}
