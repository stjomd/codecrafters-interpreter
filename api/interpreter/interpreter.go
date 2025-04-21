package interpreter

type interpreter struct { // implements spec.ExprVisitor[any, error], spec.StmtVisitor[error]
	env *environment
}

func NewInterpreter() interpreter {
	globals := newGlobalsEnv()
	env := newEnvWithParent(&globals)
	return interpreter{env: &env}
}
