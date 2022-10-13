package evaluator

import (
	"boomerang/node"
	"boomerang/utils"
)

type environment struct {
	identifiers map[string]node.Node
	parentEnv   *environment
}

func CreateEnvironment(pEnv *environment) environment {
	env := environment{identifiers: map[string]node.Node{}, parentEnv: pEnv}
	return env
}

func (e *environment) SetIdentifier(key string, value node.Node) {
	e.identifiers[key] = value
}

func (e *environment) GetIdentifier(key node.Node) (*node.Node, error) {
	identifierName := key.Value
	env := e

	for env != nil {
		if value, ok := env.identifiers[identifierName]; ok {
			value.LineNum = key.LineNum
			return &value, nil
		}
		env = env.parentEnv
	}
	return nil, utils.CreateError(key.LineNum, "undefined identifier: %s", identifierName)
}
