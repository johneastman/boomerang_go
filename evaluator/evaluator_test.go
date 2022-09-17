package evaluator

import (
	"boomerang/node"
	"boomerang/testing_utils"
	"testing"
)

func TestNumber(t *testing.T) {
	ast := []node.Node{
		{Type: "Number", Value: "5"},
	}
	evaluator := New(ast)

	actualResults := evaluator.Evaluate()
	expectedResults := []node.Node{
		{Type: "Number", Value: "5"},
	}

	testing_utils.AssertNodesEqual(expectedResults, actualResults)
}

func TestBinaryExpression(t *testing.T) {
	ast := []node.Node{
		{
			Type: "BinaryExpression",
			Params: map[string]node.Node{
				"left":     {Type: "Number", Value: "1"},
				"right":    {Type: "Number", Value: "1"},
				"operator": {Type: "PLUS", Value: "+"},
			},
		},
	}
	evaluator := New(ast)

	actualResults := evaluator.Evaluate()
	expectedResults := []node.Node{
		{Type: "Number", Value: "2"},
	}
	testing_utils.AssertNodesEqual(actualResults, expectedResults)
}
