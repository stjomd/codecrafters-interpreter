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
		if stmt, err := parser.declaration(); err == nil {
			statements = append(statements, stmt)
		} else {
			return nil, err
		}
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
		if expression, err := p.expression(); err == nil {
			expr = expression
		} else {
			return nil, err
		}
	}
	if _, err := p.consume(spec.Semicolon, "Expect ';' after variable declaration"); err != nil {
		return nil, err
	}
	return spec.DeclareStmt{Identifier: identifier, Expr: expr}, nil
}

func (p *parser) statement() (spec.Stmt, error) {
	if p.match(spec.Print) {
		return p.printStatement()
	} else if p.match(spec.LeftBrace) {
		return p.blockStatement()
	} else if p.match(spec.If) {
		return p.ifStatement()
	} else if p.match(spec.While) {
		return p.whileStatement()
	} else if p.match(spec.For) {
		return p.forStatement()
	}
	return p.expressionStatement()
}

func (p *parser) expressionStatement() (spec.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(spec.Semicolon, "Expect ';' after value"); err != nil {
		return nil, err
	}
	return spec.ExprStmt{Expr: expr}, nil
}

func (p *parser) printStatement() (spec.Stmt, error) {
	expr, err := p.expression()
	if err != nil { return nil, err }
	if _, err := p.consume(spec.Semicolon, "Expect ';' after value"); err != nil {
		return nil, err
	}
	return spec.PrintStmt{Expr: expr}, nil
}

func (p *parser) blockStatement() (spec.Stmt, error) {
	var statements []spec.Stmt
	for !p.check(spec.RightBrace) && (p.peek().Type != spec.EOF) {
		if stmt, err := p.declaration(); err == nil {
			statements = append(statements, stmt)
		} else {
			return nil, err
		}
	}
	if _, err := p.consume(spec.RightBrace, "Expect '}' after block"); err != nil {
		return nil, err
	}
	return spec.BlockStmt{Statements: statements}, nil
}

func (p *parser) ifStatement() (spec.Stmt, error) {
	// condition
	if _, err := p.consume(spec.LeftParen, "Expect '(' after 'if'."); err != nil {
		return nil, err
	}
	condition, conditionError := p.expression()
	if conditionError != nil {
		return nil, conditionError
	}
	if _, err := p.consume(spec.RightParen, "Expect ')' after if condition."); err != nil {
		return nil, err
	}
	// then branch
	thenBranch, thenError := p.statement()
	if thenError != nil {
		return nil, thenError
	}
	stmt := spec.IfStmt{Condition: condition, Then: thenBranch}
	// (optional) else branch
	if p.match(spec.Else) {
		if elseBranch, err := p.statement(); err == nil {
			stmt.Else = elseBranch
		} else {
			return nil, err
		}
	}
	return stmt, nil
}

func (p *parser) whileStatement() (spec.Stmt, error) {
	// condition
	if _, err := p.consume(spec.LeftParen, "Expect '(' after 'while'."); err != nil {
		return nil, err
	}
	condition, conditionError := p.expression()
	if conditionError != nil {
		return nil, conditionError
	}
	if _, err := p.consume(spec.RightParen, "Expect ')' after condition."); err != nil {
		return nil, err
	}
	// body
	body, bodyError := p.statement()
	if bodyError != nil {
		return nil, bodyError
	}
	return spec.WhileStmt{Condition: condition, Body: body}, nil
}

func (p *parser) forStatement() (spec.Stmt, error) {
	if _, err := p.consume(spec.LeftParen, "Expect '(' after 'for'."); err != nil {
		return nil, err
	}
	// head - initializer
	var init spec.Stmt
	if p.match(spec.Semicolon) {
		init = nil
	} else if p.match(spec.Var) {
		if initializer, err := p.varDeclaration(); err == nil {
			init = initializer
		} else {
			return nil, err
		}
	} else {
		if initializer, err := p.expressionStatement(); err == nil {
			init = initializer
		} else {
			return nil, err
		}
	}
	// head - condition
	var cond spec.Expr
	if !p.check(spec.Semicolon) {
		if condition, err := p.expression(); err == nil {
			cond = condition
		} else {
			return nil, err
		}
	}
	p.consume(spec.Semicolon, "Expect ';' after loop condition.");
	// head - increment
	var incr spec.Expr
	if !p.check(spec.RightParen) {
		if increment, err := p.expression(); err == nil {
			incr = increment
		} else {
			return nil, err
		}
	}
	p.consume(spec.RightParen, "Expect ')' after for clauses.");
	// body
	body, bodyError := p.statement()
	if bodyError != nil {
		return nil, bodyError
	}
	// desugar to a while loop
	return forLoopAsStatement(init, cond, incr, body), nil
}

// MARK: - Expressions

func (p *parser) expression() (spec.Expr, error) {
	return p.assignment()
}

func (p *parser) assignment() (spec.Expr, error) {
	expr, exprError := p.or()
	if exprError != nil {
		return nil, exprError
	}
	if p.match(spec.Equal) {
		value, valueError := p.assignment()
		if valueError != nil {
			return nil, valueError
		}
		if areTypesEqual(expr, spec.VariableExpr{}) {
			identifier := expr.(spec.VariableExpr).Identifier
			return spec.AssignmentExpr{Identifier: identifier, Expr: value}, nil
		}
	}
	return expr, nil
}

func (p *parser) or() (spec.Expr, error) {
	expr, exprError := p.and()
	if exprError != nil {
		return nil, exprError
	}
	for p.match(spec.Or) {
		operator := p.previous()
		if rightExpr, err := p.and(); err == nil {
			expr = spec.LogicalExpr{Left: expr, Opt: operator, Right: rightExpr}
		} else {
			return nil, err
		}
	}
	return expr, nil
}

func (p *parser) and() (spec.Expr, error) {
	expr, exprError := p.equality()
	if exprError != nil {
		return nil, exprError
	}
	for p.match(spec.And) {
		operator := p.previous()
		if rightExpr, err := p.equality(); err == nil {
			expr = spec.LogicalExpr{Left: expr, Opt: operator, Right: rightExpr}
		} else {
			return nil, err
		}
	}
	return expr, nil
}

func (p *parser) equality() (spec.Expr, error) {
	expr, exprError := p.comparison()
	if exprError != nil {
		return nil, exprError
	}
	for p.match(spec.EqualEqual, spec.BangEqual) {
		operator := p.previous()
		if rightExpr, err := p.comparison(); err == nil {
			expr = spec.BinaryExpr{Left: expr, Opt: operator, Right: rightExpr}
		} else {
			return nil, err
		}
	}
	return expr, nil
}

func (p *parser) comparison() (spec.Expr, error) {
	expr, exprError := p.term()
	if exprError != nil {
		return nil, exprError
	}
	for p.match(spec.Less, spec.LessEqual, spec.Greater, spec.GreaterEqual) {
		operator := p.previous()
		if rightExpr, err := p.term(); err == nil {
			expr = spec.BinaryExpr{Left: expr, Opt: operator, Right: rightExpr}
		} else {
			return nil, err
		}
	}
	return expr, nil
}

func (p *parser) term() (spec.Expr, error) {
	expr, exprError := p.factor()
	if exprError != nil {
		return nil, exprError
	}
	for p.match(spec.Plus, spec.Minus) {
		operator := p.previous()
		if rightExpr, err := p.factor(); err == nil {
			expr = spec.BinaryExpr{Left: expr, Opt: operator, Right: rightExpr}
		} else {
			return nil, err
		}
	}
	return expr, nil
}

func (p *parser) factor() (spec.Expr, error) {
	expr, exprError := p.unary()
	if exprError != nil {
		return nil, exprError
	}
	for p.match(spec.Slash, spec.Star) {
		operator := p.previous()
		if rightExpr, err := p.unary(); err == nil {
			expr = spec.BinaryExpr{Left: expr, Opt: operator, Right: rightExpr}
		} else {
			return nil, err
		}
	}
	return expr, nil
}

func (p *parser) unary() (spec.Expr, error) {
	if p.match(spec.Bang, spec.Minus) {
		operator := p.previous()
		if expr, err := p.unary(); err == nil {
			return spec.UnaryExpr{Opt: operator, Expr: expr}, nil
		} else {
			return nil, err
		}
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
		if exprError != nil {
			return nil, exprError
		}
		if _, err := p.consume(spec.RightParen, "Expect ')'"); err != nil {
			return nil, err
		}
		return spec.GroupingExpr{Expr: expr}, nil
	}
	message := fmt.Sprintf("[line %d] Error at '%v': Expect expression.", p.peek().Line, p.peek().Lexeme)
	return nil, errors.New(message)
}

// MARK: - Transformers

func forLoopAsStatement(init spec.Stmt, cond spec.Expr, incr spec.Expr, body spec.Stmt) spec.Stmt {
	// for (init; cond; incr) body
	// init; while cond { body; incr; }
	whileLoop := spec.WhileStmt{
		Condition: cond,
		Body: spec.BlockStmt{
			Statements: []spec.Stmt{
				body,
				spec.ExprStmt{Expr: incr},
			},
		},
	};
	var statements []spec.Stmt
	if init != nil {
		statements = append(statements, init)
	}
	statements = append(statements, whileLoop)
	return spec.BlockStmt{Statements: statements}
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
