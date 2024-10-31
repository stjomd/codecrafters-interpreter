package api

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/spec"
)

func Parse(tokens *[]spec.Token) spec.Expr {
	parser := parser{tokens: tokens, position: 0}
	return parser.expression()
}

type parser struct {
	tokens *[]spec.Token
	position int
}

// MARK: - Grammar rules

func (p *parser) expression() spec.Expr {
	return p.equality()
}

func (p *parser) equality() spec.Expr {
	var expr spec.Expr = p.comparison()
	for p.match(spec.EqualEqual, spec.BangEqual) {
		operation := p.previous()
		right := p.comparison()
		expr = spec.BinaryExpr{Left: expr, Opt: operation, Right: right}
	}
	return expr
}

func (p *parser) comparison() spec.Expr {
	var expr spec.Expr = p.term()
	for p.match(spec.Less, spec.LessEqual, spec.Greater, spec.GreaterEqual) {
		operation := p.previous()
		right := p.term()
		expr = spec.BinaryExpr{Left: expr, Opt: operation, Right: right}
	}
	return expr
}

func (p *parser) term() spec.Expr {
	var expr spec.Expr = p.factor()
	for p.match(spec.Plus, spec.Minus) {
		operation := p.previous()
		right := p.factor()
		expr = spec.BinaryExpr{Left: expr, Opt: operation, Right: right}
	}
	return expr
}

func (p *parser) factor() spec.Expr {
	var expr spec.Expr = p.unary()
	for p.match(spec.Slash, spec.Star) {
		operation := p.previous()
		right := p.unary()
		expr = spec.BinaryExpr{Left: expr, Opt: operation, Right: right}
	}
	return expr
}

func (p *parser) unary() spec.Expr {
	if p.match(spec.Bang, spec.Minus) {
		operation := p.previous()
		expr := p.unary()
		return spec.UnaryExpr{Opt: operation, Expr: expr}
	}
	return p.primary()
}

func (p *parser) primary() spec.Expr {
	if p.match(spec.True) {
		return spec.LiteralExpr{Value: true}
	} else if p.match(spec.False) {
		return spec.LiteralExpr{Value: false}
	} else if p.match(spec.Nil) {
		return spec.LiteralExpr{Value: nil}
	} else if p.match(spec.Number, spec.String) {
		return spec.LiteralExpr{Value: p.previous().Literal}
	} else if p.match(spec.LeftParen) {
		expr := p.expression()
		p.consume(spec.RightParen, "Expect ')'")
		return spec.GroupingExpr{Expr: expr}
	}
	// Wrongful state
	fmt.Fprintf(os.Stderr, "[line %d] Error at '%v': Expect expression.\n", p.peek().Line, p.peek().Lexeme)
	os.Exit(65)
	panic("!")
}

// MARK: - Helpers

func (p *parser) match(tokenTypes ...spec.TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true;
		}
	}
	return false;
}

func (p *parser) check(tokenType spec.TokenType) bool {
	if p.peek().Type == spec.EOF {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *parser) advance() spec.Token {
	if p.peek().Type != spec.EOF {
		p.position += 1
	}
	return p.previous()
}

func (p *parser) consume(tokenType spec.TokenType, errorMessage string) spec.Token {
	if p.check(tokenType) {
		return p.advance()
	}
	fmt.Fprintf(os.Stderr, "[line %d] Error at '%v': %s.\n", p.peek().Line, p.peek().Lexeme, errorMessage)
	os.Exit(65)
	panic("!")
}

func (p *parser) peek() spec.Token {
	return (*p.tokens)[p.position]
}

func (p *parser) previous() spec.Token {
	return (*p.tokens)[p.position - 1]
}
