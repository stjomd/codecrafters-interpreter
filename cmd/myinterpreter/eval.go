package main

func evaluate(expr Expr) any {
	return expr.Eval()
}

func isTruthy(value any) bool {
	if value == false || value == nil {
		return false
	}
	return true
}
