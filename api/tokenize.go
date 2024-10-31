package api

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"

	"github.com/codecrafters-io/interpreter-starter-go/spec"
)

func Tokenize(input *string) ([]spec.Token, []error) {
	var line uint64 = 1
	var tokens []spec.Token
	var errs []error
	runes := []rune(*input)
	for i := 0; i < len(runes); i++ {
		char := runes[i]
		// MARK: Single-character tokens
		if singleCharTokenType, isSingleCharToken := spec.SingleCharTokens[char]; isSingleCharToken {
			// handle comments too
			next, peekError := peek(&runes, i + 1)
			if peekError == nil && char == '/' && next == '/' {
				i = skipUntil(&runes, i + 1, isNewline)
				line++
			} else {
				tokens = append(tokens, spec.Token{Type: singleCharTokenType, Lexeme: string(char), Literal: nil, Line: line})
			}
		// MARK: Single- or double-character tokens
		} else if char == '!' {
			i = handleSingleDoubleCharToken(&tokens, &runes, i, line, '=', spec.BangEqual, spec.Bang)
		} else if char == '=' {
			i = handleSingleDoubleCharToken(&tokens, &runes, i, line, '=', spec.EqualEqual, spec.Equal)
		} else if char == '>' {
			i = handleSingleDoubleCharToken(&tokens, &runes, i, line, '=', spec.GreaterEqual, spec.Greater)
		} else if char == '<' {
			i = handleSingleDoubleCharToken(&tokens, &runes, i, line, '=', spec.LessEqual, spec.Less)
		// MARK: Literals
		} else if char == '"' {
			index, err := handleString(&tokens, &runes, i, line)
			if err != nil {
				message := fmt.Sprintf("[line %v] Error: %s", line, err.Error())
				errs = append(errs, errors.New(message))
			}
			i = index
		} else if unicode.IsDigit(char) {
			index := handleNumber(&tokens, &runes, i, line)
			i = index - 1
		} else if unicode.IsLetter(char) || char == '_' {
			index := handleIdentifierAndKeyword(&tokens, &runes, i, line)
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
	tokens = append(tokens, spec.Token{Type: spec.EOF, Lexeme: "", Literal: nil, Line: line})

	return tokens, errs
}

// MARK: - Helper functions

// Identifier handling
func handleIdentifierAndKeyword(tokens *[]spec.Token, runes *[]rune, currentPosition int, line uint64) int {
	slice := *runes
	index := skipUntil(runes, currentPosition + 1, isIdentifierEnd)
	lexeme := string(slice[currentPosition:index])
	keywordTokenType, presentInKeywords := spec.Keywords[lexeme]
	if presentInKeywords {
		*tokens = append(*tokens, spec.Token{Type: keywordTokenType, Lexeme: lexeme, Literal: nil, Line: line})
	} else {
		*tokens = append(*tokens, spec.Token{Type: spec.Identifier, Lexeme: lexeme, Literal: nil, Line: line})
	}
	return index
}

// Number handling
func handleNumber(tokens *[]spec.Token, runes *[]rune, currentPosition int, line uint64) int {
	slice := *runes
	index := skipUntil(runes, currentPosition + 1, isNumberEnd)
	lexeme := string(slice[currentPosition:index])
	literal, convError := strconv.ParseFloat(lexeme, 64)
	if convError != nil {
		panic("could not parse float")
	}
	*tokens = append(*tokens, spec.Token{Type: spec.Number, Lexeme: lexeme, Literal: literal, Line: line})
	return index
}

// String handling
func handleString(tokens *[]spec.Token, runes *[]rune, currentPosition int, line uint64) (int, error) {
	slice := *runes
	index := skipUntil(runes, currentPosition + 1, isStringEndOrNewline)
	if (index >= len(slice) || slice[index] == '\n') {
		//lint:ignore ST1005 spec requires capitalized message with period at the end
		return index, errors.New("Unterminated string.")
	} else {
		lexeme, literal := string(slice[currentPosition:index+1]), string(slice[currentPosition+1:index])
		*tokens = append(*tokens, spec.Token{Type: spec.String, Lexeme: lexeme, Literal: literal, Line: line})
	}
	return index, nil
}

// Single- and double-character token handling
func handleSingleDoubleCharToken(
	tokens *[]spec.Token, input *[]rune, position int, line uint64,
	match rune, tokenIfMatch spec.TokenType, tokenIfNoMatch spec.TokenType,
) int {
	var newToken spec.Token;
	character := (*input)[position]
	next, peekError := peek(input, position + 1)
	if peekError == nil && next == match {
		newToken = spec.Token{Type: tokenIfMatch, Lexeme: string(character) + string(next), Literal: nil, Line: line}
	} else {
		newToken = spec.Token{Type: tokenIfNoMatch, Lexeme: string(character), Literal: nil, Line: line}
	}
	*tokens = append(*tokens, newToken)
	return position + len(newToken.Lexeme) - 1
}

// MARK: Lookahead functions

var (
	isNewline            = func(x rune) bool { return x == '\n' }
	isStringEndOrNewline = func(x rune) bool { return x == '"' || x == '\n' }
	isNumberEnd          = func(x rune) bool { return x != '.' && !unicode.IsDigit(x) }
	isIdentifierEnd      = func(x rune) bool { return !unicode.IsLetter(x) && !unicode.IsDigit(x) && x != '_' }
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
