package main

import (
	"fmt"
	"reflect"
)

// MARK: - Expressions

type Expr interface {
	String() string
	Eval() any
}

type LiteralExpr struct {
	value any
}
func (le LiteralExpr) String() string {
	if le.value == nil {
		return "nil"
	} else if reflect.TypeOf(le.value).Kind() == reflect.Float64 {
		return Float64ToString(le.value.(float64))
	} else {
		return fmt.Sprint(le.value)
	}
}
func (le LiteralExpr) Eval() any {
	return le.value
}

type GroupingExpr struct {
	expr Expr
}
func (ge GroupingExpr) String() string {
	return fmt.Sprintf("(group %v)", ge.expr)
}
func (ge GroupingExpr) Eval() any {
	return ge.expr.Eval()
}

type UnaryExpr struct {
	operation Token
	expr Expr
}
func (ue UnaryExpr) String() string {
	return fmt.Sprintf("(%v %v)", ue.operation.Lexeme, ue.expr)
}
func (ue UnaryExpr) Eval() any {
	subvalue := ue.expr.Eval()
	switch ue.operation.Type {
	case Bang:
		return !isTruthy(subvalue)
	case Minus:
		return -subvalue.(float64)
	}
	panic("! unary eval")
}

type BinaryExpr struct {
	left Expr
	operation Token
	right Expr
}
func (be BinaryExpr) String() string {
	return fmt.Sprintf("(%v %v %v)", be.operation.Lexeme, be.left, be.right)
}
func (be BinaryExpr) Eval() any {
	panic("!")
}
