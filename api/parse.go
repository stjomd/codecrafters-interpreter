package api

import (
	"errors"
	"fmt"
	"math/rand"
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
	if (p.match(spec.Class)) {
		return p.classDesclaration()
	}
	if (p.match(spec.Fun)) {
		return p.funcDeclaration()
	}
	if (p.match(spec.Var)) {
		return p.varDeclaration() 
	}
	return p.statement()
}

func (p *parser) classDesclaration() (spec.Stmt, error) {
	name, nameError := p.consume(spec.Identifier, "Expect class name")
	if nameError != nil {
		return nil, nameError
	}

	var superclass *spec.VariableExpr
	if p.match(spec.Less) {
		superclassIdent, superclassError := p.consume(spec.Identifier, "Expect superclass name")
		if superclassError != nil {
			return nil, superclassError
		}
		superclass = &spec.VariableExpr{Identifier: superclassIdent, Occurrence: rand.Float64()}
	}

	if _, braceError := p.consume(spec.LeftBrace, "Expect '{' after function name"); braceError != nil {
		return nil, braceError
	}

	methods := []spec.FuncStmt{}
	for !p.check(spec.RightBrace) && !p.check(spec.EOF) {
		method, methodError := p.funcDeclaration()
		if methodError != nil {
			return nil, methodError
		}
		methods = append(methods, method.(spec.FuncStmt))
	}

	if _, braceError := p.consume(spec.RightBrace, "Expect '}' after function name"); braceError != nil {
		return nil, braceError
	}
	return spec.ClassStmt{Name: name, Methods: methods, Superclass: superclass}, nil
}

func (p *parser) funcDeclaration() (spec.Stmt, error) {
	name, nameError := p.consume(spec.Identifier, "Expect function name")
	if nameError != nil {
		return nil, nameError
	}
	if _, parenError := p.consume(spec.LeftParen, "Expect '(' after function name"); parenError != nil {
		return nil, parenError
	}
	params := []spec.Token{}
	if !p.check(spec.RightParen) {
		for next := true; next; next = p.match(spec.Comma) {
			if len(params) >= 255 {
				return nil, fmt.Errorf("can't have more than 255 parameters")
			}
			param, paramError := p.consume(spec.Identifier, "Expect parameter name")
			if paramError != nil {
				return nil, paramError
			}
			params = append(params, param)
		}
	}
	if _, parenError := p.consume(spec.RightParen, "Expect ')' after parameters"); parenError != nil {
		return nil, parenError
	}
	if _, braceError := p.consume(spec.LeftBrace, "Expect '{' before body"); braceError != nil {
		return nil, braceError
	}
	body, bodyError := p.blockStatement()
	if bodyError != nil {
		return nil, bodyError
	}
	return spec.FuncStmt{Name: name, Params: params, Body: body.(spec.BlockStmt).Statements}, nil
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
	} else if p.match(spec.Return) {
		return p.returnStatement()
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

func (p *parser) returnStatement() (spec.Stmt, error) {
	var keyword spec.Token = p.previous()
	var expr spec.Expr = nil;
	if !p.check(spec.Semicolon) {
		returnExpr, returnError := p.expression()
		if returnError != nil {
			return nil, returnError
		}
		expr = returnExpr
	}
	if _, scError := p.consume(spec.Semicolon, "Expect ';' after return value"); scError != nil {
		return nil, scError
	}
	return spec.ReturnStmt{Keyword: keyword, Expr: expr}, nil
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
		} else if areTypesEqual(expr, spec.GetExpr{}) {
			get := expr.(spec.GetExpr)
			return spec.SetExpr{Object: get.Object, Name: get.Name, Value: value}, nil
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
	return p.call()
}

func (p *parser) call() (spec.Expr, error) {
	expr, exprError := p.primary()
	if exprError != nil {
		return nil, exprError
	}
	for {
		if p.match(spec.LeftParen) {
			finishedCall, finishedCallError := p.finishCall(expr)
			if finishedCallError != nil {
				return nil, finishedCallError
			}
			expr = finishedCall
		} else if p.match(spec.Dot) {
			name, nameError := p.consume(spec.Identifier, "Expect property name after '.'")
			if nameError != nil {
				return nil, nameError
			}
			expr = spec.GetExpr{Object: expr, Name: name} 
		} else {
			break
		}
	}
	return expr, nil
}

func (p *parser) finishCall(callee spec.Expr) (spec.Expr, error) {
	args := []spec.Expr{};
	if !p.check(spec.RightParen) {
		for next := true; next; next = p.match(spec.Comma) {
			if len(args) >= 255 {
				return nil, fmt.Errorf("can't have more than 255 arguments")
			}
			arg, argError := p.expression()
			if argError != nil {
				return nil, argError
			}
			args = append(args, arg)
		}
	}
	paren, parenError := p.consume(spec.RightParen, "Expect ')' after arguments")
	if parenError != nil {
		return nil, parenError
	}
	return spec.CallExpr{Callee: callee, Paren: paren, Args: args}, nil
}

func (p *parser) primary() (spec.Expr, error) {
	if p.match(spec.Identifier) {
		return spec.VariableExpr{Identifier: p.previous(), Occurrence: rand.Float64()}, nil
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
	} else if p.match(spec.This) {
		return spec.ThisExpr{Keyword: p.previous()}, nil
	} else if p.match(spec.Super) {
		keyword := p.previous()
		if _, err := p.consume(spec.Dot, "Expect '.' after 'super'"); err != nil {
			return nil, err
		}
		method, methodErr := p.consume(spec.Identifier, "Expect superclass method name")
		if methodErr != nil {
			return nil, methodErr
		}
		return spec.SuperExpr{Keyword: keyword, Method: method}, nil
	}
	message := fmt.Sprintf("[line %d] Error at '%v': Expect expression.", p.peek().Line, p.peek().Lexeme)
	return nil, errors.New(message)
}

// MARK: - Transformers

func forLoopAsStatement(init spec.Stmt, cond spec.Expr, incr spec.Expr, body spec.Stmt) spec.Stmt {
	// for (init; cond; incr) body
	// init; while cond { body; incr; }
	whileLoopBody := []spec.Stmt{body}
	if incr != nil {
		incrStmt := spec.ExprStmt{Expr: incr}
		whileLoopBody = append(whileLoopBody, incrStmt)
	}
	whileLoop := spec.WhileStmt{
		Condition: cond,
		Body: spec.BlockStmt{
			Statements: whileLoopBody,
		},
	};
	// { init; while loop }:
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
