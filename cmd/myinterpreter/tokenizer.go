package main

// MARK: - Token types
type TokenType int
const (
	LEFT_PAREN = iota
	RIGHT_PAREN
	EOF
)
func (tokenType TokenType) String() string {
	switch tokenType {
	case LEFT_PAREN:
		return "LEFT_PAREN"
	case RIGHT_PAREN:
		return "RIGHT_PAREN"
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
	var tokens []Token
	for _, character := range input {
		switch character {
		case '(':
			tokens = append(tokens, Token{Type: LEFT_PAREN, Lexeme: string(character), Literal: "null"})
		case ')':
			tokens = append(tokens, Token{Type: RIGHT_PAREN, Lexeme: string(character), Literal: "null"})
		}
	}
	tokens = append(tokens, Token{Type: EOF, Lexeme: "", Literal: "null"})
	return tokens
}
