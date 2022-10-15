package tests

import (
	"boomerang/node"
	"boomerang/parser"
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
		actualAST := getAST(fmt.Sprintf("%s;", number))
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
		actualAST := getAST(fmt.Sprintf("%s;", test.Source))
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
		actualAST := getAST(fmt.Sprintf("\"%s\";", test.InputSource))
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
		{
			Source: "()",
			Params: []node.Node{},
		},
		{
			Source: "(1,)",
			Params: []node.Node{
				CreateNumber("1"),
			},
		},
		{
			Source: "(1,2,)",
			Params: []node.Node{
				CreateNumber("1"),
				CreateNumber("2"),
			},
		},
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
		actualAST := getAST(fmt.Sprintf("%s;", test.Source))
		expectedAST := []node.Node{
			CreateList(test.Params),
		}

		AssertNodesEqual(t, i, expectedAST, actualAST)
	}
}

func TestParser_NegativeNumber(t *testing.T) {
	actualAST := getAST("-66;")
	expectedAST := []node.Node{
		node.CreateUnaryExpression(
			CreateTokenFromToken(tokens.MINUS_TOKEN),
			CreateNumber("66"),
		),
	}

	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_Bang(t *testing.T) {
	actualAST := getAST("not true;")
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
			"14 == 13",
			node.CreateBinaryExpression(
				CreateNumber("14"),
				CreateTokenFromToken(tokens.EQ_TOKEN),
				CreateNumber("13"),
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
	}

	for i, test := range tests {
		source := fmt.Sprintf("%s;", test.Source)
		actualAST := getAST(source)
		expectedAST := []node.Node{
			test.ExpectedAST,
		}
		AssertNodesEqual(t, i, expectedAST, actualAST)
	}
}

func TestParser_Parentheses(t *testing.T) {
	actualAST := getAST("7 + (3);")
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

	actualAST := getAST("7 + (5 - 2);")
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
	actualAST := getAST("variable = 8 / 2;")
	expectedAST := []node.Node{
		CreateAssignmentStatement(
			"variable",
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
	actualAST := getAST("variable;")
	expectedAST := []node.Node{
		CreateIdentifier("variable"),
	}

	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_PrintStatement(t *testing.T) {
	actualAST := getAST("print(1, 2, variable);")
	expectedAST := []node.Node{
		CreatePrintStatement(
			[]node.Node{
				CreateNumber("1"),
				CreateNumber("2"),
				CreateIdentifier("variable"),
			},
		),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_PrintStatementNoArguments(t *testing.T) {
	actualAST := getAST("print();")
	expectedAST := []node.Node{
		CreatePrintStatement([]node.Node{}),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_Function(t *testing.T) {
	actualAST := getAST("func(a, b) { a + b; };")
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
	actualAST := getAST("func() { 3 + 4; };")
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
	actualAST := getAST("func() {};")
	expectedAST := []node.Node{
		CreateFunction(
			[]node.Node{},
			[]node.Node{},
		),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_FunctionCallWithNoParameters(t *testing.T) {
	actualAST := getAST("divide <- ();")
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			CreateIdentifier("divide"),
			CreateTokenFromToken(tokens.PTR_TOKEN),
			CreateList([]node.Node{}),
		),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_FunctionCallWithFunctionLiteralAndLeftPointer(t *testing.T) {
	actualAST := getAST("func(c, d) { d - c; } <- (10, 2);")

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
			CreateTokenFromToken(tokens.PTR_TOKEN),
			CreateList([]node.Node{
				CreateNumber("10"),
				CreateNumber("2"),
			}),
		),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_FunctionCallWithIdentifierAndLeftPointer(t *testing.T) {
	actualAST := getAST("multiply <- (10, 3);")
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			CreateIdentifier("multiply"),
			CreateTokenFromToken(tokens.PTR_TOKEN),
			CreateList([]node.Node{
				CreateNumber("10"),
				CreateNumber("3"),
			}),
		),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_ListIndex(t *testing.T) {
	actualAST := getAST("numbers @ 1;")
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
	actualAST := getAST("add <- (3, 4) + 3;")
	expectedAST := []node.Node{
		node.CreateBinaryExpression(
			node.CreateBinaryExpression(
				CreateIdentifier("add"),
				CreateTokenFromToken(tokens.PTR_TOKEN),
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
	actualAST := getAST("when pos { is 0 { 5; } is 1 { 10; } else { 15; } };")
	expectedAST := []node.Node{
		CreateWhenNode(
			CreateIdentifier("pos"),
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
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_WhenExpression_BooleanTrueAndNoElse(t *testing.T) {
	actualAST := getAST("when { true { 5; } false { 10; } };")
	expectedAST := []node.Node{
		CreateWhenNode(
			CreateBooleanTrue(),
			[]node.Node{
				CreateWhenCaseNode(
					CreateBooleanTrue(),
					CreateBlockStatements([]node.Node{
						CreateNumber("5"),
					}),
				),
				CreateWhenCaseNode(
					CreateBooleanFalse(),
					CreateBlockStatements([]node.Node{
						CreateNumber("10"),
					}),
				),
			},
			CreateBlockStatements([]node.Node{}),
		),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_WhenExpression_BooleanFalseAndNoElse(t *testing.T) {
	actualAST := getAST("when not { true { 5; } false { 10; } };")
	expectedAST := []node.Node{
		CreateWhenNode(
			CreateBooleanFalse(),
			[]node.Node{
				CreateWhenCaseNode(
					CreateBooleanTrue(),
					CreateBlockStatements([]node.Node{
						CreateNumber("5"),
					}),
				),
				CreateWhenCaseNode(
					CreateBooleanFalse(),
					CreateBlockStatements([]node.Node{
						CreateNumber("10"),
					}),
				),
			},
			CreateBlockStatements([]node.Node{}),
		),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_WhileLoop(t *testing.T) {
	actualAST := getAST("while i < 10 { i = i + 1; };")
	expectedAST := []node.Node{
		CreateWhileLoop(
			node.CreateBinaryExpression(
				CreateIdentifier("i"),
				CreateTokenFromToken(tokens.LT_TOKEN),
				CreateNumber("10"),
			),
			CreateBlockStatements([]node.Node{
				CreateAssignmentStatement(
					"i",
					node.CreateBinaryExpression(
						CreateIdentifier("i"),
						CreateTokenFromToken(tokens.PLUS_TOKEN),
						CreateNumber("1"),
					),
				),
			}),
		),
	}
	AssertNodesEqual(t, 0, expectedAST, actualAST)
}

func TestParser_ForLoop(t *testing.T) {
	actualAST := getAST("for e in list { print(e); };")
	expectedAST := []node.Node{
		CreateForLoop(
			CreateIdentifier("e"),
			CreateIdentifier("list"),
			CreateBlockStatements([]node.Node{
				CreatePrintStatement([]node.Node{
					CreateIdentifier("e"),
				}),
			}),
		),
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

func getAST(source string) []node.Node {
	t := tokens.NewTokenizer(source)

	p, err := parser.NewParser(t)
	if err != nil {
		panic(err.Error())
	}

	ast, err := p.Parse()
	if err != nil {
		panic(err.Error())
	}
	return *ast
}

func getParserError(t *testing.T, source string) string {
	tokenizer := tokens.NewTokenizer(source)

	p, err := parser.NewParser(tokenizer)
	if err != nil {
		panic(err.Error())
	}

	_, err = p.Parse()
	if err == nil {
		t.Fatalf("Expected error to not be nil")
	}
	return err.Error()
}
