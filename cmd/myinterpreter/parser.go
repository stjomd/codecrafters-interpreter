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
	parser := parser{tokens: tokens, position: 0}
	return parser.expression()
}

// MARK: - Parser methods

type parser struct {
	tokens *[]Token
	position int
}

func (p *parser) grouping() GroupingExpr {
	if p.match(LeftParen) {
		expr := p.expression()
		p.consume(RightParen)
		return GroupingExpr{expr: expr}
	}
	panic("?!?!?")
}

func (p *parser) literal() LiteralExpr {
	token := (*p.tokens)[p.position]
	p.position += 1
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

func (p *parser) expression() Expr {
	if p.check(LeftParen) {
		return p.grouping()
	} else {
		return p.literal()
	}
}

func (p *parser) match(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true;
		}
	}
	return false;
}

func (p *parser) check(tokenType TokenType) bool {
	if (*p.tokens)[p.position].Type == EOF {
		return false
	}
	return (*p.tokens)[p.position].Type == tokenType
}

func (p *parser) advance() Token {
	if (*p.tokens)[p.position].Type != EOF {
		p.position += 1
	}
	return (*p.tokens)[p.position - 1]
}

func (p *parser) consume(tokenType TokenType) Token {
	if p.check(tokenType) {
		return p.advance()
	}
	panic("unexp")
}
