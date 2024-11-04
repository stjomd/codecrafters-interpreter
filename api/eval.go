package api

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/codecrafters-io/interpreter-starter-go/spec"
)

func EvalWithoutEnv(expr *spec.Expr) (any, error) {
	env := newEnv()
	return (*expr).Eval(evalVisitor{env: &env})
}

func Eval(expr *spec.Expr, env *environment) (any, error) {
	return (*expr).Eval(evalVisitor{env: env})
}

// MARK: - Evaluation using visitor pattern

type evalVisitor struct { // implements spec.Visitor
	env *environment
}

func (ev evalVisitor) VisitLiteral(le spec.LiteralExpr) (any, error) {
	return le.Value, nil
}

func (ev evalVisitor) VisitGrouping(ge spec.GroupingExpr) (any, error) {
	return ge.Expr.Eval(ev)
}

func (ev evalVisitor) VisitUnary(ue spec.UnaryExpr) (any, error) {
	subvalue, suberror := ue.Expr.Eval(ev)
	if suberror != nil { return nil, suberror }

	switch ue.Opt.Type {
	case spec.Bang:
		return !isTruthy(subvalue), nil
	case spec.Minus:
		if !isNumber(subvalue) {
			return nil, runtimeError("Operand must be a number", ue.Opt.Line)
		}
		return -subvalue.(float64), nil
	}
	
	message := fmt.Sprintf("Unexpected type of unary expression: %s", ue.Opt.Type.String())
	return nil, runtimeError(message, ue.Opt.Line)
}

func (ev evalVisitor) VisitBinary(be spec.BinaryExpr) (any, error) {
	leftValue, leftError := be.Left.Eval(ev)
	rightValue, rightError := be.Right.Eval(ev)
	if leftError != nil { return nil, leftError }
	if rightError != nil { return nil, rightError }

	switch be.Opt.Type {
	case spec.Star:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeError(operandsMustBeNumbers, be.Opt.Line)
		}
		return leftValue.(float64) * rightValue.(float64), nil
	case spec.Slash:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeError(operandsMustBeNumbers, be.Opt.Line)
		}
		return leftValue.(float64) / rightValue.(float64), nil
	case spec.Minus:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeError(operandsMustBeNumbers, be.Opt.Line)
		}
		return leftValue.(float64) - rightValue.(float64), nil
	case spec.Plus:
		if isNumber(leftValue) && isNumber(rightValue) {
			return leftValue.(float64) + rightValue.(float64), nil
		}
		if isString(leftValue) && isString(rightValue) {
			return leftValue.(string) + rightValue.(string), nil
		}
		return nil, runtimeError("Operands must be two numbers or two strings", be.Opt.Line)
	case spec.Less:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeError(operandsMustBeNumbers, be.Opt.Line)
		}
		return leftValue.(float64) < rightValue.(float64), nil
	case spec.LessEqual:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeError(operandsMustBeNumbers, be.Opt.Line)
		}
		return leftValue.(float64) <= rightValue.(float64), nil
	case spec.Greater:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeError(operandsMustBeNumbers, be.Opt.Line)
		}
		return leftValue.(float64) > rightValue.(float64), nil
	case spec.GreaterEqual:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeError(operandsMustBeNumbers, be.Opt.Line)
		}
		return leftValue.(float64) >= rightValue.(float64), nil
	case spec.EqualEqual:
		return isEqual(leftValue, rightValue), nil
	case spec.BangEqual:
		return !isEqual(leftValue, rightValue), nil
	}

	message := fmt.Sprintf("Unexpected type of binary expression: %s", be.Opt.Type.String())
	return nil, runtimeError(message, be.Opt.Line)
}

func (ev evalVisitor) VisitVariable(be spec.VariableExpr) (any, error) {
	value, err := ev.env.get(be.Identifier.Lexeme)
	if err != nil { return nil, runtimeError(err.Error(), be.Identifier.Line) }
	return value, nil
}

func (ev evalVisitor) VisitAssignment(ae spec.AssignmentExpr) (any, error) {
	value, evalError := ae.Expr.Eval(ev)
	if evalError != nil { return nil, runtimeError(evalError.Error(), ae.Identifier.Line) }
	assignError := ev.env.assign(ae.Identifier.Lexeme, value)
	if assignError != nil { return nil, runtimeError(assignError.Error(), ae.Identifier.Line) }
	return value, nil
}

func (ev evalVisitor) VisitLogical(le spec.LogicalExpr) (any, error) {
	left, leftError := le.Left.Eval(ev)
	if leftError != nil { return nil, leftError }
	// short circuit
	if le.Opt.Type == spec.And && !isTruthy(left) {
		return left, nil
	} else if le.Opt.Type == spec.Or && isTruthy(left) {
		return left, nil
	}
	return le.Right.Eval(ev)
}


// MARK: - Helpers

const operandsMustBeNumbers = "Operands must be numbers"
func runtimeError(message string, line uint64) error {
	errorMessage := fmt.Sprintf("%s.\n[line %d]", message, line)
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
