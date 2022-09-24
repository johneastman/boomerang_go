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

	// Builtin bariables
	env["pi"] = node.CreateNumber(fmt.Sprintf("%v", math.Pi))

	// Builtin functions
	env["len"] = node.CreateBuiltinFunction(node.BUILTIN_LEN)
	env["unwrap"] = node.CreateBuiltinFunction(node.BUILTIN_UNWRAP)

	return environment{env: env}
}

func (e *environment) SetIdentifier(key string, value node.Node) {
	e.env[key] = value
}

func (e *environment) GetIdentifier(key string) (*node.Node, error) {
	if value, ok := e.env[key]; ok {
		return &value, nil
	}
	return nil, fmt.Errorf("Undefined variable: %s", key)
}
