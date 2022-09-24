package main

import (
	"boomerang/evaluator"
	"boomerang/parser"
	"boomerang/tokens"
	"fmt"
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

	parser, err := parser.New(tokenizer)
	if err != nil {
		fmt.Println(err.Error())
	}

	statements, err := parser.Parse()
	if err != nil {
		fmt.Println(err.Error())
	}

	eval := evaluator.New(*statements)
	_, err = eval.Evaluate()
	if err != nil {
		fmt.Println(err.Error())
	}
}
