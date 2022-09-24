package evaluator

import (
	"boomerang/node"
	"boomerang/tokens"
	"boomerang/utils"
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

func (e *evaluator) Evaluate() (*[]node.Node, error) {
	return e.evaluateStatements(e.ast)
}

func (e *evaluator) evaluateStatements(stmts []node.Node) (*[]node.Node, error) {
	results := []node.Node{}
	for _, stmt := range stmts {
		result, isExpr, err := e.evaluateStatement(stmt)
		if err != nil {
			return nil, utils.CreateError(err)
		}

		if isExpr {
			results = append(results, *result)
			if result.Type == node.RETURN {
				results = append(results, result.Params[0])
				break
			}
		}
	}
	return &results, nil
}

func (e *evaluator) evaluateStatement(stmt node.Node) (*node.Node, bool, error) {
	if stmt.Type == node.ASSIGN_STMT {
		e.evaluateAssignmentStatement(stmt)
		return nil, false, nil

	} else if stmt.Type == node.PRINT_STMT {
		e.evaluatePrintStatement(stmt)
		return nil, false, nil

	} else if stmt.Type == node.RETURN {
		param, err := e.evaluateExpression(stmt.Params[0])
		if err != nil {
			return nil, false, utils.CreateError(err)
		}
		stmt.Params[0] = *param
		return &stmt, true, nil
	}

	statementExpression, err := e.evaluateExpression(stmt)
	if err != nil {
		return nil, false, utils.CreateError(err)
	}
	return statementExpression, true, nil
}

func (e *evaluator) evaluateAssignmentStatement(stmt node.Node) *error {
	variable := stmt.GetParam(node.ASSIGN_STMT_IDENTIFIER)
	value, err := e.evaluateExpression(stmt.GetParam(node.EXPR))
	if err != nil {
		return &err
	}
	e.env.SetIdentifier(variable.Value, *value)
	return nil
}

func (e *evaluator) evaluatePrintStatement(stmt node.Node) *error {
	for i, node := range stmt.Params {
		evaluatedParam, err := e.evaluateExpression(node)
		if err != nil {
			return &err
		}

		if i < len(stmt.Params)-1 {
			fmt.Printf("%s ", evaluatedParam.String())
		} else {
			fmt.Println(evaluatedParam.String())
		}
	}
	return nil
}

func (e *evaluator) evaluateExpression(expr node.Node) (*node.Node, error) {

	switch expr.Type {

	case node.NUMBER:
		return &expr, nil

	case node.BOOLEAN:
		return &expr, nil

	case node.FUNCTION:
		return &expr, nil

	case node.STRING:
		return e.evaluateString(expr)

	case node.LIST:
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
		return nil, fmt.Errorf("invalid type %#v", expr.Type)
	}
}

func (e *evaluator) evaluateParameter(parameterExpression node.Node) (*node.Node, error) {
	for i := range parameterExpression.Params {

		parameter, err := e.evaluateExpression(parameterExpression.Params[i])
		if err != nil {
			return nil, utils.CreateError(err)
		}
		parameterExpression.Params[i] = *parameter
	}
	return &parameterExpression, nil
}

func (e *evaluator) evaluateString(stringExpression node.Node) (*node.Node, error) {
	for i, param := range stringExpression.Params {
		value, err := e.evaluateExpression(param)
		if err != nil {
			return nil, utils.CreateError(err)
		}
		stringExpression.Value = strings.Replace(stringExpression.Value, fmt.Sprintf("<%d>", i), value.Value, 1)
	}

	node := node.CreateString(stringExpression.Value, []node.Node{})
	return &node, nil
}

func (e *evaluator) evaluateUnaryExpression(unaryExpression node.Node) (*node.Node, error) {
	expression, err := e.evaluateExpression(unaryExpression.GetParam(node.EXPR))
	if err != nil {
		return nil, utils.CreateError(err)
	}
	operator := unaryExpression.GetParam(node.OPERATOR)
	if operator.Type == tokens.MINUS {
		expressionValue := -e.toFloat(expression.Value)

		node := e.createNumberNode(expressionValue)
		return &node, nil
	}

	return nil, fmt.Errorf("invalid unary operator: %s", unaryExpression.Type)
}

func (e *evaluator) evaluateBinaryExpression(binaryExpression node.Node) (*node.Node, error) {
	left, err := e.evaluateExpression(binaryExpression.GetParam(node.LEFT))
	if err != nil {
		return nil, utils.CreateError(err)
	}

	right, err := e.evaluateExpression(binaryExpression.GetParam(node.RIGHT))
	if err != nil {
		return nil, utils.CreateError(err)
	}

	op := binaryExpression.GetParam(node.OPERATOR)

	switch op.Type {

	case tokens.PLUS_TOKEN.Type:
		return e.add(*left, *right)

	case tokens.MINUS_TOKEN.Type:
		return e.subtract(*left, *right)

	case tokens.ASTERISK_TOKEN.Type:
		return e.multuply(*left, *right)

	case tokens.FORWARD_SLASH:
		return e.divide(*left, *right)

	case tokens.PTR_TOKEN.Type:
		return e.leftPointer(*left, *right)

	default:
		return nil, fmt.Errorf("invalid Operator: %s (%s)", op.Type, op.Value)
	}
}

func (e *evaluator) evaluateFunctionCall(functionCallExpression node.Node) (*node.Node, error) {
	callParams := functionCallExpression.GetParam(node.CALL_PARAMS) // Parameters pass to function

	function := functionCallExpression.GetParamByKeys([]string{node.IDENTIFIER, node.FUNCTION})

	if function.Type == node.IDENTIFIER {
		// If the function object is an identifier, retireve the actual function object from the environment
		identifierFunction, err := e.env.GetIdentifier(function.Value)
		if err != nil {
			return nil, utils.CreateError(err)
		}
		function = *identifierFunction
	}

	// Assert that the function object is, in fact, a callable function
	if function.Type != node.FUNCTION {
		return nil, fmt.Errorf("cannot make function call on type %s", function.Type)
	}

	// Check that the number of arguments passed to the function matches the number of arguments in the function definition
	functionParams := function.GetParam(node.LIST) // Parameters included in function definition
	if len(callParams.Params) != len(functionParams.Params) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(functionParams.Params), len(callParams.Params))
	}

	tmpEnv := e.env
	e.env = CreateEnvironment()

	// Set parameters to environment
	for i := range callParams.Params {
		functionParam := functionParams.Params[i]
		callParam := callParams.Params[i]

		value, err := e.evaluateExpression(callParam)
		if err != nil {
			return nil, utils.CreateError(err)
		}
		e.env.SetIdentifier(functionParam.Value, *value)
	}

	functionResults, err := e.evaluateStatements(function.GetParam(node.STMTS).Params)
	if err != nil {
		return nil, utils.CreateError(err)
	}
	e.env = tmpEnv

	var returnStatement *node.Node = nil
	if len(*functionResults) > 0 {
		returnStatement = &(*functionResults)[len(*functionResults)-1]
	}
	node := node.CreateReturnValue(returnStatement)
	return &node, nil
}

func (e *evaluator) add(left node.Node, right node.Node) (*node.Node, error) {
	if left.Type == node.NUMBER && right.Type == node.NUMBER {
		result := e.toFloat(left.Value) + e.toFloat(right.Value)

		node := e.createNumberNode(result)
		return &node, nil
	}
	return nil, fmt.Errorf("cannot add types %s and %s", left.Type, right.Type)
}

func (e *evaluator) subtract(left node.Node, right node.Node) (*node.Node, error) {
	if left.Type == node.NUMBER && right.Type == node.NUMBER {
		result := e.toFloat(left.Value) - e.toFloat(right.Value)

		node := e.createNumberNode(result)
		return &node, nil
	}
	return nil, fmt.Errorf("cannot subtract types %s and %s", left.Type, right.Type)
}

func (e *evaluator) multuply(left node.Node, right node.Node) (*node.Node, error) {
	if left.Type == node.NUMBER && right.Type == node.NUMBER {
		result := e.toFloat(left.Value) * e.toFloat(right.Value)

		node := e.createNumberNode(result)
		return &node, nil
	}
	return nil, fmt.Errorf("cannot subtract types %s and %s", left.Type, right.Type)
}

func (e *evaluator) divide(left node.Node, right node.Node) (*node.Node, error) {
	if left.Type == node.NUMBER && right.Type == node.NUMBER {

		if right.Value == "0" {
			return nil, fmt.Errorf("cannot divide by zero")
		}
		result := e.toFloat(left.Value) / e.toFloat(right.Value)

		node := e.createNumberNode(result)
		return &node, nil
	}
	return nil, fmt.Errorf("cannot subtract types %s and %s", left.Type, right.Type)
}

func (e *evaluator) leftPointer(left node.Node, right node.Node) (*node.Node, error) {
	if left.Type == node.FUNCTION && right.Type == node.LIST {
		functionCall := node.CreateFunctionCall(left, right.Params)

		return e.evaluateExpression(functionCall)

	} else if left.Type == node.BUILTIN_FUNC && right.Type == node.LIST {
		// For builtin functions, Node.Value stores the builtin-function
		return e.evaluateBuiltinFunction(left.Value, right)
	}

	return nil, fmt.Errorf("cannot use left pointer on types %s and %s", left.Type, right.Type)
}

func (e *evaluator) evaluateBuiltinFunction(builtinFunctionType string, parameters node.Node) (*node.Node, error) {

	switch builtinFunctionType {
	case node.BUILTIN_LEN:
		value := len(parameters.Params)

		node := node.CreateNumber(fmt.Sprint(value))
		return &node, nil

	case node.BUILTIN_UNWRAP:

		if parameters.Type != node.LIST {
			return nil, fmt.Errorf("invalid type for unwrap: %s. Expected %s", parameters.Type, node.LIST)
		}

		// TODO: Incorporate boolean value into computation
		returnParam := parameters.Params[0]
		defaultValue := parameters.Params[1]
		if len(returnParam.Params) == 1 {
			return e.evaluateExpression(defaultValue)
		}
		return e.evaluateExpression(returnParam.Params[1]) // Params[0] contains the boolean value

	default:
		return nil, fmt.Errorf("undefined builtin function: %s", builtinFunctionType)
	}
}

func (e *evaluator) toFloat(s string) float64 {
	floatVal, err := strconv.ParseFloat(s, 64)
	if err != nil {
		// TODO: May need to change return type to (*float64, error) if type conversion is introduced
		panic(fmt.Sprintf("Cannot convert string to number: %s", s))
	}
	return floatVal
}

func (e *evaluator) createNumberNode(value float64) node.Node {
	return node.CreateNumber(fmt.Sprint(value))
}
