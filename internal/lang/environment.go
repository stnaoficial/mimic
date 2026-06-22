package lang

import (
	"errors"
	"mimic/internal/lang/env"
)

type Function func(args []any) (any, error)

type Environment struct {
	Vars    map[string]string
	Prompts map[string]string
	Funcs   map[string]Function
}

func buildDefaultVars() map[string]string {
	return map[string]string{}
}

func buildDefaultPrompts() map[string]string {
	return map[string]string{}
}

func buildDefaultFuncs() map[string]Function {
	return map[string]Function{
		"upper": func(args []any) (any, error) {
			proto := "upper(value: string): string"

			if len(args) != 1 {
				return nil, errors.New(proto)
			}

			arg0, ok := args[0].(string)

			if !ok {
				return nil, errors.New(proto)
			}

			return env.Upper(arg0), nil
		},
		"lower": func(args []any) (any, error) {
			proto := "lower(value: string): string"

			if len(args) != 1 {
				return nil, errors.New(proto)
			}

			arg0, ok := args[0].(string)

			if !ok {
				return nil, errors.New(proto)
			}

			return env.Lower(arg0), nil
		},

		"proper": func(args []any) (any, error) {
			proto := "proper(value: string): string"

			if len(args) != 1 {
				return nil, errors.New(proto)
			}

			arg0, ok := args[0].(string)

			if !ok {
				return nil, errors.New(proto)
			}

			return env.Proper(arg0), nil
		},
		"title": func(args []any) (any, error) {
			proto := "title(value: string): string"

			if len(args) != 1 {
				return nil, errors.New(proto)
			}

			arg0, ok := args[0].(string)

			if !ok {
				return nil, errors.New(proto)
			}

			return env.Title(arg0), nil
		},
		"capitalize": func(args []any) (any, error) {
			proto := "capitalize(value: string): string"

			if len(args) != 1 {
				return nil, errors.New(proto)
			}

			arg0, ok := args[0].(string)

			if !ok {
				return nil, errors.New(proto)
			}

			return env.Capitalize(arg0), nil
		},

		"pascal": func(args []any) (any, error) {
			proto := "pascal(value: string): string"

			if len(args) != 1 {
				return nil, errors.New(proto)
			}

			arg0, ok := args[0].(string)

			if !ok {
				return nil, errors.New(proto)
			}

			return env.Pascal(arg0), nil
		},
		"camel": func(args []any) (any, error) {
			proto := "camel(value: string): string"

			if len(args) != 1 {
				return nil, errors.New(proto)
			}

			arg0, ok := args[0].(string)

			if !ok {
				return nil, errors.New(proto)
			}

			return env.Camel(arg0), nil
		},
		"flat": func(args []any) (any, error) {
			proto := "flat(value: string): string"

			if len(args) != 1 {
				return nil, errors.New(proto)
			}

			arg0, ok := args[0].(string)

			if !ok {
				return nil, errors.New(proto)
			}

			return env.Flat(arg0), nil
		},

		"kebab": func(args []any) (any, error) {
			proto := "kebab(value: string): string"

			if len(args) != 1 {
				return nil, errors.New(proto)
			}

			arg0, ok := args[0].(string)

			if !ok {
				return nil, errors.New(proto)
			}

			return env.Kebab(arg0), nil
		},
		"snake": func(args []any) (any, error) {
			proto := "snake(value: string): string"

			if len(args) != 1 {
				return nil, errors.New(proto)
			}

			arg0, ok := args[0].(string)

			if !ok {
				return nil, errors.New(proto)
			}

			return env.Snake(arg0), nil
		},

		"before": func(args []any) (any, error) {
			proto := "before(value: string, target: string): string"

			if len(args) != 2 {
				return nil, errors.New(proto)
			}

			arg0, ok := args[0].(string)
			arg1, ok := args[1].(string)

			if !ok {
				return nil, errors.New(proto)
			}

			return env.Before(arg0, arg1), nil
		},
		"after": func(args []any) (any, error) {
			proto := "after(value: string, target: string): string"

			if len(args) != 2 {
				return nil, errors.New(proto)
			}

			arg0, ok := args[0].(string)
			arg1, ok := args[1].(string)

			if !ok {
				return nil, errors.New(proto)
			}

			return env.After(arg0, arg1), nil
		},
		"between": func(args []any) (any, error) {
			proto := "between(value: string, first: string, last: string): string"

			if len(args) != 3 {
				return nil, errors.New(proto)
			}

			arg0, ok := args[0].(string)
			arg1, ok := args[1].(string)
			arg2, ok := args[2].(string)

			if !ok {
				return nil, errors.New(proto)
			}

			return env.Between(arg0, arg1, arg2), nil
		},
		"replace": func(args []any) (any, error) {
			proto := "replace(value: string, old: string, new: string): string"

			if len(args) != 3 {
				return nil, errors.New(proto)
			}

			arg0, ok := args[0].(string)
			arg1, ok := args[1].(string)
			arg2, ok := args[2].(string)

			if !ok {
				return nil, errors.New(proto)
			}

			return env.Replace(arg0, arg1, arg2), nil
		},
		"normalize": func(args []any) (any, error) {
			proto := "normalize(value: string): string"

			if len(args) != 1 {
				return nil, errors.New(proto)
			}

			arg0, ok := args[0].(string)

			if !ok {
				return nil, errors.New(proto)
			}

			return env.Normalize(arg0), nil
		},
		"delimit": func(args []any) (any, error) {
			proto := "delimit(value: string, delimiter: string): string"

			if len(args) != 2 {
				return nil, errors.New(proto)
			}

			arg0, ok := args[0].(string)
			arg1, ok := args[1].(string)

			if !ok {
				return nil, errors.New(proto)
			}

			return env.Delimit(arg0, arg1), nil
		},
		"pad_left": func(args []any) (any, error) {
			proto := "pad_left(value: any, length: number, pad: string): string"

			if len(args) != 3 {
				return nil, errors.New(proto)
			}

			arg1, ok := args[1].(float64)
			arg2, ok := args[2].(string)

			if !ok {
				return nil, errors.New(proto)
			}

			return env.PadLeft(args[0], int(arg1), arg2)
		},
		"pad_right": func(args []any) (any, error) {
			proto := "pad_right(value: any, length: number, pad: string): string"

			if len(args) != 3 {
				return nil, errors.New(proto)
			}

			arg1, ok := args[1].(float64)
			arg2, ok := args[2].(string)

			if !ok {
				return nil, errors.New(proto)
			}

			return env.PadRight(args[0], int(arg1), arg2)
		},
		"pad": func(args []any) (any, error) {
			proto := "pad(value: any, length: number, pad: string): string"

			if len(args) != 3 {
				return nil, errors.New(proto)
			}

			arg1, ok := args[1].(float64)
			arg2, ok := args[2].(string)

			if !ok {
				return nil, errors.New(proto)
			}

			return env.PadBoth(args[0], int(arg1), arg2)
		},
		"zero_fill": func(args []any) (any, error) {
			proto := "zero_fill(value: any, length: number): string"

			if len(args) != 2 {
				return nil, errors.New(proto)
			}

			arg1, ok := args[1].(float64)

			if !ok {
				return nil, errors.New(proto)
			}

			return env.ZeroFill(args[0], int(arg1))
		},
	}
}

func NewEnvironment() *Environment {
	return &Environment{
		Vars:    buildDefaultVars(),
		Prompts: buildDefaultPrompts(),
		Funcs:   buildDefaultFuncs(),
	}
}
