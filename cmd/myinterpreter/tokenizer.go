package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"unicode"
)

// MARK: - Tokens

type token struct {
	tType tokenType
	lexeme string
	literal any
}
func (token token) String() string {
	literalString := ""
	if (token.literal == nil) {
		literalString = "null"
	} else if reflect.TypeOf(token.literal).Kind() == reflect.Float64 {
		// if the underlying type of `literal` is float64, do some custom handling
		literalString = fmt.Sprintf("%g", token.literal)
		if token.literal == float64(int(token.literal.(float64))) {
			literalString = literalString + ".0"
		}
	} else {
		literalString = fmt.Sprintf("%v", token.literal)
	}
	return fmt.Sprintf("%v %v %v", token.tType.String(), token.lexeme, literalString)
}

// MARK: - Tokenizer function

func tokenize(input string) ([]token, []error) {
	var line uint64 = 1
	var tokens []token
	var errs []error
	runes := []rune(input)
	for i := 0; i < len(runes); i++ {
		char := runes[i]
		// MARK: Single-character tokens
		if char == '/' {
			// comment handling
			var next, peekError = peek(&runes, i + 1)
			if peekError == nil && next == '/' {
				i = skipUntil(&runes, i + 1, IsNewline)
				line++
			}
		} else if singleCharTokenType, isSingleCharToken := singleCharTokens[char]; isSingleCharToken {
			tokens = append(tokens, token{tType: singleCharTokenType, lexeme: string(char), literal: nil})
		// MARK: Single- or double-character tokens
		} else if char == '!' {
			i = handleSingleDoubleCharToken(&tokens, &runes, i, '=', BangEqual, Bang)
		} else if char == '=' {
			i = handleSingleDoubleCharToken(&tokens, &runes, i, '=', EqualEqual, Equal)
		} else if char == '>' {
			i = handleSingleDoubleCharToken(&tokens, &runes, i, '=', GreaterEqual, Greater)
		} else if char == '<' {
			i = handleSingleDoubleCharToken(&tokens, &runes, i, '=', LessEqual, Less)
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
		} else if unicode.IsLetter(char) || char == '_' {
			index := handleIdentifierAndKeyword(&tokens, &runes, i)
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
	tokens = append(tokens, token{tType: EOF, lexeme: "", literal: nil})

	return tokens, errs
}

// MARK: - Helper functions

// Identifier handling
func handleIdentifierAndKeyword(tokens *[]token, runes *[]rune, currentPosition int) int {
	slice := *runes
	index := skipUntil(runes, currentPosition + 1, IsIdentifierEnd)
	lexeme := string(slice[currentPosition:index])
	keywordTokenType, presentInKeywords := keywords[lexeme]
	if presentInKeywords {
		*tokens = append(*tokens, token{tType: keywordTokenType, lexeme: lexeme, literal: nil})
	} else {
		*tokens = append(*tokens, token{tType: Identifier, lexeme: lexeme, literal: nil})
	}
	return index
}

// Number handling
func handleNumber(tokens *[]token, runes *[]rune, currentPosition int) int {
	slice := *runes
	index := skipUntil(runes, currentPosition + 1, IsNumberEnd)
	lexeme := string(slice[currentPosition:index])
	literal, convError := strconv.ParseFloat(lexeme, 64)
	if convError != nil {
		panic("could not parse float")
	}
	*tokens = append(*tokens, token{tType: Number, lexeme: lexeme, literal: literal})
	return index
}

// String handling
func handleString(tokens *[]token, runes *[]rune, currentPosition int) (int, error) {
	slice := *runes
	index := skipUntil(runes, currentPosition + 1, IsStringEndOrNewline)
	if (index >= len(slice) || slice[index] == '\n') {
		//lint:ignore ST1005 spec requires capitalized message with period at the end
		return index, errors.New("Unterminated string.")
	} else {
		lexeme, literal := string(slice[currentPosition:index+1]), string(slice[currentPosition+1:index])
		*tokens = append(*tokens, token{tType: String, lexeme: lexeme, literal: literal})
	}
	return index, nil
}

// Single- and double-character token handling
func handleSingleDoubleCharToken(
	tokens *[]token, input *[]rune, position int, match rune, tokenIfMatch tokenType, tokenIfNoMatch tokenType,
) int {
	var newToken token;
	character := (*input)[position]
	next, peekError := peek(input, position + 1)
	if peekError == nil && next == match {
		newToken = token{tType: tokenIfMatch, lexeme: string(character) + string(next), literal: nil}
	} else {
		newToken = token{tType: tokenIfNoMatch, lexeme: string(character), literal: nil}
	}
	*tokens = append(*tokens, newToken)
	return position + len(newToken.lexeme) - 1
}

// MARK: Lookahead functions

var (
	IsNewline            = func(x rune) bool { return x == '\n' }
	IsStringEndOrNewline = func(x rune) bool { return x == '"' || x == '\n' }
	IsNumberEnd          = func(x rune) bool { return x != '.' && !unicode.IsDigit(x) }
	IsIdentifierEnd      = func(x rune) bool { return !unicode.IsLetter(x) && !unicode.IsDigit(x) && x != '_' }
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
