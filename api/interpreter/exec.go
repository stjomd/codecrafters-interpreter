package interpreter

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/spec"
)

func (intp *interpreter) VisitPrint(ps spec.PrintStmt) error {
	value, evalError := ps.Expr.Eval(intp)
	if evalError != nil { return evalError }
	if value == nil {
		fmt.Println("nil")
	} else {
		fmt.Println(value)
	}
	return nil
}

func (intp *interpreter) VisitExpr(es spec.ExprStmt) error {
	if es.Expr == nil { return nil }
	if _, evalError := es.Expr.Eval(intp); evalError != nil {
		return evalError
	}
	return nil
}

func (intp *interpreter) VisitDeclare(ds spec.DeclareStmt) error {
	value, evalError := ds.Expr.Eval(intp)
	if evalError != nil { return evalError }
	intp.env.define(ds.Identifier.Lexeme, value)
	return nil
}

func (intp *interpreter) VisitBlock(bs spec.BlockStmt) error {
	outerEnv := intp.env
	innerEnv := newEnvWithParent(intp.env)
	intp.env = &innerEnv
	for _, stmt := range bs.Statements {
		if err := stmt.Exec(intp); err != nil {
			return err;
		}
	}
	intp.env = outerEnv
	return nil
}

func (intp *interpreter) VisitIf(is spec.IfStmt) error {
	condition, conditionError := is.Condition.Eval(intp)
	if conditionError != nil { return conditionError }
	if isTruthy(condition) {
		if err := is.Then.Exec(intp); err != nil {
			return err
		}
	} else if is.Else != nil {
		if err := is.Else.Exec(intp); err != nil {
			return err
		}
	}
	return nil
}

func (intp *interpreter) VisitWhile(ws spec.WhileStmt) error {
	fulfiled, err := ws.Condition.Eval(intp)
	if err != nil { return err }
	for isTruthy(fulfiled) {
		if err := ws.Body.Exec(intp); err != nil {
			return err
		}
		fulfiled, err = ws.Condition.Eval(intp)
		if err != nil { return err }
	}
	return nil
}

func (intp *interpreter) VisitFunc(fs spec.FuncStmt) error {
	intp.env.define(
		fs.Name.Lexeme,
		Function {
			declaration: fs,
		},
	)
	return nil
}

func (intp *interpreter) VisitReturn(rs spec.ReturnStmt) error {
	if rs.Expr == nil {
		return Return{value: nil}
	}
	value, evalError := rs.Expr.Eval(intp)
	if evalError != nil {
		return evalError
	}
	return Return{value: value}
}
