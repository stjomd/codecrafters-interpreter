package interpreter

import (
	"fmt"
)

type Class struct { // implements Callable
	Name string
	Methods map[string]Function
}
func (class Class) String() string {
	return class.Name
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
	} else if method, contains := inst.findMethod(name); contains {
		return method, nil
	} else {
		return nil, fmt.Errorf("undefined property %v", name)
	}
}
func (inst ClassInstance) set(name string, value any) error {
	inst.Fields[name] = value
	return nil
}

func (inst ClassInstance) findMethod(name string) (Function, bool) {
	function, contains := inst.Class.Methods[name]
	return function, contains // Go...
}

// MARK: - Class Callable
func (class Class) arity() int {
	return 0
}
func (class Class) call(intp *Interpreter, args []any) (any, error) {
	inst := ClassInstance{Class: &class, Fields: make(map[string]any)}
	return inst, nil
}
