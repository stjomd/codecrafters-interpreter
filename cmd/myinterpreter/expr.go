package main

import (
	"fmt"
	"reflect"
)

// MARK: - Expressions

type Expr interface {
	String() string
	Eval() (any, error)
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
func (le LiteralExpr) Eval() (any, error) {
	return le.value, nil
}

type GroupingExpr struct {
	expr Expr
}
func (ge GroupingExpr) String() string {
	return fmt.Sprintf("(group %v)", ge.expr)
}
func (ge GroupingExpr) Eval() (any, error) {
	return ge.expr.Eval()
}

type UnaryExpr struct {
	operation Token
	expr Expr
}
func (ue UnaryExpr) String() string {
	return fmt.Sprintf("(%v %v)", ue.operation.Lexeme, ue.expr)
}
func (ue UnaryExpr) Eval() (any, error) {
	subvalue, suberror := ue.expr.Eval()
	if suberror != nil { return nil, suberror }

	switch ue.operation.Type {
	case Bang:
		return !isTruthy(subvalue), nil
	case Minus:
		if !isNumber(subvalue) {
			return nil, runtimeError("Operand must be a number.", ue.operation.Line)
		}
		return -subvalue.(float64), nil
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
func (be BinaryExpr) Eval() (any, error) {
	leftValue, leftError := be.left.Eval() 
	rightValue, rightError := be.right.Eval()
	if leftError != nil { return nil, leftError }
	if rightError != nil { return nil, rightError }

	switch be.operation.Type {
	case Star:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeError("Operands must be numbers.", be.operation.Line)
		}
		return leftValue.(float64) * rightValue.(float64), nil
	case Slash:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeError("Operands must be numbers.", be.operation.Line)
		}
		return leftValue.(float64) / rightValue.(float64), nil
	case Minus:
		return leftValue.(float64) - rightValue.(float64), nil
	case Plus:
		if isString(leftValue) && isString(rightValue) {
			return leftValue.(string) + rightValue.(string), nil
		}
		return leftValue.(float64) + rightValue.(float64), nil
	case Less:
		return leftValue.(float64) < rightValue.(float64), nil
	case LessEqual:
		return leftValue.(float64) <= rightValue.(float64), nil
	case Greater:
		return leftValue.(float64) > rightValue.(float64), nil
	case GreaterEqual:
		return leftValue.(float64) >= rightValue.(float64), nil
	case EqualEqual:
		return isEqual(leftValue, rightValue), nil
	case BangEqual:
		return !isEqual(leftValue, rightValue), nil
	}

	panic("! binary eval")
}
