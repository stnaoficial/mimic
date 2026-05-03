package lang

import (
	"fmt"
	"maps"
	"math/rand/v2"
	"mimic/core/env"
	"path/filepath"
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

func (e *Environment) DefineScopeVars(fileName string, fileData string) {
	e.Vars["__16_DIGIT_INT__"] = fmt.Sprintf("%016d", rand.IntN(10_000_000_000_000_000))
	e.Vars["__8_DIGIT_INT__"] = fmt.Sprintf("%08d", rand.IntN(100_000_000))
	e.Vars["__4_DIGIT_INT__"] = fmt.Sprintf("%04d", rand.IntN(10_000))
	e.Vars["__FILENAME__"] = fileName
	e.Vars["__DIRNAME__"] = filepath.Dir(fileName)
	e.Vars["__FILE__"] = fileData
}
