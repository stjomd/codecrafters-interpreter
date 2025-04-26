package interpreter

type Class struct {
	Name string
}
func (class Class) String() string {
	return class.Name
}
