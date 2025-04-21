package api

import "github.com/codecrafters-io/interpreter-starter-go/spec"

type Callable interface {
	arity() int
	call(execVisitor spec.StmtVisitor[error], args []any) any
}

type Function struct { // implements Callable
	declaration spec.FuncStmt
}
func (f Function) arity() int {
	return len(f.declaration.Params)
}
func (f Function) call(executor spec.StmtVisitor[error], args []any) any {
	env := newEnv()
	for i, param := range f.declaration.Params {
		env.define(param.Lexeme, args[i])
	}
	return f.declaration.Exec(executor)
}

type NativeFunction struct {
	_arity int
	_func func(args []any) any
}
func (nf NativeFunction) arity() int {
	return nf._arity
}
func (nf NativeFunction) call(executor spec.StmtVisitor[error], args []any) any {
	return nf._func(args)
}
