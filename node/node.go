package node

import (
	"boomerang/tokens"
	"fmt"
	"strings"
)

type Node struct {
	Type    string
	Value   string
	LineNum int
	Params  []Node
}

const (

	// Statements
	STMTS            = "Statements" // Super type
	RETURN_VALUE     = "ReturnValue"
	PRINT_STMT       = "PrintStatement"
	ASSIGN_STMT      = "Assign"
	IF_STMT          = "IfStatement"
	TRUE_BRANCH      = "TrueBranch"
	FALSE_BRANCH     = "FalseBranch"
	CONDITION        = "Condition"
	BLOCK_STATEMENTS = "BlockStatements"

	// Expressions
	EXPR                   = "Expression" // Super type
	UNARY_EXPR             = "UnaryExpression"
	OPERATOR               = "Operator"
	CALL_PARAMS            = "FunctionCallParameters"
	FUNCTION               = "Function"
	FUNCTION_CALL          = "FunctionCall"
	BIN_EXPR               = "BinaryExpression"
	LEFT                   = "Left"
	RIGHT                  = "Right"
	ASSIGN_STMT_IDENTIFIER = "Identifier"
	SWITCH                 = "Switch"
	SWITCH_VALUE           = "SwitchValue"
	SWITCH_CASES           = "SwitchCases"
	SWITCH_CASES_DEFAULT   = "SwitchCasesDefault"
	CASE                   = "Case"
	CASE_VALUE             = "CaseValue"
	CASE_STMTS             = "CaseStatements"

	// Factors
	NUMBER     = "Number"
	STRING     = "String"
	BOOLEAN    = "Boolean"
	IDENTIFIER = "Identifier"
	LIST       = "List"
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
	IF_STMT: {
		CONDITION:    0,
		TRUE_BRANCH:  1,
		FALSE_BRANCH: 2,
	},
	SWITCH: {
		SWITCH_VALUE:         0,
		SWITCH_CASES:         1,
		SWITCH_CASES_DEFAULT: 2,
	},
	CASE: {
		CASE_VALUE: 0,
		CASE_STMTS: 1,
	},
}

func (n *Node) ErrorDisplay() string {
	return fmt.Sprintf("%s (%#v)", n.Type, n.Value)
}

func (n *Node) GetParam(key string) Node {
	node, err := n.getParam(key)
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
		node, err := n.getParam(key)

		// If err is nil, the node has been found
		if err == nil {
			return *node
		}
	}
	panic(fmt.Sprintf("no keys matching provided keys: %s", strings.Join(keys, ", ")))
}

func (n *Node) getParam(key string) (*Node, error) {
	if stmtIndices, stmtOk := indexMap[n.Type]; stmtOk {
		if paramIndex, paramOk := stmtIndices[key]; paramOk {
			return &n.Params[paramIndex], nil
		}
		panic(fmt.Sprintf("invalid parameter: %s", key))
	}
	panic(fmt.Sprintf("invalid statement type: %s", n.Type))
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
		doubleQuoteLiteral := tokens.DOUBLE_QUOTE_TOKEN.Literal
		return fmt.Sprintf("%s%s%s", doubleQuoteLiteral, n.Value, doubleQuoteLiteral)

	default:
		// NUMBER, BOOLEAN
		return n.Value
	}
}

func (n Node) Ptr() *Node {
	return &n
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

func CreateList(lineNum int, parameters []Node) Node {
	return Node{Type: LIST, Params: parameters, LineNum: lineNum}
}

func CreatePrintStatement(lineNum int, params []Node) Node {
	return Node{
		Type:    PRINT_STMT,
		Params:  params,
		LineNum: lineNum,
	}
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

func CreateAssignmentStatement(lineNum int, variableName string, value Node) Node {
	return Node{
		Type:    ASSIGN_STMT,
		LineNum: lineNum,
		Params: []Node{
			{Type: IDENTIFIER, Value: variableName},
			value,
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

func CreateIfStatement(lineNum int, condition Node, trueStatements Node, falseStatements Node) Node {
	return Node{
		Type:    IF_STMT,
		LineNum: lineNum,
		Params: []Node{
			condition,
			trueStatements,
			falseStatements,
		},
	}
}

func CreateFunctionReturnValue(linenum int, statement *Node) Node {

	var parameters []Node

	if statement == nil {
		parameters = []Node{
			CreateBooleanFalse(linenum),
		}
	} else {
		parameters = []Node{
			CreateBooleanTrue(linenum),
			*statement,
		}
	}
	return CreateList(linenum, parameters)
}

func CreateBlockStatements(lineNum int, statements []Node) Node {
	return Node{
		Type:    BLOCK_STATEMENTS,
		LineNum: lineNum,
		Params:  statements,
	}
}

func CreateSwitchNode(lineNum int, expression Node, caseNodes []Node, elseStatements Node) Node {
	return Node{
		Type:    SWITCH,
		LineNum: lineNum,
		Params: []Node{
			expression,
			{Type: SWITCH_CASES, LineNum: lineNum, Params: caseNodes},
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
