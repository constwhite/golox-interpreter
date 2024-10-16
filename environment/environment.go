package environment

import (
	"errors"

	t "github.com/constwhite/golox-interpreter/token"
)

type Environment struct {
	Values    map[string]interface{}
	Enclosing *Environment
}

var ErrorUndefinedVar = errors.New("variable is not defined")

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{Values: make(map[string]interface{}), Enclosing: enclosing}
}

func (env *Environment) Define(name string, value interface{}) {
	env.Values[name] = value
}

func (env *Environment) Get(name t.Token) (interface{}, error) {
	if value, ok := env.Values[name.Lexeme]; ok {
		return value, nil
	}
	if env.Enclosing != nil {

		return env.Enclosing.Get(name)
	}
	return nil, ErrorUndefinedVar
}
func (env *Environment) GetAt(distance int, name string) interface{} {
	return env.ancestor(distance).Values[name]
}
func (env *Environment) ancestor(distance int) *Environment {
	environment := env
	for i := 0; i < distance; i++ {
		environment = environment.Enclosing
	}
	return environment
}

func (env *Environment) Assign(name t.Token, value interface{}) error {
	if _, ok := env.Values[name.Lexeme]; ok {
		env.Values[name.Lexeme] = value
		return nil
	}
	if env.Enclosing != nil {
		return env.Enclosing.Assign(name, value)
	}

	return ErrorUndefinedVar

}

func (env *Environment) AssignAt(distance int, name t.Token, value interface{}) {
	env.ancestor(distance).Values[name.Lexeme] = value
}
