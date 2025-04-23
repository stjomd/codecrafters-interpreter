package api

import (
	intp "github.com/codecrafters-io/interpreter-starter-go/api/interpreter"
	"github.com/codecrafters-io/interpreter-starter-go/spec"
)

func Eval(expr *spec.Expr) (any, error) {
	intp := intp.NewInterpreter()
	return (*expr).Eval(&intp)
}

func Exec(stmts *[]spec.Stmt) error {
	intp := intp.NewInterpreter()
	for _, stmt := range *stmts {
		if err := stmt.Exec(&intp); err != nil {
			return err
		}
	}
	return nil
}

func NewInterpreter() intp.Interpreter {
	return intp.NewInterpreter()
}

func ExecWithIntp(intp *intp.Interpreter, stmts *[]spec.Stmt) error {
	for _, stmt := range *stmts {
		if err := stmt.Exec(intp); err != nil {
			return err
		}
	}
	return nil
}

func ResolveWithIntp(intp *intp.Interpreter, stmts *[]spec.Stmt) {
	scopes := stack[map[string]bool]{slice: []map[string]bool{}}
	rslv := resolver{intp: intp, scopes: scopes}
	rslv.resolveStmts(stmts)
}
