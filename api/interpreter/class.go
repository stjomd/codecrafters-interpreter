package interpreter

type Class struct { // implements Callable
	Name string
}
func (class Class) String() string {
	return class.Name
}

type ClassInstance struct {
	Class *Class
}
func (instance ClassInstance) String() string {
	return instance.Class.Name + " instance"
}

// MARK: - Class Callable
func (class Class) arity() int {
	return 0
}
func (class Class) call(intp *Interpreter, args []any) (any, error) {
	instance := ClassInstance{Class: &class}
	return instance, nil
}
