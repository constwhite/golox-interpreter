package environment

import (
	"errors"

	t "github.com/constwhite/golox-interpreter/token"
)

type Environment struct {
	Values map[string]interface{}
}

var ErrorUndefinedVar = errors.New("variable is not defined")

func NewEnvironment() *Environment {
	return &Environment{Values: make(map[string]interface{})}
}

func (env *Environment) Define(name string, value interface{}) {
	env.Values[name] = value
}

func (env *Environment) Get(name t.Token) (interface{}, error) {
	if value, ok := env.Values[name.Lexeme]; ok {
		return value, nil
	}

	return nil, ErrorUndefinedVar
}

func (env *Environment) Assign(name t.Token, value interface{}) error {
	if _, ok := env.Values[name.Lexeme]; ok {
		env.Values[name.Lexeme] = value
		return nil
	}
	return ErrorUndefinedVar

}
