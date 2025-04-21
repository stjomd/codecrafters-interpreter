package api

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/spec"
)

func Exec(statements *[]spec.Stmt) error {
	env := newEnv()
	evaluator := evalVisitor{env: &env}
	executor := executor{evaluator: &evaluator}
	for _, stmt := range *statements {
		err := stmt.Exec(&executor)
		if err != nil { return err }
	}
	return nil
}

// MAKE: - Execution using visitor pattern

type executor struct { // implements spec.ExecVisitor
	// env *environment
	evaluator *evalVisitor
}

func (exec *executor) VisitPrint(ps spec.PrintStmt) error {
	value, evalError := Eval(&ps.Expr, exec.evaluator.env)
	if evalError != nil { return evalError }
	if value == nil {
		fmt.Println("nil")
	} else {
		fmt.Println(value)
	}
	return nil
}

func (exec *executor) VisitExpr(es spec.ExprStmt) error {
	if es.Expr == nil { return nil }
	if _, evalError := Eval(&es.Expr, exec.evaluator.env); evalError != nil {
		return evalError
	}
	return nil
}

func (exec *executor) VisitDeclare(ds spec.DeclareStmt) error {
	value, evalError := Eval(&ds.Expr, exec.evaluator.env)
	if evalError != nil { return evalError }
	exec.evaluator.env.define(ds.Identifier.Lexeme, value)
	return nil
}

func (exec *executor) VisitBlock(bs spec.BlockStmt) error {
	outerEnv := exec.evaluator.env
	innerEnv := newEnvWithParent(exec.evaluator.env)
	exec.evaluator.env = &innerEnv
	for _, stmt := range bs.Statements {
		err := stmt.Exec(exec)
		if err != nil { return err }
	}
	exec.evaluator.env = outerEnv
	return nil
}

func (exec *executor) VisitIf(is spec.IfStmt) error {
	condition, conditionError := Eval(&is.Condition, exec.evaluator.env)
	if conditionError != nil { return conditionError }
	if isTruthy(condition) {
		is.Then.Exec(exec)
	} else if is.Else != nil {
		is.Else.Exec(exec)
	}
	return nil
}

func (exec *executor) VisitWhile(ws spec.WhileStmt) error {
	fulfiled, err := Eval(&ws.Condition, exec.evaluator.env)
	if err != nil { return err }
	for isTruthy(fulfiled) {
		ws.Body.Exec(exec)
		fulfiled, err = Eval(&ws.Condition, exec.evaluator.env)
		if err != nil { return err }
	}
	return nil
}

func (exec *executor) VisitFunc(fs spec.FuncStmt) error {
	function := Function{declaration: fs}
	exec.evaluator.env.define(
		fs.Name.Lexeme,
		function,
	)
	return nil
}
