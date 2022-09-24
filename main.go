package main

import (
	"boomerang/evaluator"
	"boomerang/parser"
	"boomerang/tokens"
	"fmt"
	"log"
	"os"
)

func getSource(path string) string {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(fileContent)
}

func main() {
	source := getSource("source.bmg")
	tokenizer := tokens.New(source)

	parser, err := parser.New(tokenizer)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	statements, err := parser.Parse()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	eval := evaluator.New(*statements)
	_, err = eval.Evaluate()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
