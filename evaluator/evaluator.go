package evaluator

import (
	"boomerang/node"
	"boomerang/tokens"
	"fmt"
	"strconv"
	"strings"
)

type evaluator struct {
	ast []node.Node
	env environment
}

func New(ast []node.Node) evaluator {
	return evaluator{
		ast: ast,
		env: CreateEnvironment(),
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

			if result.Type == node.RETURN {
				results = append(results, result.Params[0])
				break
			}
		}
	}
	return results
}

func (e *evaluator) evaluateStatement(stmt node.Node) (*node.Node, bool) {
	if stmt.Type == node.ASSIGN_STMT {
		e.evaluateAssignmentStatement(stmt)
		return nil, false

	} else if stmt.Type == node.PRINT_STMT {
		e.evaluatePrintStatement(stmt)
		return nil, false

	} else if stmt.Type == node.RETURN {
		stmt.Params[0] = e.evaluateExpression(stmt.Params[0])
		return &stmt, true
	}

	statementExpression := e.evaluateExpression(stmt)
	return &statementExpression, true
}

func (e *evaluator) evaluateAssignmentStatement(stmt node.Node) {
	variable := stmt.GetParam(node.ASSIGN_STMT_IDENTIFIER)
	value := e.evaluateExpression(stmt.GetParam(node.EXPR))
	e.env.SetIdentifier(variable.Value, value)
}

func (e *evaluator) evaluatePrintStatement(stmt node.Node) {
	for i, node := range stmt.Params {
		evaluatedParam := e.evaluateExpression(node)

		if i < len(stmt.Params)-1 {
			fmt.Printf("%s ", evaluatedParam.String())
		} else {
			fmt.Println(evaluatedParam.String())
		}
	}
}

func (e *evaluator) evaluateExpression(expr node.Node) node.Node {

	switch expr.Type {

	case node.NUMBER:
		return expr

	case node.BOOLEAN:
		return expr

	case node.STRING:
		return e.evaluateString(expr)

	case node.FUNCTION:
		return expr

	case node.PARAMETER:
		return e.evaluateParameter(expr)

	case node.IDENTIFIER:
		return e.env.GetIdentifier(expr.Value)

	case node.UNARY_EXPR:
		return e.evaluateUnaryExpression(expr)

	case node.BIN_EXPR:
		return e.evaluateBinaryExpression(expr)

	case node.FUNCTION_CALL:
		return e.evaluateFunctionCall(expr)

	default:
		panic(fmt.Sprintf("Invalid type %#v", expr.Type))
	}
}

func (e *evaluator) evaluateParameter(parameterExpression node.Node) node.Node {
	for i := range parameterExpression.Params {
		parameterExpression.Params[i] = e.evaluateExpression(parameterExpression.Params[i])
	}
	return parameterExpression
}

func (e *evaluator) evaluateString(stringExpression node.Node) node.Node {
	for i, param := range stringExpression.Params {
		value := e.evaluateExpression(param)
		stringExpression.Value = strings.Replace(stringExpression.Value, fmt.Sprintf("<%d>", i), value.Value, 1)
	}
	return node.CreateString(stringExpression.Value, []node.Node{})
}

func (e *evaluator) evaluateUnaryExpression(unaryExpression node.Node) node.Node {
	expression := e.evaluateExpression(unaryExpression.GetParam(node.EXPR))
	operator := unaryExpression.GetParam(node.OPERATOR)
	if operator.Type == tokens.MINUS {
		expressionValue := -e.toFloat(expression.Value)
		return e.createNumberNode(expressionValue)
	}
	panic(fmt.Sprintf("Invalid unary operator: %s", unaryExpression.Type))
}

func (e *evaluator) evaluateBinaryExpression(binaryExpression node.Node) node.Node {
	left := e.evaluateExpression(binaryExpression.GetParam(node.LEFT))
	right := e.evaluateExpression(binaryExpression.GetParam(node.RIGHT))
	op := binaryExpression.GetParam(node.OPERATOR)

	switch op.Type {

	case tokens.PLUS:
		return e.add(left, right)

	case tokens.MINUS:
		return e.subtract(left, right)

	case tokens.ASTERISK:
		return e.multuply(left, right)

	case tokens.FORWARD_SLASH:
		return e.divide(left, right)

	case tokens.LEFT_PTR:
		return e.leftPointer(left, right)

	case tokens.RIGHT_PTR:
		return e.rightPointer(left, right)

	default:
		panic(fmt.Sprintf("Invalid Operator: %s (%s)", op.Type, op.Value))
	}
}

func (e *evaluator) evaluateFunctionCall(functionCallExpression node.Node) node.Node {
	callParams := functionCallExpression.GetParam(node.CALL_PARAMS) // Parameters pass to function
	function := functionCallExpression.GetParamByKeys([]string{node.IDENTIFIER, node.FUNCTION})

	if function.Type == node.IDENTIFIER {
		// If the function object is an identifier, retireve the actual function object from the environment
		function = e.env.GetIdentifier(function.Value)
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
	e.env = CreateEnvironment()

	// Set parameters to environment
	for i := range callParams.Params {
		functionParam := functionParams.Params[i]
		callParam := callParams.Params[i]

		e.env.SetIdentifier(functionParam.Value, e.evaluateExpression(callParam))
	}

	functionResults := e.evaluateStatements(function.GetParam(node.STMTS).Params)
	e.env = tmpEnv

	if len(functionResults) == 0 {
		panic("Function returns nothing")
	}
	return functionResults[len(functionResults)-1] // Return the results of the last statement in the function
}

func (e *evaluator) add(left node.Node, right node.Node) node.Node {
	if left.Type == node.NUMBER && right.Type == node.NUMBER {
		result := e.toFloat(left.Value) + e.toFloat(right.Value)
		return e.createNumberNode(result)
	}
	panic(fmt.Sprintf("cannot add types %s and %s", left.Type, right.Type))
}

func (e *evaluator) subtract(left node.Node, right node.Node) node.Node {
	if left.Type == node.NUMBER && right.Type == node.NUMBER {
		result := e.toFloat(left.Value) - e.toFloat(right.Value)
		return e.createNumberNode(result)
	}
	panic(fmt.Sprintf("cannot subtract types %s and %s", left.Type, right.Type))
}

func (e *evaluator) multuply(left node.Node, right node.Node) node.Node {
	if left.Type == node.NUMBER && right.Type == node.NUMBER {
		result := e.toFloat(left.Value) * e.toFloat(right.Value)
		return e.createNumberNode(result)
	}
	panic(fmt.Sprintf("cannot subtract types %s and %s", left.Type, right.Type))
}

func (e *evaluator) divide(left node.Node, right node.Node) node.Node {
	if left.Type == node.NUMBER && right.Type == node.NUMBER {

		if right.Value == "0" {
			panic("Cannot divide by zero.")
		}
		result := e.toFloat(left.Value) / e.toFloat(right.Value)
		return e.createNumberNode(result)
	}
	panic(fmt.Sprintf("cannot subtract types %s and %s", left.Type, right.Type))
}

func (e *evaluator) leftPointer(left node.Node, right node.Node) node.Node {
	if left.Type == node.FUNCTION && right.Type == node.PARAMETER {
		functionCall := node.CreateFunctionCall(left, right.Params)
		return e.evaluateExpression(functionCall)

	} else if left.Type == node.BUILTIN_FUNC && right.Type == node.PARAMETER {
		// For builtin functions, Node.Value stores the builtin-function
		return e.evaluateBuiltinFunction(left.Value, right)
	}

	panic(fmt.Sprintf("cannot use left pointer on types %s and %s", left.Type, right.Type))
}

func (e *evaluator) evaluateBuiltinFunction(builtinFunctionType string, right node.Node) node.Node {

	switch builtinFunctionType {
	case node.BUILTIN_LEN:
		value := len(right.Params)
		return node.CreateNumber(fmt.Sprint(value))
	default:
		panic(fmt.Sprintf("Undefined builtin function: %s", builtinFunctionType))
	}
}

func (e *evaluator) rightPointer(left node.Node, right node.Node) node.Node {
	if left.Type == node.PARAMETER && right.Type == node.FUNCTION {
		functionCall := node.CreateFunctionCall(right, left.Params)
		return e.evaluateExpression(functionCall)

	} else if left.Type == node.PARAMETER && right.Type == node.BUILTIN_FUNC {
		// For builtin functions, Node.Value stores the builtin-function
		return e.evaluateBuiltinFunction(right.Value, left)
	}
	panic(fmt.Sprintf("cannot use right pointer on types %s and %s", left.Type, right.Type))
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
