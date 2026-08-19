package lang

import (
	"fmt"
	"reflect"
)

type Function func(args []any) (any, error)

type Environment struct {
	Vars    map[string]string
	Prompts map[string]string
	Funcs   map[string]Function
}

func registerDefaultVars() map[string]string {
	return map[string]string{}
}

func registerDefaultPrompts() map[string]string {
	return map[string]string{}
}

func registerDefaultFuncs() map[string]Function {
	funcs := make(map[string]Function)

	for protoName, fn := range Registry {
		fnVal := reflect.ValueOf(fn)
		fnType := reflect.TypeOf(fn)

		funcs[protoName] = func(args []any) (any, error) {
			if fnVal.Kind() != reflect.Func {
				return nil, fmt.Errorf("%s is registered but is not a function", protoName)
			}

			if len(args) != fnType.NumIn() {
				return nil, fmt.Errorf("Function %s expects %d arguments, but got %d", protoName, fnType.NumIn(), len(args))
			}

			reflectArgs := make([]reflect.Value, len(args))

			for i := 0; i < len(args); i++ {
				expectedType := fnType.In(i)

				givenVal := reflect.ValueOf(args[i])

				if !givenVal.IsValid() {
					reflectArgs[i] = reflect.Zero(expectedType)
					continue
				}

				if givenVal.Type() != expectedType {
					if givenVal.Type().ConvertibleTo(expectedType) {
						givenVal = givenVal.Convert(expectedType)
					} else {
						return nil, fmt.Errorf("Function %s expected argument %d type to be %s, but got %s", protoName, i, expectedType, givenVal.Type())
					}
				}

				reflectArgs[i] = givenVal
			}

			results := fnVal.Call(reflectArgs)

			switch len(results) {
			case 0:
				return nil, nil
			case 1:
				return results[0].Interface(), nil
			case 2:
				val := results[0].Interface()
				errVal := results[1].Interface()

				if errVal == nil {
					return val, nil
				}

				if err, ok := errVal.(error); ok {
					return val, err
				}

				return val, fmt.Errorf("Unknown second return error type")
			default:
				return nil, fmt.Errorf("Unsupported multiple return counts for function %s", protoName)
			}
		}
	}

	return funcs
}

func NewEnvironment() *Environment {
	return &Environment{
		Vars:    registerDefaultVars(),
		Prompts: registerDefaultPrompts(),
		Funcs:   registerDefaultFuncs(),
	}
}
