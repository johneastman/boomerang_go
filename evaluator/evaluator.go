package evaluator

import (
	"boomerang/node"
	"boomerang/tokens"
	"boomerang/utils"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	BUILTIN_LEN    = "len"
	BUILTIN_UNWRAP = "unwrap"
)

var builtinFunctions = []string{
	BUILTIN_LEN,
	BUILTIN_UNWRAP,
}

var builtinVariables = map[string]string{
	"pi": fmt.Sprintf("%v", math.Pi),
}

func isBuiltinFunction(value string) bool {
	for _, builtinFunction := range builtinFunctions {
		if builtinFunction == value {
			return true
		}
	}
	return false
}

func getBuiltinVariable(value string) *string {
	if value, ok := builtinVariables[value]; ok {
		return &value
	}
	return nil
}

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
				return nil, utils.CreateError(
					result.LineNum,
					"%s statements not allowed in the global scope",
					tokens.RETURN_TOKEN.Literal,
				)
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

	if evaluatedCondition.Type != node.BOOLEAN {
		return nil, utils.CreateError(
			evaluatedCondition.LineNum,
			"invalid type for if-statement condition: %s",
			evaluatedCondition.ErrorDisplay(),
		)
	}

	if evaluatedCondition.String() == tokens.TRUE_TOKEN.Literal {
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
		return e.evaluateIdentifier(expr)

	case node.UNARY_EXPR:
		return e.evaluateUnaryExpression(expr)

	case node.BIN_EXPR:
		return e.evaluateBinaryExpression(expr)

	case node.FUNCTION_CALL:
		return e.evaluateFunctionCall(expr)

	default:
		// This error will only happen if the developer has not implemented an expression type
		panic(fmt.Sprintf("invalid type %#v", expr.Type))
	}
}

func (e *evaluator) evaluateIdentifier(identifierExpression node.Node) (*node.Node, error) {

	identifierName := identifierExpression.Value

	// Check for builtin variables
	variableValue := getBuiltinVariable(identifierName)
	if variableValue != nil {
		node := node.CreateNumber(identifierExpression.LineNum, *variableValue)
		return &node, nil
	}

	if isBuiltinFunction(identifierName) {
		return &identifierExpression, nil
	}

	// Variable defined in Boomerang file
	return e.env.GetIdentifier(identifierExpression)
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

	node := node.CreateRawString(stringExpression.LineNum, stringExpression.Value)
	return &node, nil
}

func (e *evaluator) evaluateUnaryExpression(unaryExpression node.Node) (*node.Node, error) {
	expression, err := e.evaluateExpression(unaryExpression.GetParam(node.EXPR))
	if err != nil {
		return nil, err
	}
	operator := unaryExpression.GetParam(node.OPERATOR)
	if operator.Type == tokens.MINUS_TOKEN.Type {

		if expression.Type != node.NUMBER {
			return nil, utils.CreateError(expression.LineNum, "invalid type for minus operator: %s", expression.ErrorDisplay())
		}

		expressionValue := -e.toFloat(expression.Value)

		node := e.createNumberNode(expressionValue, unaryExpression.LineNum)
		return &node, nil

	} else if operator.Type == tokens.NOT_TOKEN.Type {

		if expression.Type != node.BOOLEAN {
			return nil, utils.CreateError(
				expression.LineNum,
				"invalid type for bang operator: %s",
				expression.ErrorDisplay(),
			)
		}

		booleanValue := expression.Value

		var literal string
		if booleanValue == tokens.TRUE_TOKEN.Literal {
			literal = tokens.FALSE_TOKEN.Literal
		} else {
			literal = tokens.TRUE_TOKEN.Literal
		}
		node := node.CreateBoolean(literal, expression.LineNum)
		return &node, nil
	}

	return nil, utils.CreateError(
		unaryExpression.LineNum,
		"invalid unary operator: %s",
		operator.ErrorDisplay(),
	)
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
		return e.pointer(*left, *right)

	case tokens.AT_TOKEN.Type:
		return e.index(*left, *right)

	case tokens.EQ_TOKEN.Type:
		return e.compare(*left, *right)

	default:
		return nil, utils.CreateError(
			op.LineNum,
			"invalid binary operator: %s",
			op.ErrorDisplay(),
		)
	}
}

func (e *evaluator) evaluateFunctionCall(functionCallExpression node.Node) (*node.Node, error) {
	callParams := functionCallExpression.GetParam(node.CALL_PARAMS) // Parameters pass to function

	function := functionCallExpression.GetParamByKeys([]string{node.IDENTIFIER, node.FUNCTION})

	if function.Type == node.IDENTIFIER {
		switch function.Value {
		case BUILTIN_LEN:
			return e.evaluateBuiltinLen(function.LineNum, callParams.Params)

		case BUILTIN_UNWRAP:
			return e.evaluateBuiltinUnwrap(callParams.Params)

		default:
			// If the function object is an identifier, retireve the actual function object from the environment
			identifierFunction, err := e.env.GetIdentifier(function)
			if err != nil {
				return nil, err
			}
			function = *identifierFunction
		}
	}

	// Assert that the function object is, in fact, a callable function
	if function.Type != node.FUNCTION {
		return nil, utils.CreateError(
			function.LineNum,
			"cannot make function call on type %s",
			function.ErrorDisplay(),
		)
	}

	// Check that the number of arguments passed to the function matches the number of arguments in the function definition
	functionParams := function.GetParam(node.LIST) // Parameters included in function definition
	if len(callParams.Params) != len(functionParams.Params) {
		return nil, utils.CreateError(
			function.LineNum,
			"expected %d arguments, got %d",
			len(functionParams.Params),
			len(callParams.Params),
		)
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
		node := node.CreateFunctionReturnValue(callParams.LineNum, nil)
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

	node := node.CreateFunctionReturnValue(returnStatement.LineNum, returnStatement)
	return &node, nil
}

func (e *evaluator) evaluateBuiltinUnwrap(callParameters []node.Node) (*node.Node, error) {
	/*
		I originally wanted "unwrap" to be implemented in pure Boomerang code, but because custom functions
		return a list and the purpose of unwrap is to extract the return value from that list, this implementation
		needs to be a builtin method.
	*/
	returnValueList, err := e.evaluateExpression(callParameters[0]) // Params[0] contains the boolean value
	if err != nil {
		return nil, err
	}

	returnValueListFirst, err := e.evaluateExpression(returnValueList.Params[0])
	if err != nil {
		return nil, err
	}

	// If the boolean value in the first element of the list is "true", return the function's actual return value
	if returnValueListFirst.Value == tokens.TRUE_TOKEN.Literal {
		return e.evaluateExpression(returnValueList.Params[1]) // Params[1] contains the actual return value, if Params[0] is true
	}

	// Otherwise, return the provided default value
	return e.evaluateExpression(callParameters[1])
}

func (e *evaluator) evaluateBuiltinLen(lineNum int, callParameters []node.Node) (*node.Node, error) {
	value := len(callParameters)

	node := node.CreateNumber(
		lineNum,
		fmt.Sprint(value),
	)
	return &node, nil
}

func (e *evaluator) compare(left node.Node, right node.Node) (*node.Node, error) {

	var booleanValue string
	if left.String() == right.String() {
		booleanValue = tokens.TRUE_TOKEN.Literal
	} else {
		booleanValue = tokens.FALSE_TOKEN.Literal
	}

	booleanNode := node.CreateBoolean(booleanValue, left.LineNum)
	return &booleanNode, nil
}

func (e *evaluator) index(left node.Node, right node.Node) (*node.Node, error) {
	if left.Type == node.LIST && right.Type == node.NUMBER {

		index, err := strconv.Atoi(right.Value)
		if err != nil {
			return nil, utils.CreateError(left.LineNum, "index must be an integer")
		}

		if index >= len(left.Params) {
			return nil, utils.CreateError(left.LineNum, "index %d out of range. Length of list: %d", index, len(left.Params))
		}
		return node.Ptr(left.Params[index]), nil
	}
	return nil, utils.CreateError(
		left.LineNum,
		"invalid types for index: %s and %s",
		left.ErrorDisplay(),
		right.ErrorDisplay(),
	)
}

func (e *evaluator) add(left node.Node, right node.Node) (*node.Node, error) {
	if left.Type == node.NUMBER && right.Type == node.NUMBER {
		result := e.toFloat(left.Value) + e.toFloat(right.Value)

		return node.Ptr(e.createNumberNode(result, left.LineNum)), nil
	}
	return nil, utils.CreateError(
		left.LineNum,
		"cannot add types %s and %s",
		left.ErrorDisplay(),
		right.ErrorDisplay(),
	)
}

func (e *evaluator) subtract(left node.Node, right node.Node) (*node.Node, error) {
	if left.Type == node.NUMBER && right.Type == node.NUMBER {
		result := e.toFloat(left.Value) - e.toFloat(right.Value)

		node := e.createNumberNode(result, left.LineNum)
		return &node, nil
	}
	return nil, utils.CreateError(
		left.LineNum,
		"cannot subtract types %s and %s",
		left.ErrorDisplay(),
		right.ErrorDisplay(),
	)
}

func (e *evaluator) multuply(left node.Node, right node.Node) (*node.Node, error) {
	if left.Type == node.NUMBER && right.Type == node.NUMBER {
		result := e.toFloat(left.Value) * e.toFloat(right.Value)

		node := e.createNumberNode(result, left.LineNum)
		return &node, nil
	}
	return nil, utils.CreateError(
		left.LineNum,
		"cannot multiply types %s and %s",
		left.ErrorDisplay(),
		right.ErrorDisplay(),
	)
}

func (e *evaluator) divide(left node.Node, right node.Node) (*node.Node, error) {
	if left.Type == node.NUMBER && right.Type == node.NUMBER {

		if right.Value == "0" {
			return nil, utils.CreateError(left.LineNum, "cannot divide by zero")
		}
		result := e.toFloat(left.Value) / e.toFloat(right.Value)

		node := e.createNumberNode(result, left.LineNum)
		return &node, nil
	}
	return nil, utils.CreateError(
		left.LineNum,
		"cannot divide types %s and %s",
		left.ErrorDisplay(),
		right.ErrorDisplay(),
	)
}

func (e *evaluator) pointer(left node.Node, right node.Node) (*node.Node, error) {
	if (left.Type == node.FUNCTION || left.Type == node.IDENTIFIER) && right.Type == node.LIST {
		functionCall := node.CreateFunctionCall(left.LineNum, left, right.Params)
		return e.evaluateExpression(functionCall)

	} else if left.Type == node.LIST {

		nodes := left.Params
		if right.Type == node.LIST {
			nodes = append(nodes, right.Params...)
		} else {
			nodes = append(nodes, right)
		}

		listNode := node.CreateList(left.LineNum, nodes)
		return &listNode, nil
	}

	return nil, utils.CreateError(
		left.LineNum,
		"cannot use pointer on types %s and %s",
		left.ErrorDisplay(),
		right.ErrorDisplay(),
	)
}

func (e *evaluator) toFloat(s string) float64 {
	floatVal, err := strconv.ParseFloat(s, 64)
	if err != nil {
		// TODO: May need to change return type to (*float64, error) if type conversion is introduced
		panic(fmt.Sprintf("Cannot convert string to number: %s", s))
	}
	return floatVal
}

func (e *evaluator) createNumberNode(value float64, lineNum int) node.Node {
	return node.CreateNumber(lineNum, fmt.Sprint(value))
}
