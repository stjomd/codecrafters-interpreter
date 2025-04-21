package api

import (
	"github.com/codecrafters-io/interpreter-starter-go/api/interpreter"
	"github.com/codecrafters-io/interpreter-starter-go/spec"
)

func Eval(expr *spec.Expr) (any, error) {
	intp := interpreter.NewInterpreter()
	return (*expr).Eval(intp)
}

func Exec(stmts *[]spec.Stmt) error {
	intp := interpreter.NewInterpreter()
	for _, stmt := range *stmts {
		if err := stmt.Exec(&intp); err != nil {
			return err
		}
	}
	return nil
}
