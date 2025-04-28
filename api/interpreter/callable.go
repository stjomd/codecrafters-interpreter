package interpreter

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/spec"
)

type Callable interface {
	arity() int
	call(interpreter *Interpreter, args []any) (any, error)
}

type FunctionType int
const (
	FtNone FunctionType = iota
	FtMethod
	FtInitializer
	FtStandalone
)

// MARK: - Lox Functions

type Function struct { // implements Callable
	declaration spec.FuncStmt
	closure *environment
	isInit bool
}
func (f Function) arity() int {
	return len(f.declaration.Params)
}
func (f Function) call(interpreter *Interpreter, args []any) (any, error) {
	origEnv := interpreter.env
	subenv := newEnvWithParent(f.closure)
	interpreter.env = &subenv
	defer func(){ interpreter.env = origEnv }();

	for i, param := range f.declaration.Params {
		interpreter.env.define(param.Lexeme, args[i])
	}

	execError := interpreter.ExecBlock(&f.declaration.Body, &subenv)
	if returnValue, ok := execError.(Return); ok {
		if f.isInit {
			return f.closure.getAt(0, "this")
		}
		return returnValue.value, nil
	} else if execError != nil {
		return nil, execError
	}

	if f.isInit {
		return f.closure.getAt(0, "this")
	}
	return nil, nil
}
func (f Function) bind(inst ClassInstance) Function {
	closure := newEnvWithParent(f.closure)
	closure.define("this", inst)
	return Function{declaration: f.declaration, closure: &closure, isInit: f.isInit}
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
func (nf NativeFunction) call(interpreter *Interpreter, args []any) (any, error) {
	return nf._func(args), nil
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
