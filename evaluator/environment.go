package evaluator

import (
	"boomerang/node"
	"boomerang/utils"
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

	/*
		I originally wanted "unwrap" to be implemented in pure Boomerang code, but because custom functions
		return a list and the purpose of unwrap is to extract the return value from that list, this implementation
		needs to be a builtin method.
	*/
	env["unwrap"] = node.CreateBuiltinFunction(node.BUILTIN_UNWRAP)

	return environment{env: env}
}

func (e *environment) SetIdentifier(key string, value node.Node) {
	e.env[key] = value
}

func (e *environment) GetIdentifier(key node.Node) (*node.Node, error) {
	identifierName := key.Value
	if value, ok := e.env[identifierName]; ok {
		return &value, nil
	}
	return nil, utils.CreateError(key.LineNum, "undefined identifier: %s", identifierName)
}
