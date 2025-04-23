package interpreter

import (
	"errors"
)

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
	if _, isPresent := env.variables[name]; isPresent {
		env.variables[name] = value
		return nil
	}
	if env.parent != nil {
		return env.parent.assign(name, value)
	}
	return errors.New("Undefined variable '" + name + "'")
}

func (env *environment) assignAt(distance int, name string, value any) error {
	return env.ancestor(distance).assign(name, value)
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

func (env *environment) getAt(distance int, name string) (any, error) {
	return env.ancestor(distance).get(name)
}

func (env *environment) getGlobalsEnv() *environment {
	current := env
	for current.parent != nil {
		current = current.parent
	}
	return current
}

func (env *environment) ancestor(distance int) *environment {
	current := env
	for range distance {
		current = current.parent
	}
	return current
}
