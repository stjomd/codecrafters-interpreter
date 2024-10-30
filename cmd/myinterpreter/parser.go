package main

import (
	"fmt"
	"os"
)

func parse(tokens *[]Token) Expr {
	parser := parser{tokens: tokens, position: 0}
	return parser.expression()
}

type parser struct {
	tokens *[]Token
	position int
}

// MARK: - Grammar rules

func (p *parser) expression() Expr {
	return p.equality()
}

func (p *parser) equality() Expr {
	var expr Expr = p.comparison()
	for p.match(EqualEqual, BangEqual) {
		operation := p.previous()
		right := p.comparison()
		expr = BinaryExpr{left: expr, operation: operation, right: right}
	}
	return expr
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
		p.consume(RightParen, "Expect ')'")
		return GroupingExpr{expr: expr}
	}
	// Wrongful state
	fmt.Fprintf(os.Stderr, "[line %d] Error at '%v': Expect expression.\n", p.peek().Line, p.peek().Lexeme)
	os.Exit(65)
	panic("!")
}

// MARK: - Helpers

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

func (p *parser) consume(tokenType TokenType, errorMessage string) Token {
	if p.check(tokenType) {
		return p.advance()
	}
	fmt.Fprintf(os.Stderr, "[line %d] Error at '%v': %s.\n", p.peek().Line, p.peek().Lexeme, errorMessage)
	os.Exit(65)
	panic("!")
}

func (p *parser) peek() Token {
	return (*p.tokens)[p.position]
}

func (p *parser) previous() Token {
	return (*p.tokens)[p.position - 1]
}
