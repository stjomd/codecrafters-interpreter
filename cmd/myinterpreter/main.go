package main

import (
	"fmt"
	"os"
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
	}
}

// MARK: - Commands

func parseCommand(input *string) {
	tokens, _ := tokenize(input)
	expr, errors := parse(&tokens)
	fmt.Println(expr)
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(65)
	}
}

func tokenizeCommand(input *string) {
	tokens, tokenizeErrors := tokenize(input)
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
