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

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	if command == "tokenize" {
		tokenizeCommand(&fileContents)
	}
}

func tokenizeCommand(fileContents *[]byte) {
	var fileString = string(*fileContents)
	var tokens, tokenizeErrors = tokenize(&fileString)
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
