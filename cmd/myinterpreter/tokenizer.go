package main

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
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
	NUMBER
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
	case NUMBER:
		return "NUMBER"
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
	runes := []rune(input)
	for i := 0; i < len(runes); i++ {
		char := runes[i]
		// MARK: Single-character tokens
		if char == '(' {
			tokens = append(tokens, Token{Type: LEFT_PAREN, Lexeme: string(char), Literal: "null"})
		} else if char == ')' {
			tokens = append(tokens, Token{Type: RIGHT_PAREN, Lexeme: string(char), Literal: "null"})
		} else if char == '{' {
			tokens = append(tokens, Token{Type: LEFT_BRACE, Lexeme: string(char), Literal: "null"})
		} else if char == '}' {
			tokens = append(tokens, Token{Type: RIGHT_BRACE, Lexeme: string(char), Literal: "null"})
		} else if char == ',' {
			tokens = append(tokens, Token{Type: COMMA, Lexeme: string(char), Literal: "null"})
		} else if char == '.' {
			tokens = append(tokens, Token{Type: DOT, Lexeme: string(char), Literal: "null"})
		} else if char == '-' {
			tokens = append(tokens, Token{Type: MINUS, Lexeme: string(char), Literal: "null"})
		} else if char == '+' {
			tokens = append(tokens, Token{Type: PLUS, Lexeme: string(char), Literal: "null"})
		} else if char == '/' {
			var next, peekError = peek(&runes, i + 1)
			if peekError == nil && next == '/' {
				i = skipUntil(&runes, i + 1, IS_NEWLINE)
				line++
			} else {
				tokens = append(tokens, Token{Type: SLASH, Lexeme: string(char), Literal: "null"})
			}
		} else if char == ';' {
			tokens = append(tokens, Token{Type: SEMICOLON, Lexeme: string(char), Literal: "null"})
		} else if char == '*' {
			tokens = append(tokens, Token{Type: STAR, Lexeme: string(char), Literal: "null"})
		// MARK: Single- or double-character tokens
		} else if char == '!' {
			token, size := handleSingleDoubleCharToken(&runes, i, '=', BANG_EQUAL, BANG)
			tokens = append(tokens, token)
			i += size - 1
		} else if char == '=' {
			token, size := handleSingleDoubleCharToken(&runes, i, '=', EQUAL_EQUAL, EQUAL)
			tokens = append(tokens, token)
			i += size - 1
		} else if char == '>' {
			token, size := handleSingleDoubleCharToken(&runes, i, '=', GREATER_EQUAL, GREATER)
			tokens = append(tokens, token)
			i += size - 1
		} else if char == '<' {
			token, size := handleSingleDoubleCharToken(&runes, i, '=', LESS_EQUAL, LESS)
			tokens = append(tokens, token)
			i += size - 1
		// MARK: Literals
		} else if char == '"' {
			index, err := handleString(&tokens, &runes, i)
			if err != nil {
				message := fmt.Sprintf("[line %v] Error: %s", line, err.Error())
				errs = append(errs, errors.New(message))
			}
			i = index
		} else if unicode.IsDigit(char) {
			index := handleNumber(&tokens, &runes, i)
			i = index - 1
		// MARK: Miscellaneous
		} else if char == '\n' {
			line++
		}	else if char == ' ' || char == '\t' {
			continue
		} else {
			message := fmt.Sprintf("[line %v] Error: Unexpected character: %v", line, string(char))
			errs = append(errs, errors.New(message))
		}
	}
	tokens = append(tokens, Token{Type: EOF, Lexeme: "", Literal: "null"})

	return tokens, errs
}

// MARK: - Helper functions

// Number handling
func handleNumber(tokens *[]Token, runes *[]rune, currentPosition int) int {
	slice := *runes
	index := skipUntil(runes, currentPosition + 1, IS_END_OF_NUMBER)
	lexeme := string(slice[currentPosition:index])
	literal, convError := strconv.ParseFloat(lexeme, 64)
	if convError != nil {
		panic("could not parse float")
	}
	stringLiteral := fmt.Sprintf("%g", literal)
	if literal == float64(int(literal)) {
		stringLiteral = stringLiteral + ".0"
	}
	*tokens = append(*tokens, Token{Type: NUMBER, Lexeme: lexeme, Literal: stringLiteral})
	return index
}

// String handling
func handleString(tokens *[]Token, runes *[]rune, currentPosition int) (int, error) {
	slice := *runes
	index := skipUntil(runes, currentPosition + 1, IS_STRING_END_OR_NEWLINE)
	if (index >= len(slice) || slice[index] == '\n') {
		//lint:ignore ST1005 spec requires capitalized message with period at the end
		return index, errors.New("Unterminated string.")
	} else {
		lexeme, literal := string(slice[currentPosition:index+1]), string(slice[currentPosition+1:index])
		*tokens = append(*tokens, Token{Type: STRING, Lexeme: lexeme, Literal: literal})
	}
	return index, nil
}

// Looks ahead one character and, if it matches the `match` argument, returns a token of type `tokenIfMatch`. Otherwise
// returns a token of type `tokenIfNoMatch`. Moreover, returns the size of the lexeme.
func handleSingleDoubleCharToken(
	input *[]rune, position int, match rune, tokenIfMatch TokenType, tokenIfNoMatch TokenType,
) (Token, int) {
	character := (*input)[position]
	next, peekError := peek(input, position + 1)
	if peekError == nil && next == match {
		return Token{Type: tokenIfMatch, Lexeme: string(character) + string(next), Literal: "null"}, 2
	} else {
		return Token{Type: tokenIfNoMatch, Lexeme: string(character), Literal: "null"}, 1
	}
}

// MARK: Lookahead functions

var (
	IS_NEWLINE               = func(x rune) bool { return x == '\n' }
	IS_STRING_END_OR_NEWLINE = func(x rune) bool { return x == '"' || x == '\n' }
	IS_END_OF_NUMBER         = func(x rune) bool { return x != '.' && !unicode.IsDigit(x) }
)
// Looks ahead, starting at the specified position, and until a specified condition is fulfiled or the end of input is
// reached, and returns the position.
func skipUntil(input *[]rune, startPosition int, condition func(rune) bool) int {
	slice, i := *input, startPosition
	for ; i < len(slice); i++ {
		if condition(slice[i]) { break }
	}
	return i
}

// Peeks at the rune at the specified position. Returns an error if the position is out of bounds; otherwise, returns
// the rune at that position.
func peek(input *[]rune, position int) (rune, error) {
	if (position > 0 && position < len(*input)) {
		return (*input)[position], nil
	}
	return 0, errors.New("index out of bounds")
}
