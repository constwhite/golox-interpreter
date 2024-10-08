package interpreter

import (
	"fmt"

	abs "github.com/constwhite/golox-interpreter/abstractSyntaxTree"
	env "github.com/constwhite/golox-interpreter/environment"
)

type loxCallable interface {
	call(*Interpreter, []interface{}) interface{}
	arity() int
	toString() string
}

type loxFunction struct {
	declaration abs.FunctionStmt
}

func (f loxFunction) call(interpreter *Interpreter, args []interface{}) (returnVal interface{}) {
	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(returnValue); ok {
				returnVal = v.Value
				return

			}
			panic(err)
		}
	}()
	env := env.NewEnvironment(interpreter.Globals)
	for i := 0; i < len(f.declaration.Params); i++ {
		env.Define(f.declaration.Params[i].Lexeme, args[i])
	}
	interpreter.executeBlock(f.declaration.Body, env)
	return nil
}

func (f loxFunction) arity() int {
	return len(f.declaration.Params)
}
func (f loxFunction) toString() string {
	return fmt.Sprintf("<fn %v >", f.declaration.Name.Lexeme)
}
