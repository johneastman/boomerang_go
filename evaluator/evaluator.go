package evaluator

import (
	"boomerang/node"
	"boomerang/tokens"
	"fmt"
	"strconv"
)

type evaluator struct {
	ast []node.Node
	env map[string]node.Node
}

func New(ast []node.Node) evaluator {
	return evaluator{
		ast: ast,
		env: map[string]node.Node{},
	}
}

func (e *evaluator) Evaluate() []node.Node {
	return e.evaluateStatements(e.ast)
}

func (e *evaluator) evaluateStatements(stmts []node.Node) []node.Node {
	results := []node.Node{}
	for _, stmt := range stmts {
		if result, isExpr := e.evaluateStatement(stmt); isExpr {
			results = append(results, *result)
		}
	}
	return results
}

func (e *evaluator) evaluateStatement(stmt node.Node) (*node.Node, bool) {
	if stmt.Type == node.ASSIGN_STMT {
		variable := stmt.GetParam(node.ASSIGN_STMT_IDENTIFIER)
		value := e.evaluateExpression(stmt.GetParam(node.EXPR))
		e.env[variable.Value] = value
		return nil, false

	} else if stmt.Type == node.PRINT_STMT {
		for i, node := range stmt.Params {
			evaluatedParam := e.evaluateExpression(node)

			if i < len(stmt.Params)-1 {
				fmt.Printf("%s ", evaluatedParam.Value)
			} else {
				fmt.Println(evaluatedParam.Value)
			}
		}
		return nil, false
	}

	statementExpression := e.evaluateExpression(stmt)
	return &statementExpression, true
}

func (e *evaluator) evaluateExpression(expr node.Node) node.Node {

	switch expr.Type {

	case node.NUMBER:
		return expr

	case node.FUNCTION:
		return expr

	case node.PARAMETER:
		return expr

	case node.IDENTIFIER:
		return e.getVariable(expr.Value)

	case node.UNARY_EXPR:
		expression := e.evaluateExpression(expr.GetParam(node.EXPR))
		operator := expr.GetParam(node.OPERATOR)
		if operator.Type == tokens.MINUS {
			expressionValue := -e.toFloat(expression.Value)
			return e.createNumberNode(expressionValue)
		}
		panic(fmt.Sprintf("Invalid unary operator: %s", expr.Type))

	case node.BIN_EXPR:
		left := e.evaluateExpression(expr.GetParam(node.BIN_EXPR_LEFT))
		right := e.evaluateExpression(expr.GetParam(node.BIN_EXPR_RIGHT))
		op := expr.GetParam(node.OPERATOR)

		switch op.Type {
		case tokens.PLUS:
			result := e.toFloat(left.Value) + e.toFloat(right.Value)
			return e.createNumberNode(result)
		case tokens.MINUS:
			result := e.toFloat(left.Value) - e.toFloat(right.Value)
			return e.createNumberNode(result)
		case tokens.ASTERISK:
			result := e.toFloat(left.Value) * e.toFloat(right.Value)
			return e.createNumberNode(result)
		case tokens.FORWARD_SLASH:
			if right.Value == "0" {
				panic("Cannot divide by zero.")
			}
			result := e.toFloat(left.Value) / e.toFloat(right.Value)
			return e.createNumberNode(result)
		case tokens.POINTER:
			functionCall := node.CreateFunctionCall(left, right.Params)
			return e.evaluateExpression(functionCall)

		default:
			panic(fmt.Sprintf("Invalid Operator: %s (%s)", op.Type, op.Value))
		}

	case node.FUNCTION_CALL:
		callParams := expr.GetParam(node.CALL_PARAMS) // Parameters pass to function
		function := expr.GetParamByKeys([]string{node.IDENTIFIER, node.FUNCTION})

		if function.Type == node.IDENTIFIER {
			// If the function object is an identifier, retireve the actual function object from the environment
			function = e.getVariable(function.Value)
		}

		// Assert that the function object is, in fact, a callable function
		if function.Type != node.FUNCTION {
			panic(fmt.Sprintf("Cannot make function call on type %s", function.Type))
		}

		// Check that the number of arguments passed to the function matches the number of arguments in the function definition
		functionParams := function.GetParam(node.PARAMETER) // Parameters included in function definition
		if len(callParams.Params) != len(functionParams.Params) {
			panic(fmt.Sprintf("Expected %d arguments, got %d", len(functionParams.Params), len(callParams.Params)))
		}

		tmpEnv := e.env
		e.env = map[string]node.Node{}

		// Set parameters to environment
		for i := range callParams.Params {
			functionParam := functionParams.Params[i]
			callParam := callParams.Params[i]

			e.env[functionParam.Value] = e.evaluateExpression(callParam)
		}

		functionResults := e.evaluateStatements(function.GetParam(node.STMTS).Params)
		if len(functionResults) == 0 {
			panic("Function returns nothing")
		}

		e.env = tmpEnv
		return functionResults[len(functionResults)-1] // Return the results of the last statement in the function
	}

	panic(fmt.Sprintf("Invalid type %#v", expr.Type))
}

func (e *evaluator) toFloat(s string) float64 {
	intVal, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(fmt.Sprintf("Cannot convert string to number: %s", s))
	}
	return intVal
}

func (e *evaluator) createNumberNode(value float64) node.Node {
	return node.CreateNumber(fmt.Sprint(value))
}

func (e *evaluator) getVariable(name string) node.Node {
	if value, ok := e.env[name]; ok {
		return value
	}
	panic(fmt.Sprintf("Undefined variable: %s", name))
}
