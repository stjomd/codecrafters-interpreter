package spec

import (
	"fmt"
	"reflect"
)

// MARK: - Expressions

type Expr interface {
	String() string
	Eval() (any, error) // eval.go
}

type LiteralExpr struct {
	Value any
}
func (le LiteralExpr) String() string {
	if le.Value == nil {
		return "nil"
	} else if reflect.TypeOf(le.Value).Kind() == reflect.Float64 {
		return float64ToString(le.Value.(float64))
	} else {
		return fmt.Sprint(le.Value)
	}
}

type GroupingExpr struct {
	Expr Expr
}
func (ge GroupingExpr) String() string {
	return fmt.Sprintf("(group %v)", ge.Expr)
}

type UnaryExpr struct {
	Opt Token
	Expr Expr
}
func (ue UnaryExpr) String() string {
	return fmt.Sprintf("(%v %v)", ue.Opt.Lexeme, ue.Expr)
}

type BinaryExpr struct {
	Left Expr
	Opt Token
	Right Expr
}
func (be BinaryExpr) String() string {
	return fmt.Sprintf("(%v %v %v)", be.Opt.Lexeme, be.Left, be.Right)
}
