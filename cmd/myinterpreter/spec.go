package main

// MARK: - Token types

type tokenType int
const (
	// Single-character tokens
	LeftParen tokenType = iota
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

func (tt tokenType) String() string {
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

var keywords = map[string]tokenType {
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

var singleCharTokens = map[rune]tokenType {
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
