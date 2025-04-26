package interpreter

import (
	"fmt"
	"reflect"

	"github.com/codecrafters-io/interpreter-starter-go/spec"
)

func (intp *Interpreter) VisitPrint(ps spec.PrintStmt) error {
	value, evalError := ps.Expr.Eval(intp)
	if evalError != nil { return evalError }
	if value == nil {
		fmt.Println("nil")
	} else if reflect.TypeOf(value).Kind() == reflect.Float64 {
		fmt.Println(float64ToString(value.(float64)))
	} else {
		fmt.Println(value)
	}
	return nil
}
func float64ToString(number float64) string {
	if number == float64(int(number)) {
		return fmt.Sprintf("%.0f", number)
	}
	return fmt.Sprintf("%g", number)
}


func (intp *Interpreter) VisitExpr(es spec.ExprStmt) error {
	if es.Expr == nil { return nil }
	if _, evalError := es.Expr.Eval(intp); evalError != nil {
		return evalError
	}
	return nil
}

func (intp *Interpreter) VisitDeclare(ds spec.DeclareStmt) error {
	value, evalError := ds.Expr.Eval(intp)
	if evalError != nil { return evalError }
	intp.env.define(ds.Identifier.Lexeme, value)
	return nil
}

func (intp *Interpreter) VisitBlock(bs spec.BlockStmt) error {
	env := newEnvWithParent(intp.env)
	return intp.ExecBlock(&bs.Statements, &env)
}

func (intp *Interpreter) ExecBlock(stmts *[]spec.Stmt, env *environment) error {
	origEnv := intp.env
	intp.env = env
	defer func() { intp.env = origEnv }()

	for _, stmt := range *stmts {
		if err := stmt.Exec(intp); err != nil {
			return err;
		}
	}
	return nil
}

func (intp *Interpreter) VisitIf(is spec.IfStmt) error {
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

func (intp *Interpreter) VisitWhile(ws spec.WhileStmt) error {
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

func (intp *Interpreter) VisitFunc(fs spec.FuncStmt) error {
	intp.env.define(
		fs.Name.Lexeme,
		Function {
			declaration: fs,
			closure: intp.env,
		},
	)
	return nil
}

func (intp *Interpreter) VisitReturn(rs spec.ReturnStmt) error {
	if rs.Expr == nil {
		return Return{value: nil}
	}
	value, evalError := rs.Expr.Eval(intp)
	if evalError != nil {
		return evalError
	}
	return Return{value: value}
}

func (intp *Interpreter) VisitClass(cs spec.ClassStmt) error {
	intp.env.define(cs.Name.Lexeme, nil)
	class := Class{Name: cs.Name.Lexeme}
	intp.env.assign(cs.Name.Lexeme, class)
	return nil
}
