package interpreter

import (
	"craftinginterpreters/lox/scanner"
	"fmt"
)

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnv(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    map[string]any{},
	}
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) get(name scanner.Token) any {
	if v, ok := e.values[name.Lexeme]; ok {
		return v
	}
	if e.enclosing != nil {
		return e.enclosing.get(name)
	}
	panic(fmt.Sprintf("[line %d ]Undefined variable '%s'.", name.Line, name.Lexeme))
}

func (e *Environment) assign(name scanner.Token, value any) {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return
	}
	if e.enclosing != nil {
		e.enclosing.assign(name, value)
		return
	}
	panic(fmt.Sprintf("[line %d ]Undefined variable '%s'.", name.Line, name.Lexeme))
}
