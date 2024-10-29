package main

import (
	"fmt"
	"reflect"
)

type Expr interface {
	String() string
}

type LiteralExpr struct {
	value any
}
func (le LiteralExpr) String() string {
	if le.value == nil {
		return "nil"
	} else if reflect.TypeOf(le.value).Kind() == reflect.Float64 {
		return Float64ToString(le.value.(float64))
	} else {
		return fmt.Sprint(le.value)
	}
}

type GroupingExpr struct {
	expr Expr
}
func (ge GroupingExpr) String() string {
	return fmt.Sprintf("(group %v)", ge.expr)
}

func parse(tokens *[]Token) Expr {
	scanner := scanner{tokens: tokens, position: 0}
	return expression(&scanner)
}

type scanner struct {
	tokens *[]Token
	position int
}

func grouping(scanner *scanner) GroupingExpr {
	if match(scanner, LeftParen) {
		expr := expression(scanner)
		consume(scanner, RightParen)
		return GroupingExpr{expr: expr}
	}
	panic("?!?!?")
}

func literal(scanner *scanner) LiteralExpr {
	token := (*scanner.tokens)[scanner.position]
	scanner.position += 1
	if token.Type == True {
		return LiteralExpr{value: true}
	} else if token.Type == False {
		return LiteralExpr{value: false}
	} else if token.Type == Nil {
		return LiteralExpr{value: nil}
	} else if token.Type == Number {
		return LiteralExpr{value: token.Literal}
	} else if token.Type == String {
		return LiteralExpr{value: token.Literal}
	}
	panic("! literal")
}

func expression(scanner *scanner) Expr {
	if check(scanner, LeftParen) {
		return grouping(scanner)
	} else {
		return literal(scanner)
	}
}

func match(scanner *scanner, tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if check(scanner, tokenType) {
			advance(scanner)
			return true;
		}
	}
	return false;
}

func check(scanner *scanner, tokenType TokenType) bool {
	if (*scanner.tokens)[scanner.position].Type == EOF {
		return false
	}
	return (*scanner.tokens)[scanner.position].Type == tokenType
}

func advance(scanner *scanner) Token {
	if (*scanner.tokens)[scanner.position].Type != EOF {
		scanner.position += 1
	}
	return (*scanner.tokens)[scanner.position - 1]
}

func consume(scanner *scanner, tokenType TokenType) Token {
	if check(scanner, tokenType) {
		return advance(scanner)
	}
	panic("unexp")
}
