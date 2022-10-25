package evaluator

import (
	"boomerang/node"
	"boomerang/tokens"
	"boomerang/utils"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type evaluator struct {
	ast []node.Node
	env environment
}

func NewEvaluator(ast []node.Node) evaluator {

	rand.Seed(time.Now().UnixNano()) // for builtin "random" function

	return evaluator{
		ast: ast,
		env: CreateEnvironment(nil),
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
			if result.Type == node.BREAK {
				return nil, utils.CreateError(result.LineNum, "break statements not allowed outside loops")
			}
			results = append(results, *result)
		}
	}
	return &results, nil
}

func (e *evaluator) evaluateBlockStatements(statements node.Node) (*node.Node, error) {

	if statements.Type != node.BLOCK_STATEMENTS {
		panic(fmt.Sprintf("invalid type for block statement: %s", statements.ErrorDisplay()))
	}

	var returnValue *node.Node
	lineNum := statements.LineNum
	for _, statement := range statements.Params {
		lineNum = statement.LineNum
		result, err := e.evaluateStatement(statement)
		if err != nil {
			return nil, err
		}

		if result != nil && result.Type == node.BREAK {
			return result, nil
		}

		returnValue = result
	}

	returnValue = node.CreateBlockStatementReturnValue(lineNum, returnValue).Ptr()
	return returnValue, nil
}

func (e *evaluator) evaluateStatement(stmt node.Node) (*node.Node, error) {

	switch stmt.Type {

	case node.ASSIGN_STMT:
		if err := e.evaluateAssignmentStatement(stmt); err != nil {
			return nil, err
		}
		return nil, nil

	case node.WHILE_LOOP:
		if err := e.evaluateWhileLoop(stmt); err != nil {
			return nil, err
		}
		return nil, nil

	case node.BREAK:
		return &stmt, nil

	default:
		return e.evaluateExpression(stmt)
	}
}

func (e *evaluator) evaluateAssignmentStatement(stmt node.Node) error {
	variable := stmt.GetParam(node.ASSIGN_STMT_IDENTIFIER)

	// Check that the user hasn't created a variable with the same name as a builtin function or variable
	if IsBuiltinFunction(variable.Value) || IsBuiltinVariable(variable.Value) {
		return utils.CreateError(
			stmt.LineNum,
			"%#v is a builtin function or variable",
			variable.Value,
		)
	}

	value, err := e.evaluateExpression(stmt.GetParam(node.EXPR))
	if err != nil {
		return err
	}
	e.env.SetIdentifier(variable.Value, *value)
	return nil
}

func (e *evaluator) evaluateWhileLoop(stmt node.Node) error {
	condition := stmt.GetParam(node.WHILE_LOOP_CONDITION)
	statements := stmt.GetParam(node.WHILE_LOOP_STATEMENTS)

	for {
		evaluatedCondition, err := e.evaluateExpression(condition)
		if err != nil {
			return err
		}

		if evaluatedCondition.Equals(node.CreateBooleanTrue(stmt.LineNum)) {
			stmt, err := e.evaluateBlockStatements(statements)
			if err != nil {
				return err
			}

			if stmt.Type == node.BREAK {
				break
			}
		} else {
			break
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

	case node.BUILTIN_VARIABLE:
		return getBuiltinVariable(expr), nil

	case node.BUILTIN_FUNCTION:
		// These will be evaluated later during a function call
		return &expr, nil

	case node.UNARY_EXPR:
		return e.evaluateUnaryExpression(expr)

	case node.BIN_EXPR:
		return e.evaluateBinaryExpression(expr)

	case node.FUNCTION_CALL:
		return e.evaluateFunctionCall(expr)

	case node.WHEN:
		return e.evaluateWhenExpression(expr)

	case node.FOR_LOOP:
		return e.evaluateForLoop(expr)

	default:
		// This error will only happen if the developer has not implemented an expression type
		panic(fmt.Sprintf("invalid type %#v", expr.Type))
	}
}

func (e *evaluator) evaluateIdentifier(identifierExpression node.Node) (*node.Node, error) {
	return e.env.GetIdentifier(identifierExpression) // Get the user-defined variable from the environment
}

func (e *evaluator) evaluateParameter(parameterExpression node.Node) (*node.Node, error) {

	evaluatedParameters := []node.Node{}

	for i := range parameterExpression.Params {

		parameter, err := e.evaluateExpression(parameterExpression.Params[i])
		if err != nil {
			return nil, err
		}
		evaluatedParameters = append(evaluatedParameters, *parameter)
	}
	return node.CreateList(parameterExpression.LineNum, evaluatedParameters).Ptr(), nil
}

func (e *evaluator) evaluateString(stringExpression node.Node) (*node.Node, error) {
	for i, param := range stringExpression.Params {
		value, err := e.evaluateExpression(param)
		if err != nil {
			return nil, err
		}

		// With string interpolation, the quotes around strings should not be included in the final string
		var replacementString string
		if value.Type == node.STRING {
			replacementString = value.Value
		} else {
			replacementString = value.String()
		}
		stringExpression.Value = strings.Replace(stringExpression.Value, fmt.Sprintf("<%d>", i), replacementString, 1)
	}

	return node.CreateRawString(stringExpression.LineNum, stringExpression.Value).Ptr(), nil
}

func (e *evaluator) evaluateForLoop(expr node.Node) (*node.Node, error) {
	lineNum := expr.LineNum

	elementVariableName := expr.GetParam(node.FOR_LOOP_ELEMENT)
	list := expr.GetParam(node.LIST)
	statements := expr.GetParam(node.BLOCK_STATEMENTS)

	evaluatedList, err := e.evaluateExpression(list)
	if err != nil {
		return nil, err
	}

	if evaluatedList.Type != node.LIST {
		return nil, utils.CreateError(
			lineNum,
			"invalid type for for loop: %s",
			evaluatedList.ErrorDisplay(),
		)
	}

	var values = []node.Node{}

	for _, element := range evaluatedList.Params {
		// Assign the placeholder/element variable to the value of the current list element
		placeHolderVariable := node.CreateAssignmentStatement(lineNum, elementVariableName.Value, element)
		_, err = e.evaluateStatement(placeHolderVariable)
		if err != nil {
			return nil, err
		}

		// Evaluate the block statements in the for-loop
		result, err := e.evaluateBlockStatements(statements)
		if err != nil {
			return nil, err
		}

		if result != nil && result.Type == node.BREAK {
			return node.CreateList(lineNum, values).Ptr(), nil
		}

		values = append(values, *result)
	}

	return node.CreateList(lineNum, values).Ptr(), nil
}

func (e *evaluator) evaluateUnaryExpression(unaryExpression node.Node) (*node.Node, error) {
	expression, err := e.evaluateExpression(unaryExpression.GetParam(node.EXPR))
	if err != nil {
		return nil, err
	}
	operator := unaryExpression.GetParam(node.OPERATOR)
	if operator.Type == tokens.MINUS {

		if expression.Type != node.NUMBER {
			return nil, utils.CreateError(expression.LineNum, "invalid type for minus operator: %s", expression.ErrorDisplay())
		}
		expressionValue := -e.toFloat(expression.Value)
		return node.CreateNumber(unaryExpression.LineNum, utils.FloatToString(expressionValue)).Ptr(), nil

	} else if operator.Type == tokens.NOT {

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
		return node.CreateBoolean(expression.LineNum, literal).Ptr(), nil
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

	case tokens.PLUS:
		return e.add(*left, *right)

	case tokens.MINUS:
		return e.subtract(*left, *right)

	case tokens.ASTERISK:
		return e.multuply(*left, *right)

	case tokens.FORWARD_SLASH:
		return e.divide(*left, *right)

	case tokens.SEND:
		return e.send(*left, *right)

	case tokens.AT:
		return e.index(*left, *right)

	case tokens.EQ:
		return e.compareEQ(*left, *right)

	case tokens.LT:
		return e.compareLT(*left, *right)

	case tokens.IN:
		return e.compareIn(*left, *right)

	case tokens.OR:
		return e.booleanOr(*left, *right)

	case tokens.AND:
		return e.booleanAnd(*left, *right)

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

	if function.Type == node.BUILTIN_FUNCTION {
		return evaluateBuiltinFunction(function.Value, e, function.LineNum, callParams.Params)
	}

	if function.Type == node.IDENTIFIER {
		// If the function object is an identifier, retireve the actual function object from the environment
		identifierFunction, err := e.env.GetIdentifier(function)
		if err != nil {
			return nil, err
		}
		function = *identifierFunction
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

	oldEnv := e.env
	e.env = CreateEnvironment(&oldEnv)

	// Keyword arguments do not count as arguments passed to the function
	callParamsIndex := 0
	for _, functionParam := range functionParams.Params {
		if functionParam.Type == node.IDENTIFIER {

			if callParamsIndex >= len(callParams.Params) {
				return nil, utils.CreateError(
					function.LineNum,
					"Function paramter %#v does not have a value. Either add %d more parameters to the function call or assign %#v a default value in the function definition.",
					functionParam.Value,
					callParamsIndex-len(callParams.Params)+1,
					functionParam.Value,
				)
			}

			parameterValue := callParams.Params[callParamsIndex]

			evaluatedParameterValue, err := e.evaluateExpression(parameterValue)
			if err != nil {
				return nil, err
			}
			e.env.SetIdentifier(functionParam.Value, *evaluatedParameterValue)

		} else if functionParam.Type == node.ASSIGN_STMT {

			if callParamsIndex < len(callParams.Params) {
				parameterValue := callParams.Params[callParamsIndex]

				evaluatedParameterValue, err := e.evaluateExpression(parameterValue)
				if err != nil {
					return nil, err
				}
				parameterName := functionParam.GetParam(node.ASSIGN_STMT_IDENTIFIER).Value
				e.env.SetIdentifier(parameterName, *evaluatedParameterValue)

			} else {
				e.evaluateStatement(functionParam)
			}
		}

		callParamsIndex += 1
	}

	if callParamsIndex < len(callParams.Params) {
		/*
			After evaluating each expression, the value of "callParamsIndex" will be the expected number of call parameters,
			and "len(callParams.Params)" will be the number of call parameters provided.
		*/
		return nil, utils.CreateError(
			function.LineNum,
			"expected %d arguments, got %d",
			callParamsIndex,
			len(callParams.Params),
		)
	}

	functionStatements := function.GetParam(node.STMTS)
	if len(functionStatements.Params) == 0 {
		return node.CreateBlockStatementReturnValue(callParams.LineNum, nil).Ptr(), nil
	}

	returnValue, err := e.evaluateBlockStatements(functionStatements)
	if err != nil {
		return nil, err
	}

	// Reset environment back to original scope environment
	e.env = *e.env.parentEnv

	return returnValue, nil
}

func (e *evaluator) compareEQ(left node.Node, right node.Node) (*node.Node, error) {

	var booleanValue string
	if left.Equals(right) {
		booleanValue = tokens.TRUE_TOKEN.Literal
	} else {
		booleanValue = tokens.FALSE_TOKEN.Literal
	}

	return node.CreateBoolean(left.LineNum, booleanValue).Ptr(), nil
}

func (e *evaluator) compareLT(left node.Node, right node.Node) (*node.Node, error) {

	if left.Type == node.NUMBER && right.Type == node.NUMBER {

		leftNum, err := utils.ConvertStringToFloat(left.LineNum, left.Value)
		if err != nil {
			return nil, err
		}

		rightNum, err := utils.ConvertStringToFloat(right.LineNum, right.Value)
		if err != nil {
			return nil, err
		}

		if *leftNum < *rightNum {
			return node.CreateBooleanTrue(left.LineNum).Ptr(), nil
		}
		return node.CreateBooleanFalse(left.LineNum).Ptr(), nil
	}
	return nil, utils.CreateError(
		left.LineNum,
		"invalid types for less than: %s and %s",
		left.ErrorDisplay(),
		right.ErrorDisplay(),
	)
}

func (e *evaluator) compareIn(left node.Node, right node.Node) (*node.Node, error) {
	if right.Type == node.LIST {
		for _, value := range right.Params {
			if value.Equals(left) {
				return node.CreateBooleanTrue(left.LineNum).Ptr(), nil
			}
		}
		return node.CreateBooleanFalse(left.LineNum).Ptr(), nil
	}
	return nil, utils.CreateError(
		left.LineNum,
		"right side of \"in\" must be a list. Actual type: %s",
		right.ErrorDisplay(),
	)
}

func (e *evaluator) index(left node.Node, right node.Node) (*node.Node, error) {

	if right.Type == node.NUMBER {
		index, err := utils.ConvertStringToInteger(right.LineNum, right.Value)
		if err != nil {
			return nil, err
		}
		indexLiteral := *index

		switch left.Type {
		case node.LIST:
			if err := utils.CheckOutOfRange(left.LineNum, indexLiteral, len(left.Params)); err != nil {
				return nil, err
			}
			return left.Params[indexLiteral].Ptr(), nil
		case node.STRING:
			if err := utils.CheckOutOfRange(left.LineNum, indexLiteral, len(left.Value)); err != nil {
				return nil, err
			}
			character := left.Value[indexLiteral : indexLiteral+1]
			return node.CreateRawString(left.LineNum, character).Ptr(), nil
		}
	}

	return nil, utils.CreateError(
		right.LineNum,
		"invalid types for index: %s and %s",
		left.ErrorDisplay(),
		right.ErrorDisplay(),
	)
}

func (e *evaluator) add(left node.Node, right node.Node) (*node.Node, error) {
	if left.Type == node.NUMBER && right.Type == node.NUMBER {
		result := e.toFloat(left.Value) + e.toFloat(right.Value)

		return node.CreateNumber(left.LineNum, utils.FloatToString(result)).Ptr(), nil
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
		return node.CreateNumber(left.LineNum, utils.FloatToString(result)).Ptr(), nil
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

		return node.CreateNumber(left.LineNum, utils.FloatToString(result)).Ptr(), nil
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
		return node.CreateNumber(left.LineNum, utils.FloatToString(result)).Ptr(), nil
	}
	return nil, utils.CreateError(
		left.LineNum,
		"cannot divide types %s and %s",
		left.ErrorDisplay(),
		right.ErrorDisplay(),
	)
}

func (e *evaluator) send(left node.Node, right node.Node) (*node.Node, error) {
	if (left.Type == node.FUNCTION || left.Type == node.BUILTIN_FUNCTION) && right.Type == node.LIST {
		// Need to include "node.IDENTIFIER" check for builtin functions
		functionCall := node.CreateFunctionCall(left.LineNum, left, right.Params)
		return e.evaluateExpression(functionCall)

	} else if left.Type == node.LIST {

		nodes := left.Params
		if right.Type == node.LIST {
			nodes = append(nodes, right.Params...)
		} else {
			nodes = append(nodes, right)
		}

		return node.CreateList(left.LineNum, nodes).Ptr(), nil
	}

	return nil, utils.CreateError(
		left.LineNum,
		"cannot use send on types %s and %s",
		left.ErrorDisplay(),
		right.ErrorDisplay(),
	)
}

func (e *evaluator) booleanOr(left node.Node, right node.Node) (*node.Node, error) {
	if left.Type != node.BOOLEAN || right.Type != node.BOOLEAN {
		return nil, utils.CreateError(
			left.LineNum,
			"invalid types for boolean or. left: %s, right: %s",
			left.ErrorDisplay(),
			right.ErrorDisplay(),
		)
	}

	// Line number does not matter here because we're just checking if "left" or "right" are boolean true
	trueNode := node.CreateBooleanTrue(0)

	if left.Equals(trueNode) || right.Equals(trueNode) {
		return node.CreateBooleanTrue(left.LineNum).Ptr(), nil
	}
	return node.CreateBooleanFalse(left.LineNum).Ptr(), nil
}

func (e *evaluator) booleanAnd(left node.Node, right node.Node) (*node.Node, error) {
	if left.Type != node.BOOLEAN || right.Type != node.BOOLEAN {
		return nil, utils.CreateError(
			left.LineNum,
			"invalid types for boolean and. left: %s, right: %s",
			left.ErrorDisplay(),
			right.ErrorDisplay(),
		)
	}

	if left.String() == tokens.TRUE_TOKEN.Literal && right.String() == tokens.TRUE_TOKEN.Literal {
		return node.CreateBooleanTrue(left.LineNum).Ptr(), nil
	}
	return node.CreateBooleanFalse(left.LineNum).Ptr(), nil
}

func (e *evaluator) evaluateWhenExpression(whenExpression node.Node) (*node.Node, error) {

	expression, err := e.evaluateExpression(whenExpression.GetParam(node.WHEN_VALUE))
	if err != nil {
		return nil, err
	}

	cases := whenExpression.GetParam(node.WHEN_CASES)

	for _, _case := range cases.Params {
		caseValue, err := e.evaluateExpression(_case.GetParam(node.CASE_VALUE))
		if err != nil {
			return nil, err
		}

		if caseValue.Equals(*expression) {
			return e.evaluateBlockStatements(_case.GetParam(node.CASE_STMTS))
		}
	}

	// If none of the cases match, the else/default case will be returned.
	return e.evaluateBlockStatements(whenExpression.GetParam(node.WHEN_CASES_DEFAULT))
}

func (e *evaluator) toFloat(s string) float64 {
	floatVal, err := strconv.ParseFloat(s, 64)
	if err != nil {
		// TODO: in this error message, may need to replace "number" with "float" if type conversion is introduced
		panic(fmt.Sprintf("Cannot convert string to number: %s", s))
	}
	return floatVal
}

func (e *evaluator) evaluateAndCheckType(expression node.Node, expectedType string) (*node.Node, error) {
	evaluatedExpression, err := e.evaluateExpression(expression)
	if err != nil {
		return nil, err
	}

	if evaluatedExpression.Type != expectedType {
		return nil, utils.CreateError(
			evaluatedExpression.LineNum,
			"expected %s, got %s",
			expectedType,
			evaluatedExpression.ErrorDisplay(),
		)
	}
	return evaluatedExpression, nil
}
