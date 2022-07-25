package evaluator

import (
	"log"
	"my_lang/parser"
	"my_lang/tokens"
	"strconv"
)

type evaluator struct {
	ast []parser.Statement
}

func New(ast []parser.Statement) evaluator {
	return evaluator{ast: ast}
}

func (e *evaluator) Evaluate() []int {
	results := []int{}
	for _, stmt := range e.ast {
		result := e.evaluate(stmt.Expr)
		results = append(results, result)
	}
	return results
}

func (e *evaluator) evaluate(expr parser.Expression) int {

	switch expr := expr.(type) {

	case *parser.BinaryOperation:
		left := e.evaluate(expr.Left)
		right := e.evaluate(expr.Right)

		switch expr.OP.Type {
		case tokens.PLUS:
			return left + right
		case tokens.MINUS:
			return left - right
		case tokens.ASTERISK:
			return left * right
		case tokens.FORWARD_SLASH:
			return int(left / right)
		default:
			log.Fatalf("Invalid Operator: %s (%s)", expr.OP.Type, expr.OP.Literal)
		}

	case *parser.Number:
		value, _ := strconv.Atoi(expr.Value)
		return value
	}

	log.Fatalf("Invalid type %T", expr)
	return 0
}
