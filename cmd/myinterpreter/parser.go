package main

type Expr interface {
	String() string
}

type LiteralExpr struct {
	tType TokenType
	value any
}
func (le LiteralExpr) String() string {
	switch le.tType {
	case True:
		return "true"
	case False:
		return "false"
	case Nil:
		return "nil"
	}
	panic("unsupported token type in LiteralExpr.String()")
}

func parse(tokens *[]Token) []Expr {
	var expr []Expr
	for i := 0; i < len(*tokens); i++ {
		token := (*tokens)[i]
		if (token.Type == True) {
			expr = append(expr, LiteralExpr{tType: True, value: true})
		} else if (token.Type == False) {
			expr = append(expr, LiteralExpr{tType: False, value: false})
		} else if (token.Type == Nil) {
			expr = append(expr, LiteralExpr{tType: Nil, value: nil})
		}
	}
	return expr
}
