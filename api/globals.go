package api

import "time"

func newGlobalsEnv() environment {
	return environment{variables: globals}
}

var globals = map[string]any{
	"clock": Callable{
		_arity: 0,
		_func: func(ev evalVisitor, args []any) any {
			return float64(time.Now().Unix())
		},
	},
}
