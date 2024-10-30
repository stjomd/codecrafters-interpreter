package parse

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/scan"
)

func Parse(tokens *[]scan.Token) Expr {
	parser := parser{tokens: tokens, position: 0}
	return parser.expression()
}

type parser struct {
	tokens *[]scan.Token
	position int
}

// MARK: - Grammar rules

func (p *parser) expression() Expr {
	return p.equality()
}

func (p *parser) equality() Expr {
	var expr Expr = p.comparison()
	for p.match(scan.EqualEqual, scan.BangEqual) {
		operation := p.previous()
		right := p.comparison()
		expr = BinaryExpr{left: expr, operation: operation, right: right}
	}
	return expr
}

func (p *parser) comparison() Expr {
	var expr Expr = p.term()
	for p.match(scan.Less, scan.LessEqual, scan.Greater, scan.GreaterEqual) {
		operation := p.previous()
		right := p.term()
		expr = BinaryExpr{left: expr, operation: operation, right: right}
	}
	return expr
}

func (p *parser) term() Expr {
	var expr Expr = p.factor()
	for p.match(scan.Plus, scan.Minus) {
		operation := p.previous()
		right := p.factor()
		expr = BinaryExpr{left: expr, operation: operation, right: right}
	}
	return expr
}

func (p *parser) factor() Expr {
	var expr Expr = p.unary()
	for p.match(scan.Slash, scan.Star) {
		operation := p.previous()
		right := p.unary()
		expr = BinaryExpr{left: expr, operation: operation, right: right}
	}
	return expr
}

func (p *parser) unary() Expr {
	if p.match(scan.Bang, scan.Minus) {
		operation := p.previous()
		expr := p.unary()
		return UnaryExpr{operation: operation, expr: expr}
	}
	return p.primary()
}

func (p *parser) primary() Expr {
	if p.match(scan.True) {
		return LiteralExpr{value: true}
	} else if p.match(scan.False) {
		return LiteralExpr{value: false}
	} else if p.match(scan.Nil) {
		return LiteralExpr{value: nil}
	} else if p.match(scan.Number, scan.String) {
		return LiteralExpr{value: p.previous().Literal}
	} else if p.match(scan.LeftParen) {
		expr := p.expression()
		p.consume(scan.RightParen, "Expect ')'")
		return GroupingExpr{expr: expr}
	}
	// Wrongful state
	fmt.Fprintf(os.Stderr, "[line %d] Error at '%v': Expect expression.\n", p.peek().Line, p.peek().Lexeme)
	os.Exit(65)
	panic("!")
}

// MARK: - Helpers

func (p *parser) match(tokenTypes ...scan.TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true;
		}
	}
	return false;
}

func (p *parser) check(tokenType scan.TokenType) bool {
	if p.peek().Type == scan.EOF {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *parser) advance() scan.Token {
	if p.peek().Type != scan.EOF {
		p.position += 1
	}
	return p.previous()
}

func (p *parser) consume(tokenType scan.TokenType, errorMessage string) scan.Token {
	if p.check(tokenType) {
		return p.advance()
	}
	fmt.Fprintf(os.Stderr, "[line %d] Error at '%v': %s.\n", p.peek().Line, p.peek().Lexeme, errorMessage)
	os.Exit(65)
	panic("!")
}

func (p *parser) peek() scan.Token {
	return (*p.tokens)[p.position]
}

func (p *parser) previous() scan.Token {
	return (*p.tokens)[p.position - 1]
}
