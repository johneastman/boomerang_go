package evaluator

import (
	"boomerang/node"
	"boomerang/tokens"
	"boomerang/utils"
	"fmt"
	"math"
	"math/rand"
)

const (
	// Functions
	BUILTIN_LEN        = "len"
	BUILTIN_UNWRAP     = "unwrap"
	BUILTIN_SLICE      = "slice"
	BUILTIN_UNWRAP_ALL = "unwrap_all"
	BUILTIN_RANGE      = "range"
	BUILTIN_RANDOM     = "random"
	BUILTIN_PRINT      = "print"
	BUILTIN_INPUT      = "input"

	// Variables
	BUILTIN_PI = "pi"
)

type BuiltinFunction struct {
	NumArgs  int
	Function func(*evaluator, int, []node.Node) (*node.Node, error)
}

type BuiltinVariable struct {
	Literal string
	Type    string
}

/*
Initializing builtin functions with init() method to avoid initialization cycle error, in which the keys
in "builtinFunctions" call methods that use "builtinFunctions".

More info: https://go.dev/ref/spec#Package_initialization

For consistency, builtin variables are declared in "init" as well.
*/
var builtinFunctions map[string]BuiltinFunction
var builtinVariables map[string]BuiltinVariable

// Value for "BuiltinFunction.NumArgs" for function that can take any number of arguments (like "len").
var nArgsValue = -1

func init() {
	builtinFunctions = map[string]BuiltinFunction{
		BUILTIN_LEN:        {NumArgs: nArgsValue, Function: evaluateBuiltinLen},
		BUILTIN_UNWRAP:     {NumArgs: 2, Function: evaluateBuiltinUnwrap},
		BUILTIN_UNWRAP_ALL: {NumArgs: 2, Function: evaluateBuiltinUnwrapAll},
		BUILTIN_SLICE:      {NumArgs: 3, Function: evaluateBuiltinSlice},
		BUILTIN_RANGE:      {NumArgs: 2, Function: evaluateBuiltinRange},
		BUILTIN_RANDOM:     {NumArgs: 2, Function: evaluateBuiltinRandom},
		BUILTIN_PRINT:      {NumArgs: nArgsValue, Function: evaluateBuiltinPrint},
		BUILTIN_INPUT:      {NumArgs: 1, Function: evaluateBuiltinInput},
	}

	builtinVariables = map[string]BuiltinVariable{
		BUILTIN_PI: {Literal: fmt.Sprintf("%v", math.Pi), Type: node.NUMBER},
	}
}

func IsBuiltinVariable(value string) bool {
	if _, ok := builtinVariables[value]; ok {
		return true
	}
	return false
}

func IsBuiltinFunction(value string) bool {
	if _, ok := builtinFunctions[value]; ok {
		return true
	}
	return false
}

func getBuiltinVariable(identifierNode node.Node) *node.Node {
	if value, ok := builtinVariables[identifierNode.Value]; ok {
		switch value.Type {
		case node.NUMBER:
			return node.CreateNumber(identifierNode.LineNum, value.Literal).Ptr()
		default:
			panic(fmt.Sprintf("invalid type: %s", value.Type))
		}
	}
	return nil
}

func evaluateBuiltinFunction(name string, eval *evaluator, lineNum int, callParam []node.Node) (*node.Node, error) {
	builtinFunction := builtinFunctions[name]

	/*
		Check that the number of arguments passed to the builtin function is correct. Functions where the value
		is "nArgsValue" can accept any number of arguments.
	*/
	if builtinFunction.NumArgs != nArgsValue && builtinFunction.NumArgs != len(callParam) {
		return nil, utils.CreateError(
			lineNum,
			"incorrect number of arguments. expected %d, got %d",
			builtinFunction.NumArgs,
			len(callParam),
		)
	}

	return builtinFunction.Function(eval, lineNum, callParam)
}

func evaluateBuiltinSlice(eval *evaluator, lineNum int, callParam []node.Node) (*node.Node, error) {

	collection, err := eval.evaluateExpression(callParam[0])
	if err != nil {
		return nil, err
	}

	// Get the length of the collection based on the type. This is for verifying the indices are not out of range
	var collectionLength int
	switch collection.Type {
	case node.LIST:
		collectionLength = len(collection.Params)
	case node.STRING:
		collectionLength = len(collection.Value)
	default:
		return nil, utils.CreateError(collection.LineNum, "invalid type for slice: %s", collection.ErrorDisplay())
	}

	// Start Index
	startIndex, err := eval.evaluateAndCheckType(callParam[1], node.NUMBER)
	if err != nil {
		return nil, err
	}

	start := utils.ConvertStringToInteger(startIndex.Value)
	if start == nil {
		return nil, utils.CreateError(startIndex.LineNum, "start index must be an integer")
	}
	startLiteral := *start

	if err := utils.CheckOutOfRange(startIndex.LineNum, startLiteral, collectionLength); err != nil {
		return nil, err
	}

	// End Index
	endIndex, err := eval.evaluateAndCheckType(callParam[2], node.NUMBER)
	if err != nil {
		return nil, err
	}

	end := utils.ConvertStringToInteger(endIndex.Value)
	if end == nil {
		return nil, utils.CreateError(endIndex.LineNum, "end index must be an integer")
	}
	endLiteral := *end

	if err := utils.CheckOutOfRange(endIndex.LineNum, endLiteral, collectionLength); err != nil {
		return nil, err
	}

	if startLiteral > endLiteral {
		return nil, utils.CreateError(startIndex.LineNum, "start index cannot be greater than end index")
	}

	var returnNode node.Node

	/*
		This switch does not need a default case because that is handled in the above in the switch statement that
		gets the collection length
	*/
	switch collection.Type {

	case node.LIST:
		listValues := collection.Params
		slicedList := listValues[startLiteral : endLiteral+1]
		returnNode = node.CreateList(collection.LineNum, slicedList)

	case node.STRING:
		listValues := collection.Value
		slicedString := listValues[startLiteral : endLiteral+1]
		returnNode = node.CreateRawString(collection.LineNum, slicedString)
	}

	return returnNode.Ptr(), nil
}

func evaluateBuiltinUnwrap(eval *evaluator, lineNum int, callParameters []node.Node) (*node.Node, error) {
	/*
		I originally wanted "unwrap" to be implemented in pure Boomerang code, but because custom functions
		return a list and the purpose of unwrap is to extract the return value from that list, this implementation
		needs to be a builtin method.

		callParameters[0] contains the list returned by the function ("(true, <VALUE>)" or "(false)")
		callParameters[1] contains the default value, if the function returns "(false)"
	*/
	returnValueList, err := eval.evaluateExpression(callParameters[0])
	if err != nil {
		return nil, err
	}

	// Check that the first value passed to "unwrap" is a list
	if err := utils.CheckTypeError(lineNum, returnValueList.Type, node.LIST); err != nil {
		return nil, err
	}

	// "returnValueList.Params[0]" contains the boolean value, denoting whether the function returned a value
	returnValueListFirst, err := eval.evaluateExpression(returnValueList.Params[0])
	if err != nil {
		return nil, err
	}

	// Check that the first value in the  value monad is a boolean value
	if err := utils.CheckTypeError(lineNum, returnValueListFirst.Type, node.BOOLEAN); err != nil {
		return nil, err
	}

	// If the boolean value in the first element of the list is "true", return the function's actual return value
	if returnValueListFirst.Value == tokens.TRUE_TOKEN.Literal {
		// "returnValueList.Params[1]" contains the actual return value, if "returnValueList.Params[0]" is "true"
		return eval.evaluateExpression(returnValueList.Params[1])
	}

	// if "returnValueList.Params[0]" is "false", return the default value given to "unwrap".
	return eval.evaluateExpression(callParameters[1])
}

func evaluateBuiltinUnwrapAll(eval *evaluator, lineNum int, callParameters []node.Node) (*node.Node, error) {
	/*
		This function could easily be implemented in pure Boomerang code; for example:
		```
		unwrap_all = func(list, default) {
			newList = ()
			for e in list {
				newList = newList <- (unwrap <- (e, 0));
			};
			newList;
		};

		list = ((true, 1), (true, 2), (true, 3));
		unwrap_all <- (list, -1);
		```
		However, in the example above, because "unwrap" always returns a valid value, the return value would always be
		"(true, newList)". So, I decided "unwrap_all" should be a builtin method that just returns the list of values.
	*/
	list, err := eval.evaluateExpression(callParameters[0])
	if err != nil {
		return nil, err
	}

	if err := utils.CheckTypeError(lineNum, list.Type, node.LIST); err != nil {
		return nil, err
	}

	defaultValue := callParameters[1] // will be evaluated in "evaluateBuiltinUnwrap"

	unwrappedList := []node.Node{}

	for _, param := range list.Params {
		/*
			"param" is a monad ("(false)" or "(true, <VALUE>)"), so to utilize "evaluateBuiltinUnwrap",
			"param" and the default value are sent to "evaluateBuiltinUnwrap" as the call parameters.

			unwrap_all is essentially just calling "unwrap" on every element in the list (see example in comment above).
			```
		*/
		value, err := evaluateBuiltinUnwrap(eval, lineNum, []node.Node{param, defaultValue})
		if err != nil {
			return nil, err
		}
		unwrappedList = append(unwrappedList, *value)
	}
	return node.CreateList(lineNum, unwrappedList).Ptr(), nil
}

func evaluateBuiltinLen(eval *evaluator, lineNum int, callParameters []node.Node) (*node.Node, error) {
	value := len(callParameters)
	return node.CreateNumber(lineNum, utils.IntToString(value)).Ptr(), nil
}

func evaluateBuiltinRange(eval *evaluator, lineNum int, callParameters []node.Node) (*node.Node, error) {

	startNumber, err := eval.evaluateExpression(callParameters[0])
	if err != nil {
		return nil, err
	}

	if err := utils.CheckTypeError(lineNum, startNumber.Type, node.NUMBER); err != nil {
		return nil, err
	}

	startValue := utils.ConvertStringToInteger(startNumber.Value)
	if startValue == nil {
		return nil, utils.CreateError(lineNum, "start value must be an integer")
	}

	endNumber, err := eval.evaluateExpression(callParameters[1])
	if err != nil {
		return nil, err
	}

	if err := utils.CheckTypeError(lineNum, endNumber.Type, node.NUMBER); err != nil {
		return nil, err
	}

	endValue := utils.ConvertStringToInteger(endNumber.Value)
	if endValue == nil {
		return nil, utils.CreateError(lineNum, "end value must be an integer")
	}

	numbersNodeValues := []node.Node{}
	for i := *startValue; i <= *endValue; i++ {
		numberNode := node.CreateNumber(lineNum, utils.IntToString(i))
		numbersNodeValues = append(numbersNodeValues, numberNode)
	}

	return node.CreateList(lineNum, numbersNodeValues).Ptr(), nil
}

func evaluateBuiltinRandom(eval *evaluator, lineNum int, callParameters []node.Node) (*node.Node, error) {
	minNumber, err := eval.evaluateExpression(callParameters[0])
	if err != nil {
		return nil, err
	}

	if err := utils.CheckTypeError(lineNum, minNumber.Type, node.NUMBER); err != nil {
		return nil, err
	}

	minValue := utils.ConvertStringToInteger(minNumber.Value)
	if minValue == nil {
		return nil, utils.CreateError(lineNum, "min value must be an integer")
	}

	maxNumber, err := eval.evaluateExpression(callParameters[1])
	if err != nil {
		return nil, err
	}

	if err := utils.CheckTypeError(lineNum, maxNumber.Type, node.NUMBER); err != nil {
		return nil, err
	}

	maxValue := utils.ConvertStringToInteger(maxNumber.Value)
	if maxValue == nil {
		return nil, utils.CreateError(lineNum, "max value must be an integer")
	}

	if *minValue > *maxValue {
		return nil, utils.CreateError(
			minNumber.LineNum,
			"the minimum number, %d, cannot be greater than the maximum number, %d",
			*minValue,
			*maxValue,
		)
	}

	// "+ 1" ensures the generated number includes the maximum value
	randomValue := rand.Intn(*maxValue-*minValue+1) + *minValue
	return node.CreateNumber(minNumber.LineNum, utils.IntToString(randomValue)).Ptr(), nil
}

func evaluateBuiltinPrint(eval *evaluator, lineNum int, callParameters []node.Node) (*node.Node, error) {
	for i, value := range callParameters {
		evaluatedParam, err := eval.evaluateExpression(value)
		if err != nil {
			return nil, err
		}

		if i < len(callParameters)-1 {
			fmt.Printf("%s ", evaluatedParam.String())
		} else {
			fmt.Println(evaluatedParam.String())
		}
	}
	return node.CreateBlockStatementReturnValue(lineNum, nil).Ptr(), nil
}

func evaluateBuiltinInput(eval *evaluator, lineNum int, callParameters []node.Node) (*node.Node, error) {

	prompt, err := eval.evaluateExpression(callParameters[0])
	if err != nil {
		return nil, err
	}

	if err := utils.CheckTypeError(lineNum, prompt.Type, node.STRING); err != nil {
		return nil, err
	}

	inputValue := utils.UserInput(prompt.Value)

	return node.CreateRawString(lineNum, inputValue).Ptr(), nil
}
