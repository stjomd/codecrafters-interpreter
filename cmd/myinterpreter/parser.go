package main

import (
	"fmt"
	"reflect"
)

type Expr interface {
	String() string
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

func parse(tokens *[]Token) Expr {
	for i := 0; i < len(*tokens); i++ {
		token := (*tokens)[i]
		if token.Type == True {
			return LiteralExpr{value: true}
		} else if token.Type == False {
			return LiteralExpr{value: false}
		} else if token.Type == Nil {
			return LiteralExpr{value: nil}
		} else if token.Type == Number {
			return LiteralExpr{value: token.Literal}
		} else if token.Type == String {
			return LiteralExpr{value: token.Literal}
		}
	}
	panic("?")
}
