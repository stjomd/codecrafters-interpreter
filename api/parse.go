package api

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/codecrafters-io/interpreter-starter-go/spec"
)

func ParseExpr(tokens *[]spec.Token) (spec.Expr, error) {
	parser := parser{tokens: tokens, position: 0}
	return parser.expression()
}

func ParseStmts(tokens *[]spec.Token) ([]spec.Stmt, error) {
	var statements []spec.Stmt
	parser := parser{tokens: tokens, position: 0}
	for parser.peek().Type != spec.EOF {
		stmt, err := parser.declaration()
		if err != nil { return nil, err }
		statements = append(statements, stmt)
	}
	return statements, nil
}

type parser struct {
	tokens *[]spec.Token
	position int
}

// MARK: - Grammar rules
// MARK: Statements

func (p *parser) declaration() (spec.Stmt, error) {
	if (p.match(spec.Var)) {
		return p.varDeclaration() 
	}
	return p.statement()
}

func (p *parser) varDeclaration() (spec.Stmt, error) {
	identifier, consumeError := p.consume(spec.Identifier, "Expect variable name")
	if consumeError != nil { return nil, consumeError }
	var expr spec.Expr = spec.LiteralExpr{Value: nil}
	if p.match(spec.Equal) {
		expression, expressionError := p.expression()
		if expressionError != nil { return nil, expressionError }
		expr = expression
	}
	_, consumeError = p.consume(spec.Semicolon, "Expect ';' after variable declaration")
	if consumeError != nil { return nil, consumeError }
	return spec.DeclareStmt{Identifier: identifier, Expr: expr}, nil
}

func (p *parser) statement() (spec.Stmt, error) {
	if p.match(spec.Print) {
		return p.printStatement()
	} else if p.match(spec.LeftBrace) {
		return p.blockStatement()
	}
	return p.expressionStatement()
}

func (p *parser) expressionStatement() (spec.Stmt, error) {
	expr, err := p.expression()
	if err != nil { return nil, err }
	_, consumeError := p.consume(spec.Semicolon, "Expect ';' after value")
	if consumeError != nil { return nil, consumeError }
	return spec.ExprStmt{Expr: expr}, nil
}

func (p *parser) printStatement() (spec.Stmt, error) {
	expr, err := p.expression()
	if err != nil { return nil, err }
	_, consumeError := p.consume(spec.Semicolon, "Expect ';' after value")
	if consumeError != nil { return nil, consumeError }
	return spec.PrintStmt{Expr: expr}, nil
}

func (p *parser) blockStatement() (spec.Stmt, error) {
	var statements []spec.Stmt
	for !p.check(spec.RightBrace) && (p.peek().Type != spec.EOF) {
		stmt, stmtError := p.declaration()
		if stmtError != nil { return nil, stmtError }
		statements = append(statements, stmt)
	}
	_, consumeError := p.consume(spec.RightBrace, "Expect '}' after block")
	if consumeError != nil { return nil, consumeError }
	return spec.BlockStmt{Statements: statements}, nil
}

// MARK: Expressions

func (p *parser) expression() (spec.Expr, error) {
	return p.assignment()
}

func (p *parser) assignment() (spec.Expr, error) {
	expr, exprError := p.equality()
	if exprError != nil { return nil, exprError }
	if p.match(spec.Equal) {
		value, valueError := p.assignment()
		if valueError != nil { return nil, valueError }
		if areTypesEqual(expr, spec.VariableExpr{}) {
			identifier := expr.(spec.VariableExpr).Identifier
			return spec.AssignmentExpr{Identifier: identifier, Expr: value}, nil
		}
	}
	return expr, nil
}

func (p *parser) equality() (spec.Expr, error) {
	expr, exprError := p.comparison()
	if exprError != nil { return nil, exprError }
	for p.match(spec.EqualEqual, spec.BangEqual) {
		operation := p.previous()
		right, rightError := p.comparison()
		if rightError != nil { return nil, rightError }
		expr = spec.BinaryExpr{Left: expr, Opt: operation, Right: right}
	}
	return expr, nil
}

func (p *parser) comparison() (spec.Expr, error) {
	expr, exprError := p.term()
	if exprError != nil { return nil, exprError }
	for p.match(spec.Less, spec.LessEqual, spec.Greater, spec.GreaterEqual) {
		operation := p.previous()
		right, rightError := p.term()
		if rightError != nil { return nil, rightError }
		expr = spec.BinaryExpr{Left: expr, Opt: operation, Right: right}
	}
	return expr, nil
}

func (p *parser) term() (spec.Expr, error) {
	expr, exprError := p.factor()
	if exprError != nil { return nil, exprError }
	for p.match(spec.Plus, spec.Minus) {
		operation := p.previous()
		right, rightError := p.factor()
		if rightError != nil { return nil, rightError }
		expr = spec.BinaryExpr{Left: expr, Opt: operation, Right: right}
	}
	return expr, nil
}

func (p *parser) factor() (spec.Expr, error) {
	expr, exprError := p.unary()
	if exprError != nil { return nil, exprError }
	for p.match(spec.Slash, spec.Star) {
		operation := p.previous()
		right, rightError := p.unary()
		if rightError != nil { return nil, rightError }
		expr = spec.BinaryExpr{Left: expr, Opt: operation, Right: right}
	}
	return expr, nil
}

func (p *parser) unary() (spec.Expr, error) {
	if p.match(spec.Bang, spec.Minus) {
		operation := p.previous()
		expr, exprError := p.unary()
		if exprError != nil { return nil, exprError }
		return spec.UnaryExpr{Opt: operation, Expr: expr}, nil
	}
	return p.primary()
}

func (p *parser) primary() (spec.Expr, error) {
	if p.match(spec.Identifier) {
		return spec.VariableExpr{Identifier: p.previous()}, nil
	} else if p.match(spec.True) {
		return spec.LiteralExpr{Value: true}, nil
	} else if p.match(spec.False) {
		return spec.LiteralExpr{Value: false}, nil
	} else if p.match(spec.Nil) {
		return spec.LiteralExpr{Value: nil}, nil
	} else if p.match(spec.Number, spec.String) {
		return spec.LiteralExpr{Value: p.previous().Literal}, nil
	} else if p.match(spec.LeftParen) {
		expr, exprError := p.expression()
		if exprError != nil { return nil, exprError }
		_, consumeError := p.consume(spec.RightParen, "Expect ')'")
		if consumeError != nil { return nil, consumeError }
		return spec.GroupingExpr{Expr: expr}, nil
	}
	message := fmt.Sprintf("[line %d] Error at '%v': Expect expression.", p.peek().Line, p.peek().Lexeme)
	return nil, errors.New(message)
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

func (p *parser) consume(tokenType spec.TokenType, errorMessage string) (spec.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	message := fmt.Sprintf("[line %d] Error at '%v': %s.", p.peek().Line, p.peek().Lexeme, errorMessage)
	return spec.Token{}, errors.New(message)
}

func (p *parser) peek() spec.Token {
	return (*p.tokens)[p.position]
}

func (p *parser) previous() spec.Token {
	return (*p.tokens)[p.position - 1]
}

// MARK: - Crimes against humanity

func areTypesEqual(a any, b any) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}
