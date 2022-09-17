package evaluator

import (
	"boomerang/node"
	"boomerang/tokens"
	"fmt"
	"strconv"
)

type evaluator struct {
	ast []node.Node
}

func New(ast []node.Node) evaluator {
	return evaluator{ast: ast}
}

func (e *evaluator) Evaluate() []node.Node {
	results := []node.Node{}
	for _, stmt := range e.ast {
		result := e.evaluate(stmt)
		results = append(results, result)
	}
	return results
}

func (e *evaluator) evaluate(expr node.Node) node.Node {

	switch expr.Type {

	case "BinaryExpression":
		left := e.evaluate(expr.GetParam("left"))
		right := e.evaluate(expr.GetParam("right"))
		op := expr.GetParam("operator")

		switch op.Type {
		case tokens.PLUS:
			result := e.toInt(left.Value) + e.toInt(right.Value)
			return e.createNumberNode(result)
		case tokens.MINUS:
			result := e.toInt(left.Value) - e.toInt(right.Value)
			return e.createNumberNode(result)
		case tokens.ASTERISK:
			result := e.toInt(left.Value) * e.toInt(right.Value)
			return e.createNumberNode(result)
		case tokens.FORWARD_SLASH:
			if right.Value == "0" {
				panic("Cannot divide by zero.")
			}
			result := e.toInt(left.Value) / e.toInt(right.Value)
			return e.createNumberNode(result)
		default:
			panic(fmt.Sprintf("Invalid Operator: %s (%s)", op.Type, op.Value))
		}

	case "Number":
		return expr
	}

	panic(fmt.Sprintf("Invalid type %T", expr.Type))
}

func (e *evaluator) toInt(s string) int {
	intVal, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("Cannot convert string to int: %s", s))
	}
	return intVal
}

func (e *evaluator) createNumberNode(value int) node.Node {
	return node.Node{Type: "Number", Value: fmt.Sprint(value)}
}
