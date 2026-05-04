package lang

import (
	"maps"
	"mimic/core/env"
)

type Environment struct {
	Vars  map[string]string
	Funcs map[string]func(args []string) string
}

func ReservedVariables() map[string]string {
	return map[string]string{}
}

func ReservedFunctions() map[string]func(args []string) string {
	funcs := map[string]func(args []string) string{}

	maps.Copy(funcs, env.Cases)
	maps.Copy(funcs, env.Modifiers)
	maps.Copy(funcs, env.Utils)

	return funcs
}

func NewEnvironment() *Environment {
	return &Environment{
		Vars:  ReservedVariables(),
		Funcs: ReservedFunctions(),
	}
}
