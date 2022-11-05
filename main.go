package main

import (
	"boomerang/evaluator"
	"boomerang/parser"
	"boomerang/tokens"
	"boomerang/utils"
	"fmt"
)

func main() {
	source := utils.GetSource("source.bmg")
	tokenizer := tokens.NewTokenizer(source)

	parser, err := parser.NewParser(tokenizer)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	statements, err := parser.Parse()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	eval := evaluator.NewEvaluator(*statements)
	_, err = eval.Evaluate()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
