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
	VisitAssignment(assignmentExpr AssignmentExpr) (R, E)
	VisitBinary(binaryExpr BinaryExpr) (R, E)
	VisitCall(callExpr CallExpr) (R, E)
	VisitGrouping(groupingExpr GroupingExpr) (R, E)
	VisitLiteral(literalExpr LiteralExpr) (R, E)
	VisitLogical(logicalExpr LogicalExpr) (R, E)
	VisitUnary(unaryExpr UnaryExpr) (R, E)
	VisitVariable(variableExpr VariableExpr) (R, E)
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
	Identifier Token
}
func (ve VariableExpr) String() string {
	return fmt.Sprintf("(var %v)", ve.Identifier.Lexeme)
}
func (ve VariableExpr) Eval(evaluator ExprVisitor[any, error]) (any, error) {
	return evaluator.VisitVariable(ve)
}

type AssignmentExpr struct {
	Identifier Token
	Expr Expr
}
func (ae AssignmentExpr) String() string {
	return fmt.Sprintf("(assign %v %v)", ae.Identifier.Lexeme, ae.Expr)
}
func (ae AssignmentExpr) Eval(evaluator ExprVisitor[any, error]) (any, error) {
	return evaluator.VisitAssignment(ae)
}

type LogicalExpr struct {
	Left Expr
	Opt Token
	Right Expr
}
func (le LogicalExpr) String() string {
	return fmt.Sprintf("(%v %v %v)", le.Opt.Lexeme, le.Left, le.Right)
}
func (le LogicalExpr) Eval(evaluator ExprVisitor[any, error]) (any, error) {
	return evaluator.VisitLogical(le)
}

type CallExpr struct {
	Callee Expr
	Paren Token
	Args []Expr
}
func (ce CallExpr) String() string {
	return fmt.Sprintf("%v(%v)", ce.Callee, ce.Args)
}
func (ce CallExpr) Eval(evaluator ExprVisitor[any, error]) (any, error) {
	return evaluator.VisitCall(ce)
}
