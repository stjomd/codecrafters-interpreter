package api

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/spec"
)

func Exec(statements *[]spec.Stmt) error {
	executor := execVisitor{}
	for _, stmt := range *statements {
		err := stmt.Exec(executor)
		if err != nil { return err }
	}
	return nil
}

type execVisitor struct {} // implements spec.ExecVisitor

func (ev execVisitor) VisitPrint(ps spec.PrintStmt) error {
	value, evalError := Eval(&ps.Expr)
	if evalError != nil { return evalError }
	fmt.Println(value)
	return nil
}

func (ev execVisitor) VisitExpr(es spec.ExprStmt) error {
	if _, evalError := Eval(&es.Expr); evalError != nil {
		return evalError
	}
	return nil
}
