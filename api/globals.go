package api

import (
	"time"
)

func newGlobalsEnv() environment {
	return environment{variables: globals}
}

var globals = map[string]any{
	"clock": NativeFunction {
		_arity: 0,
		_func: func(args []any) any {
			return float64(time.Now().Unix())
		},
	},
	"echo": NativeFunction {
		_arity: 1,
		_func: func(args []any) any {
			return args[0]
		},
	},
}
