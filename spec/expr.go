package spec

import (
	"fmt"
	"reflect"
)

// MARK: - Expressions

type Expr interface {
	String() string
	Eval(evaluator ExprVisitor[any, error]) (any, error)
}

type ExprVisitor[R any, E error] interface {
	VisitLiteral(le LiteralExpr) (R, E)
	VisitGrouping(ge GroupingExpr) (R, E)
	VisitUnary(ue UnaryExpr) (R, E)
	VisitBinary(be BinaryExpr) (R, E)
	VisitVariable(ve VariableExpr) (R, E)
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
func (le LiteralExpr) Eval(evaluator ExprVisitor[any, error]) (any, error) {
	return evaluator.VisitLiteral(le)
}

type GroupingExpr struct {
	Expr Expr
}
func (ge GroupingExpr) String() string {
	return fmt.Sprintf("(group %v)", ge.Expr)
}
func (ge GroupingExpr) Eval(evaluator ExprVisitor[any, error]) (any, error) {
	return evaluator.VisitGrouping(ge)
}

type UnaryExpr struct {
	Opt Token
	Expr Expr
}
func (ue UnaryExpr) String() string {
	return fmt.Sprintf("(%v %v)", ue.Opt.Lexeme, ue.Expr)
}
func (ue UnaryExpr) Eval(evaluator ExprVisitor[any, error]) (any, error) {
	return evaluator.VisitUnary(ue)
}

type BinaryExpr struct {
	Left Expr
	Opt Token
	Right Expr
}
func (be BinaryExpr) String() string {
	return fmt.Sprintf("(%v %v %v)", be.Opt.Lexeme, be.Left, be.Right)
}
func (be BinaryExpr) Eval(evaluator ExprVisitor[any, error]) (any, error) {
	return evaluator.VisitBinary(be)
}

type VariableExpr struct {
	Identifier string
}
func (ve VariableExpr) String() string {
	return ve.Identifier
}
func (ve VariableExpr) Eval(evaluator ExprVisitor[any, error]) (any, error) {
	return evaluator.VisitVariable(ve)
}
