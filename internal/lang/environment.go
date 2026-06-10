package lang

import (
	"fmt"
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
			if len(args) != 1 {
				return nil, fmt.Errorf("upper(value: string): string")
			}

			arg0, ok := args[0].(string)

			if !ok {
				return nil, fmt.Errorf("upper(value: string): string")
			}

			return env.Upper(arg0), nil
		},
		"lower": func(args []any) (any, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("lower(value: string): string")
			}

			arg0, ok := args[0].(string)

			if !ok {
				return nil, fmt.Errorf("lower(value: string): string")
			}

			return env.Lower(arg0), nil
		},

		"proper": func(args []any) (any, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("proper(value: string): string")
			}

			arg0, ok := args[0].(string)

			if !ok {
				return nil, fmt.Errorf("proper(value: string): string")
			}

			return env.Proper(arg0), nil
		},
		"title": func(args []any) (any, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("title(value: string): string")
			}

			arg0, ok := args[0].(string)

			if !ok {
				return nil, fmt.Errorf("title(value: string): string")
			}

			return env.Title(arg0), nil
		},
		"capitalize": func(args []any) (any, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("capitalize(value: string): string")
			}

			arg0, ok := args[0].(string)

			if !ok {
				return nil, fmt.Errorf("capitalize(value: string): string")
			}

			return env.Capitalize(arg0), nil
		},

		"pascal": func(args []any) (any, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("pascal(value: string): string")
			}

			arg0, ok := args[0].(string)

			if !ok {
				return nil, fmt.Errorf("pascal(value: string): string")
			}

			return env.Pascal(arg0), nil
		},
		"camel": func(args []any) (any, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("camel(value: string): string")
			}

			arg0, ok := args[0].(string)

			if !ok {
				return nil, fmt.Errorf("camel(value: string): string")
			}

			return env.Camel(arg0), nil
		},
		"flat": func(args []any) (any, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("flat(value: string): string")
			}

			arg0, ok := args[0].(string)

			if !ok {
				return nil, fmt.Errorf("flat(value: string): string")
			}

			return env.Flat(arg0), nil
		},

		"kebab": func(args []any) (any, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("kebab(value: string): string")
			}

			arg0, ok := args[0].(string)

			if !ok {
				return nil, fmt.Errorf("kebab(value: string): string")
			}

			return env.Kebab(arg0), nil
		},
		"snake": func(args []any) (any, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("snake(value: string): string")
			}

			arg0, ok := args[0].(string)

			if !ok {
				return nil, fmt.Errorf("snake(value: string): string")
			}

			return env.Snake(arg0), nil
		},

		"replace": func(args []any) (any, error) {
			if len(args) != 3 {
				return nil, fmt.Errorf("replace(value: string, old: string, new: string): string")
			}

			arg0, ok := args[0].(string)
			arg1, ok := args[0].(string)
			arg2, ok := args[0].(string)

			if !ok {
				return nil, fmt.Errorf("replace(value: string, old: string, new: string): string")
			}

			return env.Replace(arg0, arg1, arg2), nil
		},
		"normalize": func(args []any) (any, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("normalize(value: string): string")
			}

			arg0, ok := args[0].(string)

			if !ok {
				return nil, fmt.Errorf("normalize(value: string): string")
			}

			return env.Normalize(arg0), nil
		},
		"join": func(args []any) (any, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("join(value: string, separator: string): string")
			}

			arg0, ok := args[0].(string)
			arg1, ok := args[0].(string)

			if !ok {
				return nil, fmt.Errorf("join(value: string, separator: string): string")
			}

			return env.Join(arg0, arg1), nil
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
