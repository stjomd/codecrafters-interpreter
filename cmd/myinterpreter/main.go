package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/scan"
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
	}

}

// MARK: - Commands

func evaluateCommand(input *string) {
	tokens, _ := scan.Tokenize(input)
	expr := parse(&tokens)
	value, evalError := expr.Eval()

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
	tokens, _ := scan.Tokenize(input)
	expr := parse(&tokens)
	fmt.Println(expr)
}

func tokenizeCommand(input *string) {
	tokens, tokenizeErrors := scan.Tokenize(input)
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
