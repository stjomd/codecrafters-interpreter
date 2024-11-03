package api

import "errors"

type environment struct {
	parent *environment
	variables map[string]any
}

func newEnv() environment {
	return environment{variables: make(map[string]any)}
}

func newEnvWithParent(parent *environment) environment {
	return environment{parent: parent, variables: make(map[string]any)}
}

func (env *environment) define(name string, value any) {
	env.variables[name] = value
}

func (env *environment) assign(name string, value any) error {
	_, isPresent := env.variables[name]
	if isPresent {
		env.variables[name] = value
	}
	if env.parent != nil {
		return env.parent.assign(name, value)
	}
	return nil
}

func (env *environment) get(name string) (any, error) {
	value, isPresent := env.variables[name]
	if isPresent {
		return value, nil
	}
	if env.parent != nil {
		return env.parent.get(name)
	}
	return nil, errors.New("Undefined variable '" + name + "'")
}
