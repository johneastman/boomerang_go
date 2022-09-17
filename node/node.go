package node

import "fmt"

const (
	NUMBER                 = "Number"
	BIN_EXPR               = "BinaryExpression"
	EXPR                   = "Expression"
	IDENTIFIER             = "Identifier"
	OPERATOR               = "Operator"
	BIN_EXPR_LEFT          = "Left"
	BIN_EXPR_RIGHT         = "Right"
	UNARY_EXPR             = "UnaryExpression"
	ASSIGN_STMT            = "Assign"
	ASSIGN_STMT_IDENTIFIER = "Identifier"
)

type Node struct {
	Type   string
	Value  string
	Params map[string]Node
}

func (n *Node) GetParam(key string) Node {
	node, ok := n.Params[key]
	if !ok {
		panic(fmt.Sprintf("Key not in node params: %s", key))
	}
	return node
}

func (n *Node) String() string {
	return fmt.Sprintf("Node(Type: %s, Value: %s)", n.Type, n.Value)
}
