package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/api"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "tokenize":
		input := readFile(os.Args[2])
		tokenizeCommand(&input)
	case "parse":
		input := readFile(os.Args[2])
		parseCommand(&input)
	case "evaluate":
		input := readFile(os.Args[2])
		evaluateCommand(&input)
	case "run":
		input := readFile(os.Args[2])
		runCommand(&input)
	}

}

// MARK: - Commands

func runCommand(input *string) {
	tokens, tokenizeErrors := api.Tokenize(input)
	if len(tokenizeErrors) > 0 {
		for _, err := range tokenizeErrors {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(65)
	}
	statements, parseError := api.ParseStmts(&tokens)
	if parseError != nil {
		fmt.Fprintln(os.Stderr, parseError)
		os.Exit(65)
	}
	api.Exec(&statements)
}

func evaluateCommand(input *string) {
	tokens, _ := api.Tokenize(input)
	expression := api.ParseExpr(&tokens)
	value, evalError := api.Eval(&expression)

	if evalError != nil {
		fmt.Fprintln(os.Stderr, evalError)
		os.Exit(70)
	}

	if value == nil {
		fmt.Println("nil")
	} else {
		fmt.Println(value)
	}
}

func parseCommand(input *string) {
	tokens, _ := api.Tokenize(input)
	expr := api.ParseExpr(&tokens)
	fmt.Println(expr)
}

func tokenizeCommand(input *string) {
	tokens, tokenizeErrors := api.Tokenize(input)
	if len(tokenizeErrors) > 0 {
		for _, err := range tokenizeErrors {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	for _, token := range tokens {
		fmt.Println(token)
	}
	if len(tokenizeErrors) > 0 {
		os.Exit(65)
	}
}

// MARK: - Helpers

func readFile(filename string) string {
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	return string(fileContents)
}
