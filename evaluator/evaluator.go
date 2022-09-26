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

func (e *evaluator) Evaluate() (*[]node.Node, error) {
	return e.evaluateGlobalStatements(e.ast)
}

func (e *evaluator) evaluateGlobalStatements(stmts []node.Node) (*[]node.Node, error) {
	results := []node.Node{}
	for _, stmt := range stmts {
		result, err := e.evaluateStatement(stmt)
		if err != nil {
			return nil, err
		}

		// If 'result' is not nil, then the statement returned a value (likely an expression statement)
		if result != nil {
			results = append(results, *result)
			if result.Type == node.RETURN {
				return nil, fmt.Errorf("%s statements not allowed in the global scope", tokens.RETURN_TOKEN.Literal)
			}
		}
	}
	return &results, nil
}

func (e *evaluator) evaluateBlockStatements(statements []node.Node) (*node.Node, error) {
	var returnValue *node.Node
	for _, statement := range statements {
		result, err := e.evaluateStatement(statement)
		if err != nil {
			return nil, err
		}
		returnValue = result

		/*
			Stop evaluating the statements if a return statement is found.

			Not all statements return a value, so check that the value is not nil before checking if the value
			is a return node.
		*/
		if result != nil && result.Type == node.RETURN {
			break
		}
	}
	return returnValue, nil
}

func (e *evaluator) evaluateStatement(stmt node.Node) (*node.Node, error) {
	if stmt.Type == node.ASSIGN_STMT {
		if err := e.evaluateAssignmentStatement(stmt); err != nil {
			return nil, err
		}
		return nil, nil

	} else if stmt.Type == node.PRINT_STMT {
		if err := e.evaluatePrintStatement(stmt); err != nil {
			return nil, err
		}
		return nil, nil

	} else if stmt.Type == node.RETURN {
		returnValue, err := e.evaluateExpression(stmt.GetParam(node.RETURN_VALUE))
		if err != nil {
			return nil, err
		}
		stmt.Params[0] = *returnValue
		return &stmt, nil

	} else if stmt.Type == node.IF_STMT {
		ifStatement, err := e.evaluateIfStatement(stmt)
		if err != nil {
			return nil, err
		}

		if ifStatement != nil {
			return ifStatement, nil
		}
		return nil, nil
	}

	statementExpression, err := e.evaluateExpression(stmt)
	if err != nil {
		return nil, err
	}
	return statementExpression, nil
}

func (e *evaluator) evaluateIfStatement(ifStatement node.Node) (*node.Node, error) {
	condition := ifStatement.GetParam(node.CONDITION)
	trueStatements := ifStatement.GetParam(node.TRUE_BRANCH)

	evaluatedCondition, err := e.evaluateExpression(condition)
	if err != nil {
		return nil, err
	}

	if evaluatedCondition.Value == tokens.TRUE_TOKEN.Literal {
		return e.evaluateBlockStatements(trueStatements.Params)
	}
	return nil, nil
}

func (e *evaluator) evaluateAssignmentStatement(stmt node.Node) error {
	variable := stmt.GetParam(node.ASSIGN_STMT_IDENTIFIER)
	value, err := e.evaluateExpression(stmt.GetParam(node.EXPR))
	if err != nil {
		return err
	}
	e.env.SetIdentifier(variable.Value, *value)
	return nil
}

func (e *evaluator) evaluatePrintStatement(stmt node.Node) error {
	for i, node := range stmt.Params {
		evaluatedParam, err := e.evaluateExpression(node)
		if err != nil {
			return err
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
			return nil, err
		}
		parameterExpression.Params[i] = *parameter
	}
	return &parameterExpression, nil
}

func (e *evaluator) evaluateString(stringExpression node.Node) (*node.Node, error) {
	for i, param := range stringExpression.Params {
		value, err := e.evaluateExpression(param)
		if err != nil {
			return nil, err
		}
		stringExpression.Value = strings.Replace(stringExpression.Value, fmt.Sprintf("<%d>", i), value.Value, 1)
	}

	node := node.CreateRawString(stringExpression.Value)
	return &node, nil
}

func (e *evaluator) evaluateUnaryExpression(unaryExpression node.Node) (*node.Node, error) {
	expression, err := e.evaluateExpression(unaryExpression.GetParam(node.EXPR))
	if err != nil {
		return nil, err
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
		return nil, err
	}

	right, err := e.evaluateExpression(binaryExpression.GetParam(node.RIGHT))
	if err != nil {
		return nil, err
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

	case tokens.OPEN_BRACKET_TOKEN.Type:
		return e.index(*left, *right)

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
			return nil, err
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
			return nil, err
		}
		e.env.SetIdentifier(functionParam.Value, *value)
	}

	functionStatements := function.GetParam(node.STMTS).Params
	if len(functionStatements) == 0 {
		node := node.CreateReturnValue(nil)
		return &node, nil
	}

	returnStatement, err := e.evaluateBlockStatements(functionStatements)
	if err != nil {
		return nil, err
	}

	if returnStatement != nil && returnStatement.Type == node.RETURN {
		returnStatement = node.Ptr(returnStatement.GetParam(node.RETURN_VALUE))
	}

	// Reset environment back to original scope environment
	e.env = tmpEnv

	return node.Ptr(node.CreateReturnValue(returnStatement)), nil
}

func (e *evaluator) index(left node.Node, right node.Node) (*node.Node, error) {
	if left.Type == node.LIST && right.Type == node.NUMBER {

		index, err := strconv.Atoi(right.Value)
		if err != nil {
			return nil, fmt.Errorf("index must be an integer")
		}

		if index >= len(left.Params) {
			return nil, fmt.Errorf("index out of range: %d. Length of list: %d", index, len(left.Params))
		}
		return node.Ptr(left.Params[index]), nil
	}
	return nil, fmt.Errorf("invalid types for index: %s and %s", left.Type, right.Type)
}

func (e *evaluator) add(left node.Node, right node.Node) (*node.Node, error) {
	if left.Type == node.NUMBER && right.Type == node.NUMBER {
		result := e.toFloat(left.Value) + e.toFloat(right.Value)

		return node.Ptr(e.createNumberNode(result)), nil
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
