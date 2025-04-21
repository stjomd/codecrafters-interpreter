package interpreter

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/spec"
)

type Callable interface {
	arity() int
	call(interpreter *interpreter, args []any) any
}

// MARK: - Lox Functions

type Function struct { // implements Callable
	declaration spec.FuncStmt
}
func (f Function) arity() int {
	return len(f.declaration.Params)
}
func (f Function) call(interpreter *interpreter, args []any) any {
	subenv := newEnvWithParent(interpreter.env)
	interpreter.env = &subenv
	defer func(){ interpreter.env = interpreter.env.parent }();

	for i, param := range f.declaration.Params {
		interpreter.env.define(param.Lexeme, args[i])
	}
	result := f.declaration.Body.Exec(interpreter);
	if returnValue, ok := result.(Return); ok {
		return returnValue.value
	}
	return nil
}
func (f Function) String() string {
	return fmt.Sprintf("<fn %v>", f.declaration.Name.Lexeme)
}

// MARK: - Native Functions

type NativeFunction struct {
	_name string
	_arity int
	_func func(args []any) any
}
func (nf NativeFunction) arity() int {
	return nf._arity
}
func (nf NativeFunction) call(interpreter *interpreter, args []any) any {
	return nf._func(args)
}
func (nf NativeFunction) String() string {
	return fmt.Sprintf("<nat fn %v>", nf._name)
}

// MARK: - Return "Error"

type Return struct {
	value any
}
func (r Return) Error() string {
	return fmt.Sprintf(
		"error: this should not be an error! Some function is returning the value %v, but the runtime did not catch this.",
		r.value,
	)
}
