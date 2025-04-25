package interpreter

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/spec"
)

type Interpreter struct { // implements spec.ExprVisitor[any, error], spec.StmtVisitor[error]
	env *environment
	globals *environment
	locals map[uint64]int // map of spec.Expr.Hash() -> int
}

func NewInterpreter() Interpreter {
	env := newGlobalsEnv()
	return Interpreter{env: &env, globals: &env, locals: make(map[uint64]int)}
}

func (intp *Interpreter) Resolve(expr spec.Expr, depth int) {
	intp.locals[expr.Hash()] = depth
}

func (intp *Interpreter) lookUpVar(name spec.Token, expr spec.Expr) (any, error) {
	distance, contains := intp.locals[expr.Hash()]
	if contains {
		return intp.env.getAt(distance, name.Lexeme)
	} else {
		return intp.globals.get(name.Lexeme)
	}
}

func (intp *Interpreter) ReportError(token spec.Token, message string) {
	switch token.Type {
	case spec.EOF:
		fmt.Fprintf(os.Stderr, "[line %v] Error at end: %v\n", token.Line, message)
	default:
		fmt.Fprintf(os.Stderr, "[line %v] Error at '%v': %v\n", token.Line, token.Lexeme, message)
	}
}
