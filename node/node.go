package node

import (
	"boomerang/tokens"
	"fmt"
	"strings"
)

const (
	NUMBER                 = "Number"
	BIN_EXPR               = "BinaryExpression"
	STMTS                  = "Statements"
	EXPR                   = "Expression"
	IDENTIFIER             = "Identifier"
	OPERATOR               = "Operator"
	PARAMETER              = "Parameter"
	CALL_PARAMS            = "FunctionCallParameters"
	FUNCTION               = "Function"
	FUNCTION_CALL          = "FunctionCall"
	BIN_EXPR_LEFT          = "Left"
	BIN_EXPR_RIGHT         = "Right"
	UNARY_EXPR             = "UnaryExpression"
	PRINT_STMT             = "PrintStatement"
	ASSIGN_STMT            = "Assign"
	ASSIGN_STMT_IDENTIFIER = "Identifier"
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
		BIN_EXPR_LEFT:  0,
		OPERATOR:       1,
		BIN_EXPR_RIGHT: 2,
	},
	UNARY_EXPR: {
		OPERATOR: 0,
		EXPR:     1,
	},
	FUNCTION: {
		PARAMETER: 0,
		STMTS:     1,
	},
	FUNCTION_CALL: {
		CALL_PARAMS: 0,
		FUNCTION:    1,
	},
}

type Node struct {
	Type   string
	Value  string
	Params []Node
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
	panic(fmt.Sprintf("No keys matching provided keys: %s", strings.Join(keys, ", ")))
}

func (n *Node) getParam(key string) (*Node, error) {
	if stmt_indices, stmt_ok := indexMap[n.Type]; stmt_ok {
		if param_index, param_ok := stmt_indices[key]; param_ok {
			return &n.Params[param_index], nil
		}
		return nil, fmt.Errorf("invalid parameter: %s", key)
	}
	return nil, fmt.Errorf("invalid statement type: %s", n.Type)
}

func (n *Node) String() string {
	return fmt.Sprintf("Node(Type: %s, Value: %s)", n.Type, n.Value)
}

func CreateAssignmentStatement(name tokens.Token, value Node) Node {
	return Node{
		Type: ASSIGN_STMT,
		Params: []Node{
			{Type: name.Type, Value: name.Literal},
			value,
		},
	}
}

func CreatePrintStatement(params []Node) Node {
	return Node{
		Type:   PRINT_STMT,
		Params: params,
	}
}

func CreateUnaryExpression(operator tokens.Token, expression Node) Node {
	return Node{
		Type: UNARY_EXPR,
		Params: []Node{
			{Type: operator.Type, Value: operator.Literal}, // Operator
			expression, // Expression
		},
	}
}

func CreateIdentifier(name string) Node {
	return Node{Type: IDENTIFIER, Value: name}
}

func CreateNumber(value string) Node {
	return Node{Type: NUMBER, Value: value}
}

func CreateBinaryExpression(left Node, op tokens.Token, right Node) Node {
	return Node{
		Type: BIN_EXPR,
		Params: []Node{
			left,                               // Left Expression
			{Type: op.Type, Value: op.Literal}, // Operator
			right,                              // Right Expression
		},
	}
}

func CreateFunction(parameters []Node, statements []Node) Node {
	return Node{
		Type: FUNCTION,
		Params: []Node{
			{Type: PARAMETER, Params: parameters},
			{Type: STMTS, Params: statements},
		},
	}
}
