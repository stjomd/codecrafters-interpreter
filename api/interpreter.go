package api

import "github.com/codecrafters-io/interpreter-starter-go/spec"

type interpreter struct { // implements spec.ExprVisitor[any, error], spec.StmtVisitor[error]
	env *environment
}

func NewInterpreter() interpreter {
	env := newEnv()
	return interpreter{env: &env}
}

func (itp *interpreter) Interpret(statements *[]spec.Stmt) error {
	for _, stmt := range *statements {
		if err := stmt.Exec(itp); err != nil {
			return err
		}
	}
	return nil
}

func (itp *interpreter) Evaluate(expr *spec.Expr) (any, error) {
	return (*expr).Eval(itp)
}
