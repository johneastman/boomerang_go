package main

import (
	"boomerang/evaluator"
	"boomerang/parser"
	"boomerang/tokens"
	"io/ioutil"
	"log"
)

func getSource(path string) string {
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(fileContent)
}

func main() {
	source := getSource("source.bmg")
	tokenizer := tokens.New(source)

	parser := parser.New(tokenizer)
	statements := parser.Parse()

	eval := evaluator.New(statements)
	eval.Evaluate()
}
