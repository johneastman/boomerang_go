package tests

import (
	"boomerang/evaluator"
	"boomerang/node"
	"boomerang/parser"
	"boomerang/tokens"
	"boomerang/utils"
	"os"
	"strings"
	"testing"
)

const INTEGRATION_TESTS_DIRECTORY = "integration_tests"

func TestIntegration_EditList(t *testing.T) {

	actualResults := evaluateSource(t, "edit_list.bmg")
	expectedResult := []node.Node{

		node.CreateFunction(
			1,
			[]node.Node{
				node.CreateIdentifier(1, "list"),
				node.CreateIdentifier(1, "pos"),
				node.CreateIdentifier(1, "new_value"),
			},
			node.CreateBlockStatements(
				[]node.Node{

					// FUNCTION FIRST EXPRESSION
					node.CreateAssignmentNode(
						node.CreateIdentifier(2, "list_len"),
						node.CreateBinaryExpression(
							node.CreateBuiltinFunctionIdentifier(2, "len"),
							CreateTokenWithLineNum(tokens.SEND_TOKEN, 2),
							node.CreateList(2, []node.Node{
								node.CreateIdentifier(2, "list"),
							}),
						),
					),

					// FUNCTION SECOND EXPRESSION
					node.CreateAssignmentNode(
						node.CreateIdentifier(4, "new_numbers"),
						node.CreateWhenNode(4,
							node.CreateIdentifier(4, "pos"),
							[]node.Node{

								// FIRST WHEN CASE
								node.CreateCaseNode(
									5,
									node.CreateNumber(5, "0"),
									node.CreateBlockStatements(
										[]node.Node{
											node.CreateBinaryExpression(
												node.CreateList(6, []node.Node{node.CreateIdentifier(6, "new_value")}),
												CreateTokenWithLineNum(tokens.SEND_TOKEN, 6),
												node.CreateBinaryExpression(
													node.CreateBuiltinFunctionIdentifier(6, "slice"),
													CreateTokenWithLineNum(tokens.SEND_TOKEN, 6),
													node.CreateList(6, []node.Node{
														node.CreateIdentifier(6, "list"),
														node.CreateNumber(6, "1"),
														node.CreateBinaryExpression(
															node.CreateIdentifier(6, "list_len"),
															CreateTokenWithLineNum(tokens.MINUS_TOKEN, 6),
															node.CreateNumber(6, "1"),
														),
													}),
												),
											),
										},
									),
								),

								// SECOND WHEN CASE
								node.CreateCaseNode(
									8,
									node.CreateBinaryExpression(
										node.CreateIdentifier(8, "list_len"),
										CreateTokenWithLineNum(tokens.MINUS_TOKEN, 8),
										node.CreateNumber(8, "1"),
									),
									node.CreateBlockStatements(
										[]node.Node{
											node.CreateBinaryExpression(
												node.CreateBinaryExpression(
													node.CreateBuiltinFunctionIdentifier(9, "slice"),
													CreateTokenWithLineNum(tokens.SEND_TOKEN, 9),
													node.CreateList(9, []node.Node{
														node.CreateIdentifier(9, "list"),
														node.CreateNumber(9, "0"),
														node.CreateBinaryExpression(
															node.CreateIdentifier(9, "list_len"),
															CreateTokenWithLineNum(tokens.MINUS_TOKEN, 9),
															node.CreateNumber(9, "2"),
														),
													}),
												),
												CreateTokenWithLineNum(tokens.SEND_TOKEN, 9),
												node.CreateIdentifier(9, "new_value"),
											),
										},
									),
								),
							},

							// ELSE CASE
							node.CreateBlockStatements(
								[]node.Node{
									node.CreateBinaryExpression(
										node.CreateBinaryExpression(
											node.CreateBinaryExpression(
												node.CreateBuiltinFunctionIdentifier(12, "slice"),
												CreateTokenWithLineNum(tokens.SEND_TOKEN, 12),
												node.CreateList(12, []node.Node{
													node.CreateIdentifier(12, "list"),
													node.CreateNumber(12, "0"),
													node.CreateBinaryExpression(
														node.CreateIdentifier(12, "pos"),
														CreateTokenWithLineNum(tokens.MINUS_TOKEN, 12),
														node.CreateNumber(12, "1"),
													),
												}),
											),
											CreateTokenWithLineNum(tokens.SEND_TOKEN, 12),
											node.CreateIdentifier(12, "new_value"),
										),
										CreateTokenWithLineNum(tokens.SEND_TOKEN, 12),
										node.CreateBinaryExpression(
											node.CreateBuiltinFunctionIdentifier(12, "slice"),
											CreateTokenWithLineNum(tokens.SEND_TOKEN, 12),
											node.CreateList(12, []node.Node{
												node.CreateIdentifier(12, "list"),
												node.CreateBinaryExpression(
													node.CreateIdentifier(12, "pos"),
													CreateTokenWithLineNum(tokens.PLUS_TOKEN, 12),
													node.CreateNumber(12, "1"),
												),
												node.CreateBinaryExpression(
													node.CreateIdentifier(12, "list_len"),
													CreateTokenWithLineNum(tokens.MINUS_TOKEN, 12),
													node.CreateNumber(12, "1"),
												),
											}),
										),
									),
								},
							),
						),
					),

					// FUNCTION THIRD EXPRESSION
					node.CreateReturnStatement(
						15,
						node.CreateBinaryExpression(
							node.CreateBuiltinFunctionIdentifier(15, "unwrap"),
							CreateTokenWithLineNum(tokens.SEND_TOKEN, 15),
							node.CreateList(15, []node.Node{
								node.CreateIdentifier(15, "new_numbers"),
								node.CreateList(15, []node.Node{}),
							}),
						),
					),
				},
			),
		),

		node.CreateList(18, []node.Node{
			node.CreateNumber(18, "1"),
			node.CreateNumber(18, "2"),
			node.CreateNumber(18, "3"),
			node.CreateNumber(18, "4"),
			node.CreateNumber(18, "5"),
			node.CreateNumber(18, "6"),
			node.CreateNumber(18, "7"),
		}),

		node.CreateNumber(19, "7"),

		node.CreateNumber(20, "4"),

		node.CreateList(21, []node.Node{
			node.CreateNumber(21, "1"),
			node.CreateNumber(21, "2"),
			node.CreateNumber(21, "3"),
			node.CreateNumber(21, "4"),
			node.CreateNumber(21, "20"),
			node.CreateNumber(21, "6"),
			node.CreateNumber(21, "7"),
		}),
	}

	AssertNodesEqual(t, 0, expectedResult, actualResults)
}

func evaluateSource(t *testing.T, testFileName string) []node.Node {

	path := strings.Join([]string{INTEGRATION_TESTS_DIRECTORY, testFileName}, string(os.PathSeparator))
	source := utils.GetSource(path)

	tokenizer := tokens.NewTokenizer(source)
	parserObj, err := parser.NewParser(tokenizer)
	if err != nil {
		t.Fatal(err)
	}

	ast, err := parserObj.Parse()
	if err != nil {
		t.Fatal(err)
	}

	evaluatorObj := evaluator.NewEvaluator(*ast)
	results, err := evaluatorObj.Evaluate()
	if err != nil {
		t.Fatal(err)
	}
	return results
}
