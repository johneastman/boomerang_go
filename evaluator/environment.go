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

	// Builtin variables
	env["pi"] = node.CreateNumber(0, fmt.Sprintf("%v", math.Pi))

	// Builtin functions
	env["len"] = node.CreateBuiltinFunction(node.BUILTIN_LEN, 0)

	/*
		I originally wanted "unwrap" to be implemented in pure Boomerang code, but because custom functions
		return a list and the purpose of unwrap is to extract the return value from that list, this implementation
		needs to be a builtin method.
	*/
	env["unwrap"] = node.CreateBuiltinFunction(node.BUILTIN_UNWRAP, 0)

	return environment{env: env}
}

func (e *environment) SetIdentifier(key string, value node.Node) {
	e.env[key] = value
}

func (e *environment) GetIdentifier(key node.Node) (*node.Node, error) {
	identifierName := key.Value
	if value, ok := e.env[identifierName]; ok {
		value.LineNum = key.LineNum
		return &value, nil
	}
	return nil, utils.CreateError(key.LineNum, "undefined identifier: %s", identifierName)
}
