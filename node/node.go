package node

import "fmt"

const (
	NUMBER         = "Number"
	BIN_EXPR       = "BinaryExpression"
	BIN_EXPR_LEFT  = "Left"
	BIN_EXPR_RIGHT = "Right"
	BIN_EXPR_OP    = "Operator"
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
