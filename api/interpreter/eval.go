package interpreter

import (
	"fmt"
	"reflect"

	"github.com/codecrafters-io/interpreter-starter-go/spec"
)

func (intp interpreter) VisitLiteral(le spec.LiteralExpr) (any, error) {
	return le.Value, nil
}

func (intp interpreter) VisitGrouping(ge spec.GroupingExpr) (any, error) {
	return ge.Expr.Eval(intp)
}

func (intp interpreter) VisitUnary(ue spec.UnaryExpr) (any, error) {
	subvalue, suberror := ue.Expr.Eval(intp)
	if suberror != nil { return nil, suberror }

	switch ue.Opt.Type {
	case spec.Bang:
		return !isTruthy(subvalue), nil
	case spec.Minus:
		if !isNumber(subvalue) {
			return nil, runtimeError{message: "Operand must be a number", line: ue.Opt.Line}
		}
		return -subvalue.(float64), nil
	}
	
	message := fmt.Sprintf("Unexpected type of unary expression: %s", ue.Opt.Type.String())
	return nil, runtimeError{message: message, line: ue.Opt.Line}
}

func (intp interpreter) VisitBinary(be spec.BinaryExpr) (any, error) {
	leftValue, leftError := be.Left.Eval(intp)
	rightValue, rightError := be.Right.Eval(intp)
	if leftError != nil { return nil, leftError }
	if rightError != nil { return nil, rightError }

	switch be.Opt.Type {
	case spec.Star:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeError{message: operandsMustBeNumbers, line: be.Opt.Line}
		}
		return leftValue.(float64) * rightValue.(float64), nil
	case spec.Slash:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeError{message: operandsMustBeNumbers, line: be.Opt.Line}
		}
		return leftValue.(float64) / rightValue.(float64), nil
	case spec.Minus:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeError{message: operandsMustBeNumbers, line: be.Opt.Line}
		}
		return leftValue.(float64) - rightValue.(float64), nil
	case spec.Plus:
		if isNumber(leftValue) && isNumber(rightValue) {
			return leftValue.(float64) + rightValue.(float64), nil
		}
		if isString(leftValue) && isString(rightValue) {
			return leftValue.(string) + rightValue.(string), nil
		}
		return nil, runtimeError{message: "Operands must be two numbers or two strings", line: be.Opt.Line}
	case spec.Less:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeError{message: operandsMustBeNumbers, line: be.Opt.Line}
		}
		return leftValue.(float64) < rightValue.(float64), nil
	case spec.LessEqual:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeError{message: operandsMustBeNumbers, line: be.Opt.Line}
		}
		return leftValue.(float64) <= rightValue.(float64), nil
	case spec.Greater:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeError{message: operandsMustBeNumbers, line: be.Opt.Line}
		}
		return leftValue.(float64) > rightValue.(float64), nil
	case spec.GreaterEqual:
		if !isNumber(leftValue) || !isNumber(rightValue) {
			return nil, runtimeError{message: operandsMustBeNumbers, line: be.Opt.Line}
		}
		return leftValue.(float64) >= rightValue.(float64), nil
	case spec.EqualEqual:
		return isEqual(leftValue, rightValue), nil
	case spec.BangEqual:
		return !isEqual(leftValue, rightValue), nil
	}

	message := fmt.Sprintf("Unexpected type of binary expression: %s", be.Opt.Type.String())
	return nil, runtimeError{message: message, line: be.Opt.Line}
}

func (intp interpreter) VisitVariable(be spec.VariableExpr) (any, error) {
	if value, err := intp.env.get(be.Identifier.Lexeme); err == nil {
		return value, nil
	} else {
		return nil, runtimeError{message: err.Error(), line: be.Identifier.Line, cause: err}
	}
}

func (intp interpreter) VisitAssignment(ae spec.AssignmentExpr) (any, error) {
	value, evalError := ae.Expr.Eval(intp)
	if evalError != nil {
		return nil, runtimeError{message: evalError.(runtimeError).message, line: ae.Identifier.Line, cause: evalError}
	}
	assignError := intp.env.assign(ae.Identifier.Lexeme, value)
	if assignError != nil { 
		return nil, runtimeError{message: assignError.(runtimeError).message, line: ae.Identifier.Line, cause: assignError}
	}
	return value, nil
}

func (intp interpreter) VisitLogical(le spec.LogicalExpr) (any, error) {
	left, leftError := le.Left.Eval(intp)
	if leftError != nil { return nil, leftError }
	// short circuit
	if le.Opt.Type == spec.And && !isTruthy(left) {
		return left, nil
	} else if le.Opt.Type == spec.Or && isTruthy(left) {
		return left, nil
	}
	return le.Right.Eval(intp)
}

func (intp interpreter) VisitCall(ce spec.CallExpr) (any, error) {
	callee, calleeError := ce.Callee.Eval(intp)
	if calleeError != nil { return nil, calleeError }
	args := []any{}
	for _, arg := range ce.Args {
		evaledArg, evalError := arg.Eval(intp)
		if evalError != nil {
			return nil, evalError
		}
		args = append(args, evaledArg)
	}
	function, castOk := callee.(Callable)
	if !castOk {
		return nil, runtimeError{message: "can only call functions and classes", line: ce.Paren.Line}
	}
	if len(args) != int(function.arity()) {
		msg := fmt.Sprintf("expected %v arguments but got %v", function.arity(), len(args))
		return nil, runtimeError{message: msg, line: ce.Paren.Line}
	}
	return function.call(&intp, args)
}

// MARK: - Helpers

const operandsMustBeNumbers = "Operands must be numbers"

type runtimeError struct {
	message string
	line uint64
	cause error
}
func (re runtimeError) Error() string {
	return fmt.Sprintf("%s.\n[line %d]", re.message, re.line)
	// if re.cause != nil {
	// 	return fmt.Sprintf("%s.\n[line %d]\n- caused by: %v", re.message, re.line, re.cause)
	// } else {
	// 	return fmt.Sprintf("%s.\n[line %d]", re.message, re.line)
	// }
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
