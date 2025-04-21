package api

import (
	"time"
)

func newGlobalsEnv() environment {
	env := environment{variables: make(map[string]any)}
	for _, fn := range nativeFunctions {
		env.variables[fn._name] = fn
	}
	return env
}

var nativeFunctions = []NativeFunction {
	{
		_name: "clock",
		_arity: 0,
		_func: func(args []any) any {
			return float64(time.Now().Unix())
		},
	},
	{
		_name: "echo",
		_arity: 1,
		_func: func(args []any) any {
			return args[0]
		},
	},
}
