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
	SLASH
	STAR
	// Single- or double-character tokens
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL
	// Literals
	STRING
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
	case SLASH:
		return "SLASH"
	case STAR:
		return "STAR"
	case BANG:
		return "BANG"
	case BANG_EQUAL:
		return "BANG_EQUAL"
	case EQUAL:
		return "EQUAL"
	case EQUAL_EQUAL:
		return "EQUAL_EQUAL"
	case GREATER:
		return "GREATER"
	case GREATER_EQUAL:
		return "GREATER_EQUAL"
	case LESS:
		return "LESS"
	case LESS_EQUAL:
		return "LESS_EQUAL"
	case STRING:
		return "STRING"
	case EOF:
		return "EOF"
	}
	return "?"
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
		// MARK: Single-character tokens
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
		case '/':
			var next, peekError = peek(&runes, i + 1)
			if peekError == nil && next == '/' {
				i = skipUntil(&runes, '\n', i + 1)
				line++
			} else {
				tokens = append(tokens, Token{Type: SLASH, Lexeme: string(character), Literal: "null"})
			}
		case ';':
			tokens = append(tokens, Token{Type: SEMICOLON, Lexeme: string(character), Literal: "null"})
		case '*':
			tokens = append(tokens, Token{Type: STAR, Lexeme: string(character), Literal: "null"})
		// MARK: Single- or double-character tokens
		case '!':
			var token, skip = lookahead(&runes, i, '=', BANG_EQUAL, BANG)
			tokens = append(tokens, token)
			i += skip
		case '=':
			var token, skip = lookahead(&runes, i, '=', EQUAL_EQUAL, EQUAL)
			tokens = append(tokens, token)
			i += skip
		case '>':
			var token, skip = lookahead(&runes, i, '=', GREATER_EQUAL, GREATER)
			tokens = append(tokens, token)
			i += skip
		case '<':
			var token, skip = lookahead(&runes, i, '=', LESS_EQUAL, LESS)
			tokens = append(tokens, token)
			i += skip
		// MARK: Literals
		case '"':
			index := skipUntilOnLine(&runes, '"', i + 1)
			if (index >= len(runes) || runes[index] == '\n') {
				message := fmt.Sprintf("[line %v] Error: Unterminated string.", line)
				errs = append(errs, errors.New(message))
			} else {
				lexeme, literal := string(runes[i:index+1]), string(runes[i+1:index])
				tokens = append(tokens, Token{Type: STRING, Lexeme: lexeme, Literal: literal})
			}
			i = index
		// MARK: Miscellaneous
		case '\n':
			line++
		case ' ', '\t':
			continue
		default:
			var message = fmt.Sprintf("[line %v] Error: Unexpected character: %v", line, string(character))
			errs = append(errs, errors.New(message))
		}
	}
	tokens = append(tokens, Token{Type: EOF, Lexeme: "", Literal: "null"})

	return tokens, errs
}

func skipUntil(input *[]rune, match rune, startPosition int) int {
	var slice = *input
	var i = startPosition
	for ; i < len(slice); i++ {
		if slice[i] == match {
			break
		}
	}
	return i;
}

func skipUntilOnLine(input *[]rune, match rune, startPosition int) int {
	var slice = *input
	var i = startPosition
	for ; i < len(slice); i++ {
		if slice[i] == match {
			break
		} else if slice[i] == '\n' {
			break
		}
	}
	return i;
}

// Looks ahead one character and, if it matches the `match` argument, returns a token of type `tokenIfMatch`. Otherwise
// returns a token of type `tokenIfNoMatch`.
func lookahead(input *[]rune, position int, match rune, tokenIfMatch TokenType, tokenIfNoMatch TokenType) (Token, int) {
	var character = (*input)[position]
	var next, peekError = peek(input, position + 1)
	if peekError == nil && next == match {
		return Token{Type: tokenIfMatch, Lexeme: string(character) + string(next), Literal: "null"}, 1
	} else {
		return Token{Type: tokenIfNoMatch, Lexeme: string(character), Literal: "null"}, 0
	}
}

func peek(input *[]rune, position int) (rune, error) {
	if (position > 0 && position < len(*input)) {
		return (*input)[position], nil
	}
	return 0, errors.New("index out of bounds")
}
