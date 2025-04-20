package api

type Callable struct {
	_arity uint8
	_func func(evalVisitor, []any) any
}
func (f *Callable) arity() uint8 {
	return f._arity;
}
func (f *Callable) call(evalVisitor evalVisitor, args []any) any {
	return f._func(evalVisitor, args)
}
