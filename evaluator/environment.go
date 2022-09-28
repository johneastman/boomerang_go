package evaluator

import (
	"boomerang/node"
	"boomerang/utils"
)

type environment struct {
	env map[string]node.Node
}

func CreateEnvironment() environment {
	env := map[string]node.Node{}
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
