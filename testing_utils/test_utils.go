package testing_utils

import (
	"boomerang/node"
	"fmt"
)

func AssertNodesEqual(expectedNodes []node.Node, actualNodes []node.Node) bool {
	if len(expectedNodes) != len(actualNodes) {
		fmt.Printf("Expected length: %d, Actual length: %d\n", len(expectedNodes), len(actualNodes))
	}
	for i := range expectedNodes {
		expected := expectedNodes[i]
		actual := actualNodes[i]
		if !AssertNodeEqual(expected, actual) {
			return false
		}
	}
	return true
}

func AssertNodeEqual(expected node.Node, actual node.Node) bool {
	if expected.Type != actual.Type {
		fmt.Printf("Expected type: %s, Actual type: %s\n", expected.Type, actual.Type)
		return false
	}

	if expected.Value != actual.Value {
		fmt.Printf("Expected value: %s, Actual value: %s\n", expected.Value, actual.Value)
		return false
	}

	keys := make([]string, 0, len(expected.Params))
	for k := range expected.Params {
		keys = append(keys, k)
	}

	for _, key := range keys {
		expectedParamNode := expected.GetParam(key)
		actualParamNode := actual.GetParam(key)
		return AssertNodeEqual(expectedParamNode, actualParamNode)
	}

	return true
}
