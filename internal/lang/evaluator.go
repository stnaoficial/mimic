package lang

import (
	"fmt"
	"mimic/internal/cli"
)

type Evaluator struct {
	env *Environment

	strict bool
}

func NewEvaluator(env *Environment, strict bool) *Evaluator {
	return &Evaluator{
		env: env,

		strict: strict,
	}
}

func (e *Evaluator) solve(name string) string {
	prompt, ok := e.env.Prompts[name]

	if !ok {
		prompt = fmt.Sprintf("Please enter a value for \"%s\": ", name)
	}

	return cli.MustAsk(prompt)
}

func (e *Evaluator) Eval(node Node) (any, error) {
	switch n := node.(type) {

	case Null:
		return "", nil

	case String:
		return n.Value, nil

	case Identifier:
		if value, ok := e.env.Vars[n.Name]; ok {
			return value, nil
		}

		if e.strict {
			return nil, fmt.Errorf("Undefined variable \"%s\"", n.Name)
		}

		value := e.solve(n.Name)

		e.env.Vars[n.Name] = value

		return value, nil

	case Callable:
		fn, ok := e.env.Funcs[n.Name]

		if !ok {
			return nil, fmt.Errorf("Undefined function \"%s\"", n.Name)
		}

		var args []any

		for _, arg := range n.Args {
			if value, err := e.Eval(arg); err == nil {
				args = append(args, value)
			} else {
				return nil, err
			}
		}

		value, err := fn(args)

		if err != nil {
			return nil, fmt.Errorf("Usage: %s", err)
		}

		return value, nil
	}

	return "", nil
}

func (e *Evaluator) EvalAsString(node Node) (string, error) {
	value, err := e.Eval(node)

	if err != nil {
		return "", err
	}

	if str, ok := value.(string); ok {
		return str, nil
	}

	return "", fmt.Errorf("Unsupported string conversion")
}
