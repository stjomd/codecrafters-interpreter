package main

func evaluate(expr Expr) any {
	return expr.Eval()
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
