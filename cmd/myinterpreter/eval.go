package main

import (
	"errors"
	"fmt"
	"reflect"
)

func evaluate(expr Expr) (any, error) {
	return expr.Eval()
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
