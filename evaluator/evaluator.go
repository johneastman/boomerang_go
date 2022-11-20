package evaluator

import (
	"boomerang/node"
	"boomerang/tokens"
	"boomerang/utils"
	"fmt"
	"math/rand"
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

func (e *evaluator) Evaluate() ([]node.Node, error) {
	return e.evaluateGlobalStatements(e.ast)
}

func (e *evaluator) evaluateGlobalStatements(stmts []node.Node) ([]node.Node, error) {
	results := []node.Node{}
	for _, stmt := range stmts {
		result, err := e.evaluateStatement(stmt)
		if err != nil {
			return nil, err
		}

		// If 'result' is not nil, then the statement returned a value (likely an expression statement)
		if result != nil {
			if result.Type == node.BREAK || result.Type == node.CONTINUE || result.Type == node.RETURN {
				return nil, utils.CreateError(result.LineNum, "%s statements not allowed outside loops", result.Value)
			}
			results = append(results, *result)
		}
	}
	return results, nil
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

		if result != nil {
			if result.Type == node.BREAK || result.Type == node.CONTINUE || result.Type == node.RETURN {
				return result, nil
			}
		}

		returnValue = result
	}

	returnValue = node.CreateBlockStatementReturnValue(lineNum, returnValue).Ptr()
	return returnValue, nil
}

func (e *evaluator) evaluateStatement(stmt node.Node) (*node.Node, error) {

	switch stmt.Type {

	case node.BREAK, node.CONTINUE, node.RETURN:
		return &stmt, nil

	case node.WHILE_LOOP:
		returnValue, err := e.evaluateWhileLoop(stmt)
		if err != nil {
			return nil, err
		}

		if returnValue != nil && returnValue.Type == node.RETURN {
			return returnValue, nil
		}
		return nil, nil

	default:
		return e.evaluateExpression(stmt)
	}
}

func (e *evaluator) evaluateAssignmentStatement(stmt node.Node) (*node.Node, error) {
	variable := stmt.GetParam(node.ASSIGN_STMT_IDENTIFIER)       // identifier, list of identifiers
	value, err := e.evaluateExpression(stmt.GetParam(node.EXPR)) // actual value(s)
	if err != nil {
		return nil, err
	}

	if variable.Type == node.IDENTIFIER {
		// Check that the user hasn't created a variable with the same name as a builtin construct
		if IsBuiltin(variable.Value) {
			return nil, utils.CreateError(
				stmt.LineNum,
				"%#v is a builtin function or variable",
				variable.Value,
			)
		}

		value, err := e.evaluateExpression(*value)
		if err != nil {
			return nil, err
		}
		e.env.SetIdentifier(variable.Value, *value)
		return value, nil

	} else if variable.Type == node.LIST && value.Type == node.LIST {

		evaluatedValues := []node.Node{}

		assignments := e.partitionAssignmentVariables(variable, *value)
		for _, identifierPair := range assignments {

			identifier := identifierPair[0]
			identifierValue := identifierPair[1]

			if identifier.Type != node.IDENTIFIER {
				return nil, utils.CreateError(identifier.LineNum, "invalid type for assignment: %s", identifier.ErrorDisplay())
			}

			identifierValueEvaluated, err := e.evaluateExpression(identifierValue)
			if err != nil {
				return nil, err
			}
			e.env.SetIdentifier(identifier.Value, *identifierValueEvaluated)
			evaluatedValues = append(evaluatedValues, *identifierValueEvaluated)
		}

		// multiple assignment expressions return the full list on the right side of the assignment operator
		return node.CreateList(stmt.LineNum, evaluatedValues).Ptr(), nil
	}

	return nil, utils.CreateError(
		stmt.LineNum,
		"invalid type for assignment: %s",
		variable.ErrorDisplay(),
	)
}

func (e *evaluator) partitionAssignmentVariables(identifiers, values node.Node) [][]node.Node {

	var identifierValuePairs = [][]node.Node{} // Map identifiers to their corresponding values

	/*
		Iterate though first (n - 1) identifiers. If an identifier has an associated value, pair the identifier with
		that value. Otherwise, pair the identifier with an empty monad object.
	*/
	index := 0
	for ; index < len(identifiers.Params)-1; index++ {
		identifier := identifiers.Params[index]

		var value node.Node
		if index >= len(values.Params) {
			// If the number of identifiers is greater than the number of values, set subsequent variables to an empty monad.
			value = node.CreateMonad(identifier.LineNum, nil)
		} else {
			value = values.Params[index]
		}

		identifierValuePairs = append(identifierValuePairs, []node.Node{identifier, value})
	}

	lastIdentifier := identifiers.Params[index]
	var lastIdentifierValue node.Node

	/*
		For the last identifier, if the index is greater than or equal the number of values, we can assume there are more
		identifiers than values, so pair that identifier with an empty monad object.
	*/
	if index >= len(values.Params) {
		lastIdentifierValue = node.CreateMonad(lastIdentifier.LineNum, nil)
	} else {
		/*
			Otherwise, if the length of the remaining values is 1, pair that value with the last identifier. If the number of
			remaining values is greater than 1, put the remaining values in a list and assign that list to the last identifier.
		*/
		remainingNodes := values.Params[index:]
		switch len(remainingNodes) {
		case 1:
			lastIdentifierValue = values.Params[index]
		default:
			lastIdentifierValue = node.CreateList(lastIdentifier.LineNum, values.Params[index:])
		}
	}
	identifierValuePairs = append(identifierValuePairs, []node.Node{lastIdentifier, lastIdentifierValue})

	return identifierValuePairs
}

func (e *evaluator) evaluateWhileLoop(stmt node.Node) (*node.Node, error) {
	condition := stmt.GetParam(node.WHILE_LOOP_CONDITION)
	statements := stmt.GetParam(node.WHILE_LOOP_STATEMENTS)

	for {
		evaluatedCondition, err := e.evaluateExpression(condition)
		if err != nil {
			return nil, err
		}

		if evaluatedCondition.Equals(node.CreateBooleanTrue(stmt.LineNum)) {
			stmt, err := e.evaluateBlockStatements(statements)
			if err != nil {
				return nil, err
			}

			if stmt.Type == node.BREAK {
				break
			}

			if stmt.Type == node.CONTINUE {
				continue
			}

			if stmt.Type == node.RETURN {
				return stmt, nil
			}

			// TODO: Figure out returns in while loops
		} else {
			break
		}
	}
	return nil, nil
}

func (e *evaluator) evaluateExpression(expr node.Node) (*node.Node, error) {

	switch expr.Type {

	case node.NUMBER, node.BOOLEAN, node.FUNCTION, node.BUILTIN_FUNCTION, node.MONAD:
		// Builtin functions will be evaluated later during a function call
		return &expr, nil

	case node.STRING:
		return e.evaluateString(expr)

	case node.LIST:
		return e.evaluateParameter(expr)

	case node.IDENTIFIER:
		return e.evaluateIdentifier(expr)

	case node.BUILTIN_VARIABLE:
		return evaluateBuiltinFunction(expr.Value, e, expr.LineNum, []node.Node{})

	case node.UNARY_EXPR:
		return e.evaluateUnaryExpression(expr)

	case node.BIN_EXPR:
		return e.evaluateBinaryExpression(expr)

	case node.ASSIGN_STMT:
		return e.evaluateAssignmentStatement(expr)

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

	elementVariableExpression := expr.GetParam(node.FOR_LOOP_ELEM_ASSIGN)
	if err := utils.CheckTypeError(lineNum, elementVariableExpression.Type, node.ASSIGN_STMT); err != nil {
		return nil, err
	}

	list := elementVariableExpression.GetParam(node.EXPR)
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

	variables := elementVariableExpression.GetParam(node.IDENTIFIER)
	statements := expr.GetParam(node.BLOCK_STATEMENTS)

	var values = []node.Node{}

	for _, element := range evaluatedList.Params {
		// Assign the placeholder/element variable to the value of the current list element
		placeHolderVariable := node.CreateAssignmentNode(variables, element)
		_, err = e.evaluateStatement(placeHolderVariable)
		if err != nil {
			return nil, err
		}

		// Evaluate the block statements in the for-loop
		result, err := e.evaluateBlockStatements(statements)
		if err != nil {
			return nil, err
		}

		if result != nil {
			switch result.Type {
			case node.BREAK:
				return node.CreateList(lineNum, values).Ptr(), nil
			case node.CONTINUE:
				continue
			case node.RETURN:
				return result, nil
			}
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
		floatValue := utils.ConvertStringToFloat(expression.Value)
		if floatValue == nil {
			return nil, utils.NotANumberError(expression.LineNum, expression.Value)
		}
		return node.CreateNumber(unaryExpression.LineNum, utils.FloatToString(-*floatValue)).Ptr(), nil

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

	leftNode := binaryExpression.GetParam(node.LEFT)
	op := binaryExpression.GetParam(node.OPERATOR)
	rightNode := binaryExpression.GetParam(node.RIGHT)

	left, err := e.evaluateExpression(leftNode)
	if err != nil {
		return nil, err
	}

	right, err := e.evaluateExpression(rightNode)
	if err != nil {
		return nil, err
	}

	switch op.Type {

	case tokens.PLUS:
		return e.add(*left, *right)

	case tokens.MINUS:
		return e.subtract(*left, *right)

	case tokens.ASTERISK:
		return e.multuply(*left, *right)

	case tokens.FORWARD_SLASH:
		return e.divide(*left, *right)

	case tokens.MODULO:
		return e.modulo(*left, *right)

	case tokens.SEND:
		return e.send(*left, *right)

	case tokens.AT:
		return e.index(*left, *right)

	case tokens.EQ:
		return e.compareEQ(*left, *right)

	case tokens.NE:
		return e.compareNE(*left, *right)

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
				e.evaluateExpression(functionParam)
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

	if returnValue.Type == node.CONTINUE || returnValue.Type == node.BREAK {
		return returnValue, nil
	}

	var wrappedReturnValue *node.Node
	if returnValue.Type == node.RETURN {
		wrappedReturnValue, err = e.evaluateExpression(returnValue.GetParam(node.EXPR))
		if err != nil {
			return nil, err
		}
		wrappedReturnValue = node.CreateMonad(returnValue.LineNum, wrappedReturnValue).Ptr()

	} else {
		wrappedReturnValue = node.CreateMonad(returnValue.LineNum, nil).Ptr()
	}

	// Reset environment back to original scope environment
	e.env = *e.env.parentEnv

	return wrappedReturnValue, nil
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

func (e *evaluator) compareNE(left, right node.Node) (*node.Node, error) {

	var booleanValue string
	if left.Equals(right) {
		booleanValue = tokens.FALSE_TOKEN.Literal
	} else {
		booleanValue = tokens.TRUE_TOKEN.Literal
	}

	return node.CreateBoolean(left.LineNum, booleanValue).Ptr(), nil
}

func (e *evaluator) compareLT(left node.Node, right node.Node) (*node.Node, error) {

	if left.Type == node.NUMBER && right.Type == node.NUMBER {

		leftNum := utils.ConvertStringToFloat(left.Value)
		if leftNum == nil {
			return nil, utils.CreateError(left.LineNum, "cannot convert %s to float64", left.Value)
		}

		rightNum := utils.ConvertStringToFloat(right.Value)
		if rightNum == nil {
			return nil, utils.CreateError(right.LineNum, "cannot convert %s to float64", left.Value)
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
		index := utils.ConvertStringToInteger(right.Value)
		if index == nil {
			return nil, utils.CreateError(right.LineNum, "list index must be an integer")
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

		leftValue := utils.ConvertStringToFloat(left.Value)
		if leftValue == nil {
			return nil, utils.NotANumberError(left.LineNum, left.Value)
		}

		rightValue := utils.ConvertStringToFloat(right.Value)
		if rightValue == nil {
			return nil, utils.NotANumberError(right.LineNum, right.Value)
		}

		result := *leftValue + *rightValue

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

		leftValue := utils.ConvertStringToFloat(left.Value)
		if leftValue == nil {
			return nil, utils.NotANumberError(left.LineNum, left.Value)
		}

		rightValue := utils.ConvertStringToFloat(right.Value)
		if rightValue == nil {
			return nil, utils.NotANumberError(right.LineNum, right.Value)
		}

		result := *leftValue - *rightValue

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

		leftValue := utils.ConvertStringToFloat(left.Value)
		if leftValue == nil {
			return nil, utils.NotANumberError(left.LineNum, left.Value)
		}

		rightValue := utils.ConvertStringToFloat(right.Value)
		if rightValue == nil {
			return nil, utils.NotANumberError(right.LineNum, right.Value)
		}

		result := *leftValue * *rightValue

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

		leftValue := utils.ConvertStringToFloat(left.Value)
		if leftValue == nil {
			return nil, utils.NotANumberError(left.LineNum, left.Value)
		}

		rightValue := utils.ConvertStringToFloat(right.Value)
		if rightValue == nil {
			return nil, utils.NotANumberError(right.LineNum, right.Value)
		}

		result := *leftValue / *rightValue
		return node.CreateNumber(left.LineNum, utils.FloatToString(result)).Ptr(), nil
	}
	return nil, utils.CreateError(
		left.LineNum,
		"cannot divide types %s and %s",
		left.ErrorDisplay(),
		right.ErrorDisplay(),
	)
}

func (e *evaluator) modulo(left, right node.Node) (*node.Node, error) {
	if left.Type == node.NUMBER && right.Type == node.NUMBER {

		if right.Value == "0" {
			return nil, utils.CreateError(left.LineNum, "cannot divide by zero")
		}

		leftValue := utils.ConvertStringToInteger(left.Value)
		if leftValue == nil {
			return nil, utils.CreateError(left.LineNum, "modulo only valid for whole (integer) numbers")
		}

		rightValue := utils.ConvertStringToInteger(right.Value)
		if rightValue == nil {
			return nil, utils.CreateError(right.LineNum, "modulo only valid for whole (integer) numbers")
		}

		result := *leftValue % *rightValue
		return node.CreateNumber(left.LineNum, utils.IntToString(result)).Ptr(), nil
	}
	return nil, utils.CreateError(
		left.LineNum,
		"cannot use modulus operator on types %s and %s",
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
