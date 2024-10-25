package main

import (
	"errors"
	"fmt"
)

// MARK: - Token types
type TokenType int
const (
	// Single-character tokens
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
	// Single- or double-character tokens
	EQUAL
	EQUAL_EQUAL
	// No-character tokens
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
	case EQUAL:
		return "EQUAL"
	case EQUAL_EQUAL:
		return "EQUAL_EQUAL"
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
func tokenize(input string) ([]Token, []error) {
	var line uint64 = 1
	var tokens []Token
	var errs []error

	var runes = []rune(input)
	for i := 0; i < len(runes); i++ {
		var character = runes[i]
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
		case '=':
			var next, peekError = peek(runes, i + 1)
			if peekError == nil && next == '=' {
				tokens = append(tokens, Token{Type: EQUAL_EQUAL, Lexeme: string(character) + string(next), Literal: "null"})
				i++;
			} else {
				tokens = append(tokens, Token{Type: EQUAL, Lexeme: string(character), Literal: "null"})
			}
		case '\n':
			line++
		default:
			var message = fmt.Sprintf("[line %v] Error: Unexpected character: %v", line, string(character))
			errs = append(errs, errors.New(message))
		}
	}
	tokens = append(tokens, Token{Type: EOF, Lexeme: "", Literal: "null"})

	return tokens, errs
}

func peek(input []rune, position int) (rune, error) {
	if (position > 0 && position < len(input)) {
		return input[position], nil
	}
	return 0, errors.New("index out of bounds")
}
