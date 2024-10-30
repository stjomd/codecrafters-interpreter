package scan

import (
	"fmt"
	"reflect"
)

// MARK: - Token types

type TokenType int
const (
	// Single-character tokens
	LeftParen TokenType = iota
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star
	// Single- or double-character tokens
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual
	// Literals
	Identifier
	String
	Number
	// Keywords
	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While
	// No-character tokens
	EOF
)

func (tt TokenType) String() string {
	switch tt {
	case LeftParen:
		return "LEFT_PAREN"
	case RightParen:
		return "RIGHT_PAREN"
	case LeftBrace:
		return "LEFT_BRACE"
	case RightBrace:
		return "RIGHT_BRACE"
	case Comma:
		return "COMMA"
	case Dot:
		return "DOT"
	case Minus:
		return "MINUS"
	case Plus:
		return "PLUS"
	case Semicolon:
		return "SEMICOLON"
	case Slash:
		return "SLASH"
	case Star:
		return "STAR"
	case Bang:
		return "BANG"
	case BangEqual:
		return "BANG_EQUAL"
	case Equal:
		return "EQUAL"
	case EqualEqual:
		return "EQUAL_EQUAL"
	case Greater:
		return "GREATER"
	case GreaterEqual:
		return "GREATER_EQUAL"
	case Less:
		return "LESS"
	case LessEqual:
		return "LESS_EQUAL"
	case Identifier:
		return "IDENTIFIER"
	case String:
		return "STRING"
	case Number:
		return "NUMBER"
	case And:
		return "AND"
	case Class:
		return "CLASS"
	case Else:
		return "ELSE"
	case False:
		return "FALSE"
	case Fun:
		return "FUN"
	case For:
		return "FOR"
	case If:
		return "IF"
	case Nil:
		return "NIL"
	case Or:
		return "OR"
	case Print:
		return "PRINT"
	case Return:
		return "RETURN"
	case Super:
		return "SUPER"
	case This:
		return "THIS"
	case True:
		return "TRUE"
	case Var:
		return "VAR"
	case While:
		return "WHILE"
	case EOF:
		return "EOF"
	}
	return "?"
}

var keywords = map[string]TokenType {
	"and": And,
	"class": Class,
	"else": Else,
	"false": False,
	"fun": Fun,
	"for": For,
	"if": If,
	"nil": Nil,
	"or": Or,
	"print": Print,
	"return": Return,
	"super": Super,
	"this": This,
	"true": True,
	"var": Var,
	"while": While,
}

var singleCharTokens = map[rune]TokenType {
	'(': LeftParen,
	')': RightParen,
	'{': LeftBrace,
	'}': RightBrace,
	',': Comma,
	'.': Dot,
	'-': Minus,
	'+': Plus,
	'/': Slash,
	';': Semicolon,
	'*': Star,
}

// MARK: - Token

type Token struct {
	Type TokenType
	Lexeme string
	Literal any
	Line uint64
}
func (token Token) String() string {
	literalString := ""
	if (token.Literal == nil) {
		literalString = "null"
	} else if reflect.TypeOf(token.Literal).Kind() == reflect.Float64 {
		literalString = Float64ToString(token.Literal.(float64))
	} else {
		literalString = fmt.Sprintf("%v", token.Literal)
	}
	return fmt.Sprintf("%v %v %v", token.Type.String(), token.Lexeme, literalString)
}

func Float64ToString(number float64) string {
	literalString := fmt.Sprintf("%g", number)
	if number == float64(int(number)) {
		literalString = literalString + ".0"
	}
	return literalString
}
