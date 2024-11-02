package api

import "errors"

type Environment struct {
	variables map[string]any
}

func NewEnv() Environment {
	return Environment{variables: make(map[string]any)}
}

func (env *Environment) Define(name string, value any) {
	env.variables[name] = value
}

func (env *Environment) Assign(name string, value any) error {
	if _, isPresent := env.variables[name]; !isPresent {
		return errors.New("Undefined variable '" + name + "'")
	}
	env.variables[name] = value
	return nil
}

func (env *Environment) Get(name string) (any, error) {
	value, isPresent := env.variables[name]
	if !isPresent { return nil, errors.New("Undefined variable '" + name + "'") }
	return value, nil
}
