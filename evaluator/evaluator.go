package evaluator

import (
	"fmt"
	"my_lang/parser"
	"my_lang/tokens"
	"strconv"
)

type evaluator struct {
	ast []parser.Node
}

func New(ast []parser.Node) evaluator {
	return evaluator{ast: ast}
}

func (e *evaluator) Evaluate() []int {
	results := []int{}
	for _, stmt := range e.ast {
		result := e.evaluate(stmt.GetParam("Expression"))
		results = append(results, result)
	}
	return results
}

func (e *evaluator) evaluate(expr parser.Node) int {

	switch expr.Type {

	case "Statement":
		statementExpression := expr.GetParam("Expression")
		return e.evaluate(statementExpression)

	case "BinaryExpression":
		left := e.evaluate(expr.GetParam("left"))
		right := e.evaluate(expr.GetParam("right"))
		op := expr.GetParam("operator")

		switch op.Type {
		case tokens.PLUS:
			return left + right
		case tokens.MINUS:
			return left - right
		case tokens.ASTERISK:
			return left * right
		case tokens.FORWARD_SLASH:
			return int(left / right)
		default:
			panic(fmt.Sprintf("Invalid Operator: %s (%s)", op.Type, op.Value))
		}

	case "Number":
		value, _ := strconv.Atoi(expr.Value)
		return value
	}

	panic(fmt.Sprintf("Invalid type %T", expr.Type))
}
