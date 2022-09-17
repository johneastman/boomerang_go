package node

import "fmt"

const (
	NUMBER                 = "Number"
	BIN_EXPR               = "BinaryExpression"
	EXPR                   = "Expression"
	IDENTIFIER             = "Identifier"
	OPERATOR               = "Operator"
	PARAMETER              = "Parameter"
	BIN_EXPR_LEFT          = "Left"
	BIN_EXPR_RIGHT         = "Right"
	UNARY_EXPR             = "UnaryExpression"
	PRINT_STMT             = "PrintStatement"
	ASSIGN_STMT            = "Assign"
	ASSIGN_STMT_IDENTIFIER = "Identifier"
)

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
}

type Node struct {
	Type   string
	Value  string
	Params []Node
}

func (n *Node) GetParam(key string) Node {
	if stmt_indices, stmt_ok := indexMap[n.Type]; stmt_ok {
		if param_index, param_ok := stmt_indices[key]; param_ok {
			return n.Params[param_index]
		}
		panic(fmt.Sprintf("Invalid parameter: %s", key))
	}
	panic(fmt.Sprintf("Invalid statement type: %s", n.Type))
}

func (n *Node) String() string {
	return fmt.Sprintf("Node(Type: %s, Value: %s)", n.Type, n.Value)
}
