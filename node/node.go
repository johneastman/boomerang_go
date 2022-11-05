package node

import (
	"boomerang/tokens"
	"boomerang/utils"
	"fmt"
	"strings"
)

type Node struct {
	Type    string
	Value   string
	LineNum int
	Params  []Node
}

func (n *Node) ErrorDisplay() string {
	return fmt.Sprintf("%s (%#v)", n.Type, n.Value)
}

func (n *Node) GetParam(key string) Node {
	node, err := getParam(key, *n)
	if err != nil {
		panic(err.Error())
	}
	return *node
}

func (n *Node) GetParamByKeys(keys []string) Node {
	/*
		Check if any of the keys in a list of keys exist as parameters in a node. This method is useful for evaluating
		function calls, which could be a variable (e.g., "add = func() { a + b; }; add(1, 2)") or a function literal
		(e.g., "func(a, b) { a + b; }(1, 2);")
	*/
	for _, key := range keys {
		node, err := getParam(key, *n)

		// If err is nil, the node has been found
		if err == nil {
			return *node
		}
	}
	panic(fmt.Sprintf("no keys matching provided keys: %s", strings.Join(keys, ", ")))
}

func (n *Node) UpdateLineNumbers(lineNum int) {
	n.LineNum = lineNum
	for i := range n.Params {
		n.Params[i].UpdateLineNumbers(lineNum)
	}
}

func (n *Node) String() string {

	switch n.Type {

	case LIST:
		var s string
		for i, param := range n.Params {
			if i < len(n.Params)-1 {
				s += fmt.Sprintf("%s, ", param.String())
			} else {
				s += param.String()
			}
		}
		return fmt.Sprintf("(%s)", s)

	case STRING:
		doubleQuoteLiteral := "\""
		return fmt.Sprintf("%s%s%s", doubleQuoteLiteral, n.Value, doubleQuoteLiteral)

	case FUNCTION:
		funcKeyword := tokens.FUNCTION_TOKEN.Literal
		openParen := tokens.OPEN_PAREN_TOKEN.Literal
		closedParen := tokens.CLOSED_PAREN_TOKEN.Literal
		openCurlyBracket := tokens.OPEN_CURLY_BRACKET_TOKEN.Literal
		closedCurlyBracket := tokens.CLOSED_CURLY_BRACKET_TOKEN.Literal

		s := fmt.Sprintf("%s%s", funcKeyword, openParen)

		parameterStrings := []string{}
		for _, param := range n.GetParam(LIST).Params {
			parameterStrings = append(parameterStrings, param.String())
		}
		s += strings.Join(parameterStrings, ",")

		s += closedParen

		s += fmt.Sprintf("%s...%s", openCurlyBracket, closedCurlyBracket)

		return s

	case BUILTIN_FUNCTION:
		return fmt.Sprintf("<built-in function %s>", n.Value)

	case MONAD:
		s := "Monad{"
		if len(n.Params) == 1 {
			s += n.Params[0].String()
		}
		s += "}"
		return s

	default:
		// NUMBER, BOOLEAN
		return n.Value
	}
}

func (n Node) Equals(other Node) bool {
	/*
		The line number is not checked because two values could be equal but defined on different lines.
		Also, the line number is for debugging/error handling, so that value is not relevant here.

		If the object type and string representation match, we can assume the objects are equal.

		Do not use this method for testing.
	*/
	if n.Type != other.Type {
		return false
	}

	if n.Value != other.Value {
		return false
	}

	if len(n.Params) != len(other.Params) {
		return false
	}

	for i := 0; i < len(n.Params); i++ {
		selfParam := n.Params[i]
		otherParam := other.Params[i]

		if !selfParam.Equals(otherParam) {
			return false
		}
	}
	return true
}

func (n *Node) Length() (*Node, error) {
	var length int
	switch n.Type {
	case LIST:
		length = len(n.Params)
	case STRING:
		length = len(n.Value)
	default:
		return nil, utils.CreateError(n.LineNum, "Type %s does not have a length", n.Type)
	}
	return CreateNumber(n.LineNum, fmt.Sprint(length)).Ptr(), nil
}

func (n Node) Ptr() *Node {
	return &n
}

func getParam(key string, n Node) (*Node, error) {
	if stmtIndices, stmtOk := indexMap[n.Type]; stmtOk {
		if paramIndex, paramOk := stmtIndices[key]; paramOk {
			return &n.Params[paramIndex], nil
		}
		panic(fmt.Sprintf("invalid parameter: %s", key))
	}
	panic(fmt.Sprintf("invalid statement type: %s", n.Type))
}

const (

	// Statements
	STMTS                 = "Statements"
	RETURN_VALUE          = "ReturnValue"
	ASSIGN_STMT           = "Assign"
	TRUE_BRANCH           = "TrueBranch"
	FALSE_BRANCH          = "FalseBranch"
	CONDITION             = "Condition"
	BLOCK_STATEMENTS      = "BlockStatements"
	WHILE_LOOP            = "WhileLoop"
	WHILE_LOOP_CONDITION  = "WhileLoopCondition"
	WHILE_LOOP_STATEMENTS = "WhileLoopStatements"
	BREAK                 = "Break"
	CONTINUE              = "Continue"

	// Expressions
	EXPR                   = "Expression"
	UNARY_EXPR             = "UnaryExpression"
	OPERATOR               = "Operator"
	CALL_PARAMS            = "FunctionCallParameters"
	FUNCTION               = "Function"
	FUNCTION_CALL          = "FunctionCall"
	BIN_EXPR               = "BinaryExpression"
	LEFT                   = "Left"
	RIGHT                  = "Right"
	ASSIGN_STMT_IDENTIFIER = "Identifier"
	WHEN                   = "When"
	WHEN_VALUE             = "WhenValue"
	WHEN_CASES             = "WhenCases"
	WHEN_CASES_DEFAULT     = "WhenCasesDefault"
	CASE                   = "Case"
	CASE_VALUE             = "CaseValue"
	CASE_STMTS             = "CaseStatements"
	FOR_LOOP               = "ForLoop"
	FOR_LOOP_ELEMENT       = "ForLoopElement"

	// Factors
	NUMBER           = "Number"
	STRING           = "String"
	BOOLEAN          = "Boolean"
	IDENTIFIER       = "Identifier"
	BUILTIN_VARIABLE = "BuiltinVariable"
	BUILTIN_FUNCTION = "BuiltinFunction"
	LIST             = "List"
	MONAD            = "Monad"
	MONAD_VALUE      = "MonadValue"
)

/*
Node parameters are stored in a list of Nodes (Node.Params). Originally, Node.Params was a map[string]Node,
but this implementation had to change to accomodate multiple parameters of the same type. For example, with print
statements, I would use the node type as the key, but duplicate keys can't be stored in a map, which is a problem
for printing multiple nodes of the same time (e.g., "print(1, 2, 3);"). Additionally, preserving the order
of parameters is important, which maps do not do.

Instead, I moved to a system where the node parameters are an array of nodes, but I didn't want people to have to
remember what index stores a specific parameter, so indexMap allows for finding parameters by keys. The node type
and parameter name are used in concert to find the correct parameter in Node.Params. "getParam" performs this
search: "Node.Params[indexMap[NODE_TYPE][PARAM_NAME]]".

For example, let's say you have an ASSIGN_STMT Node object and want to find the variable name (ASSIGN_STMT
nodes contain a variable name (IDENTIFIER) and a value/expression to assign to that variable name (EXPR)).
To find the variable name, "getParam" searches "indexMap" with "indexMap[ASSIGN_STMT][IDENTIFIER]]",
which is 0 and the index in Node.Params where the variable name should be stored.
*/
var indexMap = map[string]map[string]int{
	ASSIGN_STMT: {
		IDENTIFIER: 0,
		EXPR:       1,
	},
	BIN_EXPR: {
		LEFT:     0,
		OPERATOR: 1,
		RIGHT:    2,
	},
	UNARY_EXPR: {
		OPERATOR: 0,
		EXPR:     1,
	},
	FUNCTION: {
		LIST:  0,
		STMTS: 1,
	},
	FUNCTION_CALL: {
		CALL_PARAMS: 0,
		FUNCTION:    1,
		IDENTIFIER:  1,
	},
	WHEN: {
		WHEN_VALUE:         0,
		WHEN_CASES:         1,
		WHEN_CASES_DEFAULT: 2,
	},
	CASE: {
		CASE_VALUE: 0,
		CASE_STMTS: 1,
	},
	FOR_LOOP: {
		FOR_LOOP_ELEMENT: 0,
		LIST:             1,
		BLOCK_STATEMENTS: 2,
	},
	WHILE_LOOP: {
		WHILE_LOOP_CONDITION:  0,
		WHILE_LOOP_STATEMENTS: 1,
	},
	MONAD: {
		MONAD_VALUE: 0,
	},
}

func CreateTokenNode(token tokens.Token) Node {
	return Node{Type: token.Type, Value: token.Literal, LineNum: token.LineNumber}
}

func CreateNumber(lineNum int, value string) Node {
	return Node{Type: NUMBER, Value: value, LineNum: lineNum}
}

func CreateBoolean(lineNum int, value string) Node {
	return Node{Type: BOOLEAN, Value: value, LineNum: lineNum}
}

func CreateBooleanTrue(lineNum int) Node {
	return CreateBoolean(lineNum, tokens.TRUE_TOKEN.Literal)
}

func CreateBooleanFalse(lineNum int) Node {
	return CreateBoolean(lineNum, tokens.FALSE_TOKEN.Literal)
}

func CreateString(lineNum int, literal string, parameters []Node) Node {
	// Strings that contain interpolation
	return Node{Type: STRING, Value: literal, Params: parameters, LineNum: lineNum}
}

func CreateRawString(lineNum int, literal string) Node {
	// Strings with no interpolation
	return CreateString(lineNum, literal, []Node{})
}

func CreateIdentifier(lineNum int, name string) Node {
	return Node{Type: IDENTIFIER, Value: name, LineNum: lineNum}
}

func CreateBuiltinVariableIdentifier(lineNum int, name string) Node {
	return Node{Type: BUILTIN_VARIABLE, Value: name, LineNum: lineNum}
}

func CreateBuiltinFunctionIdentifier(lineNum int, name string) Node {
	return Node{Type: BUILTIN_FUNCTION, Value: name, LineNum: lineNum}
}

func CreateList(lineNum int, parameters []Node) Node {
	return Node{Type: LIST, Params: parameters, LineNum: lineNum}
}

func CreateUnaryExpression(operator tokens.Token, expression Node) Node {
	return Node{
		Type:    UNARY_EXPR,
		LineNum: operator.LineNumber,
		Params: []Node{
			CreateTokenNode(operator), // Operator
			expression,                // Expression
		},
	}
}

func CreateBinaryExpression(left Node, op tokens.Token, right Node) Node {
	return Node{
		Type:    BIN_EXPR,
		LineNum: left.LineNum,
		Params: []Node{
			left,                // Left Expression
			CreateTokenNode(op), // Operator
			right,               // Right Expression
		},
	}
}

func CreateAssignmentNode(left, right Node) Node {
	return Node{
		Type:    ASSIGN_STMT,
		LineNum: left.LineNum,
		Params: []Node{
			left,  // Identifier(s) to assign values to
			right, // Value(s) assigned to those identifier(s)
		},
	}
}

func CreateFunction(lineNum int, parameters []Node, statements Node) Node {
	return Node{
		Type:    FUNCTION,
		LineNum: lineNum,
		Params: []Node{
			{Type: LIST, Params: parameters, LineNum: lineNum},
			statements,
		},
	}
}

func CreateFunctionCall(lineNum int, function Node, callParams []Node) Node {
	return Node{
		Type:    FUNCTION_CALL,
		LineNum: lineNum,
		Params: []Node{
			{Type: CALL_PARAMS, Params: callParams, LineNum: lineNum},
			function,
		},
	}
}

func CreateBlockStatementReturnValue(linenum int, statement *Node) Node {
	/*
	 Keep this method because the name, CreateBlockStatementReturnValue, makes it clear that this is what
	 block statements return. And if I never change what block statements return in the future, I will
	 only have to change this method.
	*/
	return CreateMonad(linenum, statement)
}

func CreateBlockStatements(statements []Node) Node {
	return Node{
		Type:   BLOCK_STATEMENTS,
		Params: statements,
	}
}

func CreateWhenNode(lineNum int, expression Node, caseNodes []Node, elseStatements Node) Node {
	return Node{
		Type:    WHEN,
		LineNum: lineNum,
		Params: []Node{
			expression,
			{Type: WHEN_CASES, LineNum: lineNum, Params: caseNodes},
			elseStatements,
		},
	}
}

func CreateCaseNode(lineNum int, expression Node, statements Node) Node {
	return Node{
		Type:    CASE,
		LineNum: lineNum,
		Params: []Node{
			expression,
			statements,
		},
	}
}

func CreateForLoop(lineNum int, elementPlaceholder Node, list Node, statements Node) Node {
	return Node{
		Type:    FOR_LOOP,
		LineNum: lineNum,
		Params: []Node{
			elementPlaceholder,
			list,
			statements,
		},
	}
}

func CreateWhileLoop(lineNum int, conditionExpression Node, statements Node) Node {
	return Node{
		Type:    WHILE_LOOP,
		LineNum: lineNum,
		Params: []Node{
			conditionExpression,
			statements,
		},
	}
}

func CreateBreakStatement(lineNum int) Node {
	return Node{Value: tokens.BREAK_TOKEN.Literal, Type: BREAK, LineNum: lineNum}
}

func CreateContinueStatement(lineNum int) Node {
	return Node{Value: tokens.CONTINUE_TOKEN.Literal, Type: CONTINUE, LineNum: lineNum}
}

func CreateMonad(lineNum int, value *Node) Node {

	params := []Node{}

	if value != nil {

		// If the value is already a monad, simply return the monad
		if value.Type == MONAD {
			return *value
		}
		params = append(params, *value)
	}

	return Node{
		Type:    MONAD,
		LineNum: lineNum,
		Params:  params,
	}
}
