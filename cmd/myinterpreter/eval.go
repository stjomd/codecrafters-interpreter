package main

import (
	"errors"
	"fmt"
	"reflect"
)

func evaluate(expr Expr) (any, error) {
	return expr.Eval()
}

// MARK: - Eval() implementations

func (le LiteralExpr) Eval() (any, error) {
	return le.value, nil
}

func (ge GroupingExpr) Eval() (any, error) {
	return ge.expr.Eval()
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
	
	message := fmt.Sprintf("Unexpected type of unary expression: %s.", ue.operation.Type.String())
	return nil, runtimeError(message, ue.operation.Line)
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
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeError("Operands must be numbers.", be.operation.Line)
		}
		return leftValue.(float64) - rightValue.(float64), nil
	case Plus:
		if isNumber(leftValue) && isNumber(rightValue) {
			return leftValue.(float64) + rightValue.(float64), nil
		}
		if isString(leftValue) && isString(rightValue) {
			return leftValue.(string) + rightValue.(string), nil
		}
		return nil, runtimeError("Operands must be two numbers or two strings.", be.operation.Line)
	case Less:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeError("Operands must be numbers.", be.operation.Line)
		}
		return leftValue.(float64) < rightValue.(float64), nil
	case LessEqual:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeError("Operands must be numbers.", be.operation.Line)
		}
		return leftValue.(float64) <= rightValue.(float64), nil
	case Greater:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeError("Operands must be numbers.", be.operation.Line)
		}
		return leftValue.(float64) > rightValue.(float64), nil
	case GreaterEqual:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeError("Operands must be numbers.", be.operation.Line)
		}
		return leftValue.(float64) >= rightValue.(float64), nil
	case EqualEqual:
		return isEqual(leftValue, rightValue), nil
	case BangEqual:
		return !isEqual(leftValue, rightValue), nil
	}

	message := fmt.Sprintf("Unexpected type of binary expression: %s.", be.operation.Type.String())
	return nil, runtimeError(message, be.operation.Line)
}


// MARK: - Helpers

func runtimeError(message string, line uint64) error {
	errorMessage := fmt.Sprintf("%s\n[line %d]", message, line)
	return errors.New(errorMessage)
}

func isTruthy(value any) bool {
	if value == false || value == nil { return false }
	return true
}

func isEqual(a any, b any) bool {
	if a == nil && b == nil { return true }
	if a == nil { return false }
	return a == b
}

func isNumber(value any) bool {
	return reflect.TypeOf(value).Kind() == reflect.Float64
}

func isString(value any) bool {
	return reflect.TypeOf(value).Kind() == reflect.String
}
