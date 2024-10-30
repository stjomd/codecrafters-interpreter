package spec

import (
	"errors"
	"fmt"
	"reflect"
)

// MARK: - Eval() implementations

func (le LiteralExpr) Eval() (any, error) {
	return le.Value, nil
}

func (ge GroupingExpr) Eval() (any, error) {
	return ge.Expr.Eval()
}

func (ue UnaryExpr) Eval() (any, error) {
	subvalue, suberror := ue.Expr.Eval()
	if suberror != nil { return nil, suberror }

	switch ue.Opt.Type {
	case Bang:
		return !isTruthy(subvalue), nil
	case Minus:
		if !isNumber(subvalue) {
			return nil, runtimeError("Operand must be a number.", ue.Opt.Line)
		}
		return -subvalue.(float64), nil
	}
	
	message := fmt.Sprintf("Unexpected type of unary expression: %s.", ue.Opt.Type.String())
	return nil, runtimeError(message, ue.Opt.Line)
}

func (be BinaryExpr) Eval() (any, error) {
	leftValue, leftError := be.Left.Eval() 
	rightValue, rightError := be.Right.Eval()
	if leftError != nil { return nil, leftError }
	if rightError != nil { return nil, rightError }

	switch be.Opt.Type {
	case Star:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeErrorMustBeNumbers(be.Opt.Line)
		}
		return leftValue.(float64) * rightValue.(float64), nil
	case Slash:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeErrorMustBeNumbers(be.Opt.Line)
		}
		return leftValue.(float64) / rightValue.(float64), nil
	case Minus:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeErrorMustBeNumbers(be.Opt.Line)
		}
		return leftValue.(float64) - rightValue.(float64), nil
	case Plus:
		if isNumber(leftValue) && isNumber(rightValue) {
			return leftValue.(float64) + rightValue.(float64), nil
		}
		if isString(leftValue) && isString(rightValue) {
			return leftValue.(string) + rightValue.(string), nil
		}
		return nil, runtimeError("Operands must be two numbers or two strings.", be.Opt.Line)
	case Less:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeErrorMustBeNumbers(be.Opt.Line)
		}
		return leftValue.(float64) < rightValue.(float64), nil
	case LessEqual:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeErrorMustBeNumbers(be.Opt.Line)
		}
		return leftValue.(float64) <= rightValue.(float64), nil
	case Greater:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeErrorMustBeNumbers(be.Opt.Line)
		}
		return leftValue.(float64) > rightValue.(float64), nil
	case GreaterEqual:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeErrorMustBeNumbers(be.Opt.Line)
		}
		return leftValue.(float64) >= rightValue.(float64), nil
	case EqualEqual:
		return isEqual(leftValue, rightValue), nil
	case BangEqual:
		return !isEqual(leftValue, rightValue), nil
	}

	message := fmt.Sprintf("Unexpected type of binary expression: %s.", be.Opt.Type.String())
	return nil, runtimeError(message, be.Opt.Line)
}


// MARK: - Helpers

func runtimeErrorMustBeNumbers(line uint64) error {
	return runtimeError("Operands must be numbers.", line)
}

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
