package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"my_lang/evaluator"
	"my_lang/parser"
	"my_lang/tokens"
)

func getSource(path string) string {
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(fileContent)
}

func main() {
	source := getSource("source.txt")
	tokenizer := tokens.New(source)

	parser := parser.New(tokenizer)
	statements := parser.Parse()

	eval := evaluator.New(statements)
	results := eval.Evaluate()
	for _, result := range results {
		fmt.Println(result.Value)
	}
}
