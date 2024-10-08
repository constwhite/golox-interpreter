package interpreter

type loxCallable interface {
	call(*Interpreter, []interface{}) interface{}
	arity() int
}
