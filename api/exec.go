package api

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/spec"
)

func (exec *interpreter) VisitPrint(ps spec.PrintStmt) error {
	value, evalError := ps.Expr.Eval(exec)
	if evalError != nil { return evalError }
	if value == nil {
		fmt.Println("nil")
	} else {
		fmt.Println(value)
	}
	return nil
}

func (exec *interpreter) VisitExpr(es spec.ExprStmt) error {
	if es.Expr == nil { return nil }
	if _, evalError := es.Expr.Eval(exec); evalError != nil {
		return evalError
	}
	return nil
}

func (exec *interpreter) VisitDeclare(ds spec.DeclareStmt) error {
	value, evalError := ds.Expr.Eval(exec)
	if evalError != nil { return evalError }
	exec.env.define(ds.Identifier.Lexeme, value)
	return nil
}

func (exec *interpreter) VisitBlock(bs spec.BlockStmt) error {
	outerEnv := exec.env
	innerEnv := newEnvWithParent(exec.env)
	exec.env = &innerEnv
	for _, stmt := range bs.Statements {
		err := stmt.Exec(exec)
		if err != nil { return err }
	}
	exec.env = outerEnv
	return nil
}

func (exec *interpreter) VisitIf(is spec.IfStmt) error {
	condition, conditionError := is.Condition.Eval(exec)
	if conditionError != nil { return conditionError }
	if isTruthy(condition) {
		is.Then.Exec(exec)
	} else if is.Else != nil {
		is.Else.Exec(exec)
	}
	return nil
}

func (exec *interpreter) VisitWhile(ws spec.WhileStmt) error {
	fulfiled, err := ws.Condition.Eval(exec)
	if err != nil { return err }
	for isTruthy(fulfiled) {
		ws.Body.Exec(exec)
		fulfiled, err = ws.Condition.Eval(exec)
		if err != nil { return err }
	}
	return nil
}

func (exec *interpreter) VisitFunc(fs spec.FuncStmt) error {
	function := Function{declaration: fs}
	exec.env.define(
		fs.Name.Lexeme,
		function,
	)
	return nil
}
