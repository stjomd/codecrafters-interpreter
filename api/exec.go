package api

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/spec"
)

func Exec(statements *[]spec.Stmt) error {
	env := NewEnv()
	executor := execVisitor{env: &env}
	for _, stmt := range *statements {
		err := stmt.Exec(executor)
		if err != nil { return err }
	}
	return nil
}

// MAKE: - Execution using visitor pattern

type execVisitor struct { // implements spec.ExecVisitor
	env *Environment
}

func (ev execVisitor) VisitPrint(ps spec.PrintStmt) error {
	value, evalError := Eval(&ps.Expr, ev.env)
	if evalError != nil { return evalError }
	if value == nil {
		fmt.Println("nil")
	} else {
		fmt.Println(value)
	}
	return nil
}

func (ev execVisitor) VisitExpr(es spec.ExprStmt) error {
	if _, evalError := Eval(&es.Expr, ev.env); evalError != nil {
		return evalError
	}
	return nil
}

func (ev execVisitor) VisitDeclare(ds spec.DeclareStmt) error {
	value, evalError := Eval(&ds.Expr, ev.env)
	if evalError != nil { return evalError }
	ev.env.Define(ds.Identifier.Lexeme, value)
	return nil
}
