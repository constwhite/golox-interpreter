package interpreter

import (
	"fmt"

	t "github.com/constwhite/golox-interpreter/token"
)

type loxClass struct {
	Name    string
	methods map[string]loxFunction
}

func (c loxClass) call(interpreter *Interpreter, args []interface{}) interface{} {
	instance := loxInstance{Class: c}
	initialiser := c.findMethod("init")
	if initialiser != nil {
		initialiser.bind(&instance).call(interpreter, args)
	}
	return instance
}

func (c loxClass) arity() int {
	initialiser := c.findMethod("init")
	if initialiser == nil {
		return 0
	}
	return initialiser.arity()
}
func (c loxClass) findMethod(name string) *loxFunction {
	if method, ok := c.methods[name]; ok {
		return &method
	}
	return nil
}

// type Instance interface {
// 	Get(Interpreter, t.Token) (interface{}, error)
// }

type loxInstance struct {
	Class  loxClass
	Fields map[string]interface{}
}

func (in *loxInstance) get(name t.Token) (interface{}, error) {
	if field, ok := in.Fields[name.Lexeme]; ok {
		return field, nil
	}

	method := in.Class.findMethod(name.Lexeme)
	if method != nil {
		return method.bind(in), nil
	}
	err := fmt.Errorf("undefined property '%v'", name.Lexeme)
	return nil, err

}

func (in *loxInstance) set(name t.Token, value interface{}) {
	in.Fields[name.Lexeme] = value
}
