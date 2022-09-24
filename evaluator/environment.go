package evaluator

import (
	"boomerang/node"
	"fmt"
	"math"
)

type environment struct {
	env map[string]node.Node
}

func CreateEnvironment() environment {
	env := map[string]node.Node{}
	env["pi"] = node.CreateNumber(fmt.Sprintf("%v", math.Pi))
	env["len"] = node.CreateBuiltinFunction(node.BUILTIN_LEN)

	return environment{env: env}
}

func (e *environment) SetIdentifier(key string, value node.Node) {
	e.env[key] = value
}

func (e *environment) GetIdentifier(key string) node.Node {
	if value, ok := e.env[key]; ok {
		return value
	}
	panic(fmt.Sprintf("Undefined variable: %s", key))
}
