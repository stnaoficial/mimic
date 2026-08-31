package lang

import (
	"fmt"
	"mimic/internal/util"
	"reflect"
	"strings"
	"unicode"
)

type RegistryManager struct {
	entries map[string]any
}

func newRegistryManager() *RegistryManager {
	return &RegistryManager{
		entries: make(map[string]any),
	}
}

func (r *RegistryManager) Register(abstract string, concrete any) {
	r.entries[abstract] = concrete
}

func (r *RegistryManager) Assign(env *Environment) {
	for abstract, concrete := range r.entries {
		concreteType := reflect.TypeOf(concrete)
		concreteValue := reflect.ValueOf(concrete)

		switch concreteValue.Kind() {
		case reflect.Func:
			r.assignFunction(env, abstract, concrete, concreteType, concreteValue)
		}
	}
}

func (r *RegistryManager) assignFunction(
	env *Environment,
	abstract string,
	concrete any,
	concreteType reflect.Type,
	concreteValue reflect.Value,
) {
	env.Funcs[abstract] = func(args []any) (any, error) {
		if concreteValue.Kind() != reflect.Func {
			return nil, fmt.Errorf("%s is registered but is not a function", abstract)
		}

		if len(args) != concreteType.NumIn() {
			return nil, fmt.Errorf("Function %s expects %d arguments, but got %d", abstract, concreteType.NumIn(), len(args))
		}

		reflectArgs := make([]reflect.Value, len(args))

		for i := 0; i < len(args); i++ {
			expectedType := concreteType.In(i)

			givenVal := reflect.ValueOf(args[i])

			if !givenVal.IsValid() {
				reflectArgs[i] = reflect.Zero(expectedType)
				continue
			}

			if givenVal.Type() != expectedType {
				if givenVal.Type().ConvertibleTo(expectedType) {
					givenVal = givenVal.Convert(expectedType)
				} else {
					return nil, fmt.Errorf("Function %s expected argument %d type to be %s, but got %s", abstract, i, expectedType, givenVal.Type())
				}
			}

			reflectArgs[i] = givenVal
		}

		results := concreteValue.Call(reflectArgs)

		if len(results) < 2 {
			return nil, fmt.Errorf("Function %s error return is missing", abstract)
		}

		val := results[0].Interface()
		err := results[1].Interface()

		if err == nil {
			return val, nil
		}

		if err, ok := err.(error); ok {
			return val, err
		}

		return val, fmt.Errorf("Second return type must be an error")
	}
}

var Registry = newRegistryManager()

func init() {
	Registry.Register("upper", func(s string) (string, error) {
		return strings.ToUpper(s), nil
	})

	Registry.Register("lower", func(s string) (string, error) {
		return strings.ToLower(s), nil
	})

	Registry.Register("proper", func(s string) (string, error) {
		tokens := strings.Fields(strings.ToLower(s))

		if len(tokens) == 0 {
			return "", nil
		}

		for i := range tokens {
			runes := []rune(tokens[i])

			runes[0] = unicode.ToUpper(runes[0])

			tokens[i] = string(runes)
		}

		return strings.Join(tokens, " "), nil
	})

	Registry.Register("title", func(s string) (string, error) {
		tokens := strings.Fields(strings.ToLower(s))

		if len(tokens) == 0 {
			return "", nil
		}

		for i := range tokens {
			runes := []rune(tokens[i])

			if len(runes) > 2 || (i == 0 || i == len(tokens)-1) {
				runes[0] = unicode.ToUpper(runes[0])
			}

			tokens[i] = string(runes)
		}

		return strings.Join(tokens, " "), nil
	})

	Registry.Register("capitalize", func(s string) (string, error) {
		if s == "" {
			return "", nil
		}

		runes := []rune(strings.ToLower(s))

		runes[0] = unicode.ToUpper(runes[0])

		return string(runes), nil
	})

	Registry.Register("pascal", func(s string) (string, error) {
		tokens := strings.Fields(strings.ToLower(util.Normalize(s)))

		if len(tokens) == 0 {
			return "", nil
		}

		var builder strings.Builder

		for i := range tokens {
			runes := []rune(tokens[i])

			runes[0] = unicode.ToUpper(runes[0])

			builder.WriteString(string(runes))
		}

		return builder.String(), nil
	})

	Registry.Register("camel", func(s string) (string, error) {
		tokens := strings.Fields(strings.ToLower(util.Normalize(s)))

		if len(tokens) == 0 {
			return "", nil
		}

		var builder strings.Builder

		builder.WriteString(tokens[0])

		for i := 1; i < len(tokens); i++ {
			runes := []rune(tokens[i])

			runes[0] = unicode.ToUpper(runes[0])

			builder.WriteString(string(runes))
		}

		return builder.String(), nil
	})

	Registry.Register("flat", func(s string) (string, error) {
		return strings.Join(strings.Fields(strings.ToLower(util.Normalize(s))), ""), nil
	})

	Registry.Register("snake", func(s string) (string, error) {
		return util.Delimit(strings.ToLower(util.Normalize(s)), "_"), nil
	})

	Registry.Register("kebab", func(s string) (string, error) {
		return util.Delimit(strings.ToLower(util.Normalize(s)), "-"), nil
	})

	Registry.Register("before", func(value string, target string) (string, error) {
		before, _, ok := strings.Cut(value, target)

		if !ok {
			return value, nil
		}

		return before, nil
	})

	Registry.Register("after", func(value string, target string) (string, error) {
		_, after, ok := strings.Cut(value, target)

		if !ok {
			return value, nil
		}

		return after, nil
	})

	Registry.Register("between", func(value string, first string, last string) (string, error) {
		start := strings.Index(value, first)

		if start == -1 {
			return value, nil
		}

		start += len(first)

		end := strings.Index(value[start:], last)

		if end == -1 {
			return value, nil
		}

		return value[start : start+end], nil
	})

	Registry.Register("replace", func(str, old, new string) (string, error) {
		return strings.ReplaceAll(str, old, new), nil
	})

	Registry.Register("normalize", func(str string) (string, error) {
		return util.Normalize(str), nil
	})

	Registry.Register("delimit", func(str, del string) (string, error) {
		return util.Delimit(str, del), nil
	})

	Registry.Register("pad_left", func(value any, length int, pad string) (string, error) {
		str, err := util.StringValue(value)
		return util.Pad(str, length, pad, util.PadLeft), err
	})

	Registry.Register("pad_right", func(value any, length int, pad string) (string, error) {
		str, err := util.StringValue(value)
		return util.Pad(str, length, pad, util.PadRight), err
	})

	Registry.Register("pad_both", func(value any, length int, pad string) (string, error) {
		str, err := util.StringValue(value)
		return util.Pad(str, length, pad, util.PadBoth), err
	})

	Registry.Register("zero_fill", func(value any, length int) (string, error) {
		str, err := util.StringValue(value)
		return util.ZeroFill(str, length), err
	})
}
