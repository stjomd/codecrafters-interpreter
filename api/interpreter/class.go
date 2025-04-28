package interpreter

import (
	"fmt"
)

type Class struct { // implements Callable
	Name string
	Methods map[string]Function
	Superclass *Class
}
func (class Class) String() string {
	return class.Name
}
func (class Class) findMethod(name string) (Function, bool) {
	function, contains := class.Methods[name]
	return function, contains // Go...
}

type ClassInstance struct {
	Class *Class
	Fields map[string]any
}
func (inst ClassInstance) String() string {
	return inst.Class.Name + " instance"
}
func (inst ClassInstance) get(name string) (any, error) {
	if value, contains := inst.Fields[name]; contains {
		return value, nil
	} else if method, contains := inst.Class.findMethod(name); contains {
		return method.bind(inst), nil
	} else {
		return nil, fmt.Errorf("undefined property %v", name)
	}
}
func (inst ClassInstance) set(name string, value any) error {
	inst.Fields[name] = value
	return nil
}

type ClassType int
const (
	CtNone = iota
	CtClass
)

// MARK: - Class Callable

func (class Class) arity() int {
	if init, contains := class.findMethod("init"); contains {
		return init.arity()
	}
	return 0
}
func (class Class) call(intp *Interpreter, args []any) (any, error) {
	inst := ClassInstance{Class: &class, Fields: make(map[string]any)}
	if init, contains := inst.Class.findMethod("init"); contains {
		init.bind(inst).call(intp, args)
	}
	return inst, nil
}
