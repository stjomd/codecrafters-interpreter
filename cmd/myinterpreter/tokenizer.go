package main

import (
	"fmt"
	"os"
)

// MARK: - Token types
type TokenType int
const (
	LEFT_PAREN = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	STAR
	EOF
)
func (tokenType TokenType) String() string {
	switch tokenType {
	case LEFT_PAREN:
		return "LEFT_PAREN"
	case RIGHT_PAREN:
		return "RIGHT_PAREN"
	case LEFT_BRACE:
		return "LEFT_BRACE"
	case RIGHT_BRACE:
		return "RIGHT_BRACE"
	case COMMA:
		return "COMMA"
	case DOT:
		return "DOT"
	case MINUS:
		return "MINUS"
	case PLUS:
		return "PLUS"
	case SEMICOLON:
		return "SEMICOLON"
	case STAR:
		return "STAR"
	case EOF:
		return "EOF"
	}
	return ""
}

// MARK: - Tokens
type Token struct {
	Type TokenType
	Lexeme string
	Literal string
}
func (token Token) String() string {
	return token.Type.String() + " " + token.Lexeme + " " + token.Literal
}

// MARK: - Tokenizer function
func Tokenize(input string) []Token {
	var line uint64 = 1
	var tokens []Token
	for _, character := range input {
		switch character {
		case '(':
			tokens = append(tokens, Token{Type: LEFT_PAREN, Lexeme: string(character), Literal: "null"})
		case ')':
			tokens = append(tokens, Token{Type: RIGHT_PAREN, Lexeme: string(character), Literal: "null"})
		case '{':
			tokens = append(tokens, Token{Type: LEFT_BRACE, Lexeme: string(character), Literal: "null"})
		case '}':
			tokens = append(tokens, Token{Type: RIGHT_BRACE, Lexeme: string(character), Literal: "null"})
		case ',':
			tokens = append(tokens, Token{Type: COMMA, Lexeme: string(character), Literal: "null"})
		case '.':
			tokens = append(tokens, Token{Type: DOT, Lexeme: string(character), Literal: "null"})
		case '-':
			tokens = append(tokens, Token{Type: MINUS, Lexeme: string(character), Literal: "null"})
		case '+':
			tokens = append(tokens, Token{Type: PLUS, Lexeme: string(character), Literal: "null"})
		case ';':
			tokens = append(tokens, Token{Type: SEMICOLON, Lexeme: string(character), Literal: "null"})
		case '*':
			tokens = append(tokens, Token{Type: STAR, Lexeme: string(character), Literal: "null"})
		case '\n':
			line++
		default:
			fmt.Fprintf(os.Stderr, "[line %v] Error: Unexpected character: %v\n", line, string(character))
		}
	}
	tokens = append(tokens, Token{Type: EOF, Lexeme: "", Literal: "null"})
	return tokens
}
