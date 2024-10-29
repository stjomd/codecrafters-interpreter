package main

import (
	"fmt"
	"reflect"
)

func parse(tokens *[]Token) Expr {
	parser := parser{tokens: tokens, position: 0}
	return parser.expression()
}

// MARK: - Expressions

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

type UnaryExpr struct {
	operation Token
	expr Expr
}
func (ue UnaryExpr) String() string {
	return fmt.Sprintf("(%v %v)", ue.operation.Lexeme, ue.expr)
}

type BinaryExpr struct {
	left Expr
	operation Token
	right Expr
}
func (be BinaryExpr) String() string {
	return fmt.Sprintf("(%v %v %v)", be.operation.Lexeme, be.left, be.right)
}

// MARK: - Parser methods

type parser struct {
	tokens *[]Token
	position int
}

func (p *parser) comparison() Expr {
	var expr Expr = p.term()
	for p.match(Less, LessEqual, Greater, GreaterEqual) {
		operation := p.previous()
		right := p.term()
		expr = BinaryExpr{left: expr, operation: operation, right: right}
	}
	return expr
}

func (p *parser) term() Expr {
	var expr Expr = p.factor()
	for p.match(Plus, Minus) {
		operation := p.previous()
		right := p.factor()
		expr = BinaryExpr{left: expr, operation: operation, right: right}
	}
	return expr
}

func (p *parser) factor() Expr {
	var expr Expr = p.unary()
	for p.match(Slash, Star) {
		operation := p.previous()
		right := p.unary()
		expr = BinaryExpr{left: expr, operation: operation, right: right}
	}
	return expr
}

func (p *parser) unary() Expr {
	if p.match(Bang, Minus) {
		operation := p.previous()
		expr := p.unary()
		return UnaryExpr{operation: operation, expr: expr}
	}
	return p.primary()
}

func (p *parser) primary() Expr {
	if p.match(True) {
		return LiteralExpr{value: true}
	} else if p.match(False) {
		return LiteralExpr{value: false}
	} else if p.match(Nil) {
		return LiteralExpr{value: nil}
	} else if p.match(Number, String) {
		return LiteralExpr{value: p.previous().Literal}
	} else if p.match(LeftParen) {
		expr := p.expression()
		p.consume(RightParen)
		return GroupingExpr{expr: expr}
	}
	panic("unexp literal: " + p.peek().String())
}

func (p *parser) expression() Expr {
	return p.comparison()
}

// MARK: Helpers

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
	if p.peek().Type == EOF {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *parser) advance() Token {
	if p.peek().Type != EOF {
		p.position += 1
	}
	return p.previous()
}

func (p *parser) consume(tokenType TokenType) Token {
	if p.check(tokenType) {
		return p.advance()
	}
	panic("unexp")
}

func (p *parser) peek() Token {
	return (*p.tokens)[p.position]
}

func (p *parser) previous() Token {
	return (*p.tokens)[p.position - 1]
}
