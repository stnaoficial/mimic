package lang

import (
	"mimic/core/util"
)

type Variable = string
type Function = func(args []string) string

type Environment struct {
	Vars  map[string]Variable
	Funcs map[string]Function
}

func getReservedVariables() map[string]Variable {
	return map[string]Variable{}
}

func getReservedFunctions() map[string]Function {
	return map[string]Function{
		"sentence": func(args []string) string {
			if len(args) == 0 {
				return ""
			}
			return util.ToSentence(args[0])
		},
		"camel": func(args []string) string {
			if len(args) == 0 {
				return ""
			}
			return util.ToCamel(args[0])
		},
		"pascal": func(args []string) string {
			if len(args) == 0 {
				return ""
			}
			return util.ToPascal(args[0])
		},
		"snake": func(args []string) string {
			if len(args) == 0 {
				return ""
			}
			return util.ToSnake(args[0])
		},
		"kebab": func(args []string) string {
			if len(args) == 0 {
				return ""
			}
			return util.ToKebab(args[0])
		},
		"dot": func(args []string) string {
			if len(args) == 0 {
				return ""
			}
			return util.ToDot(args[0])
		},
		"flat": func(args []string) string {
			if len(args) == 0 {
				return ""
			}
			return util.ToFlat(args[0])
		},
		"lower": func(args []string) string {
			if len(args) == 0 {
				return ""
			}
			return util.ToLower(args[0])
		},
		"upper": func(args []string) string {
			if len(args) == 0 {
				return ""
			}
			return util.ToUpper(args[0])
		},
	}
}

func NewEnvironment() *Environment {
	return &Environment{
		Vars:  getReservedVariables(),
		Funcs: getReservedFunctions(),
	}
}
