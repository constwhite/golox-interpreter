package interpreter

import (
	abs "github.com/constwhite/golox-interpreter/abstractSyntaxTree"
	env "github.com/constwhite/golox-interpreter/environment"
)

type loxCallable interface {
	call(*Interpreter, []interface{}) interface{}
	arity() int
}

type loxFunction struct {
	Declaration   abs.FunctionStmt
	Closure       *env.Environment
	isInitialiser bool
}

func (f loxFunction) call(interpreter *Interpreter, args []interface{}) (returnVal interface{}) {
	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(returnValue); ok {
				if f.isInitialiser {
					returnVal = f.Closure.GetAt(0, "this")
					return
				}
				returnVal = v.Value
				return

			}
			panic(err)
		}
	}()
	env := env.NewEnvironment(f.Closure)
	for i := 0; i < len(f.Declaration.Params); i++ {
		env.Define(f.Declaration.Params[i].Lexeme, args[i])
	}
	interpreter.executeBlock(f.Declaration.Body, env)
	if f.isInitialiser {
		return f.Closure.GetAt(0, "this")
	}
	return nil
}

func (f loxFunction) arity() int {
	return len(f.Declaration.Params)
}
func (f loxFunction) bind(instance *loxInstance) loxFunction {
	environment := env.Environment{Enclosing: f.Closure}
	environment.Define("this", instance)
	return loxFunction{Declaration: f.Declaration, Closure: &environment, isInitialiser: f.isInitialiser}
}
