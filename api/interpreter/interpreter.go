package api

import "github.com/codecrafters-io/interpreter-starter-go/spec"

type interpreter struct { // implements spec.ExprVisitor[any, error], spec.StmtVisitor[error]
	env *environment
}

func NewInterpreter() interpreter {
	globals := newGlobalsEnv()
	env := newEnvWithParent(&globals)
	return interpreter{env: &env}
}

func (intp *interpreter) Interpret(statements *[]spec.Stmt) error {
	for _, stmt := range *statements {
		if err := stmt.Exec(intp); err != nil {
			return err
		}
	}
	return nil
}

func (intp *interpreter) Evaluate(expr *spec.Expr) (any, error) {
	return (*expr).Eval(intp)
}
