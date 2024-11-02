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
	input := readFile(os.Args[2])

	switch command {
	case "tokenize":
		tokenizeCommand(&input)
	case "parse":
		parseCommand(&input)
	case "evaluate":
		evaluateCommand(&input)
	case "run":
		runCommand(&input)
	}

}

// MARK: - Commands

func runCommand(input *string) {
	tokens, tokenizeErrors := api.Tokenize(input)
	handleErrors(tokenizeErrors, 65)
	statements, parseError := api.ParseStmts(&tokens)
	handleError(parseError, 65)
	execError := api.Exec(&statements)
	handleError(execError, 70)
}

func evaluateCommand(input *string) {
	tokens, tokenizeErrors := api.Tokenize(input)
	handleErrors(tokenizeErrors, 65)
	expr, parseError := api.ParseExpr(&tokens)
	handleError(parseError, 65)
	value, evalError := api.EvalWithoutEnv(&expr)
	handleError(evalError, 70)
	if value == nil {
		fmt.Println("nil")
	} else {
		fmt.Println(value)
	}
}

func parseCommand(input *string) {
	tokens, tokenizeErrors := api.Tokenize(input)
	handleErrors(tokenizeErrors, 65)
	expr, parseError := api.ParseExpr(&tokens)
	handleError(parseError, 65)
	fmt.Println(expr)
}

func tokenizeCommand(input *string) {
	tokens, tokenizeErrors := api.Tokenize(input)
	for _, token := range tokens {
		fmt.Println(token)
	}
	handleErrors(tokenizeErrors, 65)
}

// MARK: - Helpers

func handleErrors(errors []error, exitCode int) {
	if len(errors) == 0 { return }
	for _, err := range errors {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(exitCode)
}

func handleError(err error, exitCode int) {
	if err == nil { return }
	fmt.Fprintln(os.Stderr, err)
	os.Exit(exitCode)
}

func readFile(filename string) string {
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	return string(fileContents)
}
