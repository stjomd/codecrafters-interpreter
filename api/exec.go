package api

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/spec"
)

func Exec(statements *[]spec.Stmt) error {
	env := newEnv()
	executor := execVisitor{env: &env}
	for _, stmt := range *statements {
		err := stmt.Exec(&executor)
		if err != nil { return err }
	}
	return nil
}

// MAKE: - Execution using visitor pattern

type execVisitor struct { // implements spec.ExecVisitor
	env *environment
}

func (ev *execVisitor) VisitPrint(ps spec.PrintStmt) error {
	value, evalError := Eval(&ps.Expr, ev.env)
	if evalError != nil { return evalError }
	if value == nil {
		fmt.Println("nil")
	} else {
		fmt.Println(value)
	}
	return nil
}

func (ev *execVisitor) VisitExpr(es spec.ExprStmt) error {
	if _, evalError := Eval(&es.Expr, ev.env); evalError != nil {
		return evalError
	}
	return nil
}

func (ev *execVisitor) VisitDeclare(ds spec.DeclareStmt) error {
	value, evalError := Eval(&ds.Expr, ev.env)
	if evalError != nil { return evalError }
	ev.env.define(ds.Identifier.Lexeme, value)
	return nil
}

func (ev *execVisitor) VisitBlock(bs spec.BlockStmt) error {
	outerEnv := ev.env
	innerEnv := newEnvWithParent(ev.env)
	ev.env = &innerEnv
	for _, stmt := range bs.Statements {
		err := stmt.Exec(ev)
		if err != nil { return err }
	}
	ev.env = outerEnv
	return nil
}

func (ev *execVisitor) VisitIf(is spec.IfStmt) error {
	condition, conditionError := Eval(&is.Condition, ev.env)
	if conditionError != nil { return conditionError }
	if isTruthy(condition) {
		is.Then.Exec(ev)
	} else if is.Else != nil {
		is.Else.Exec(ev)
	}
	return nil
}

func (ev *execVisitor) VisitWhile(ws spec.WhileStmt) error {
	fulfiled, err := Eval(&ws.Condition, ev.env)
	if err != nil { return err }
	for isTruthy(fulfiled) {
		ws.Body.Exec(ev)
		fulfiled, err = Eval(&ws.Condition, ev.env)
		if err != nil { return err }
	}
	return nil
}
