package lang

import (
	"fmt"
	"mimic/core/cli"
	"strconv"
)

type Evaluator struct {
	env *Environment
}

func NewEvaluator(env *Environment) *Evaluator {
	return &Evaluator{
		env: env,
	}
}

func (e *Evaluator) Eval(node Node) (string, error) {
	switch n := node.(type) {

	case Identifier:
		if value, ok := e.env.Vars[n.Name]; ok {
			return value, nil
		}

		prompt, ok := e.env.Prompts[n.Name]

		if !ok {
			prompt = fmt.Sprintf("Please enter a value for \"%s\": ", n.Name)
		}

		value := cli.MustAsk(prompt)

		e.env.Vars[n.Name] = value

		return value, nil

	case StringLiteral:
		return strconv.Unquote(n.Value)

	case CallExpression:
		fn, ok := e.env.Funcs[n.Name]

		if !ok {
			return "", fmt.Errorf("Undefined function \"%s\"", n.Name)
		}

		var args []string

		for _, arg := range n.Args {
			if result, err := e.Eval(arg); err == nil {
				args = append(args, result)
			} else {
				return "", err
			}
		}

		return fn(args), nil
	}

	return "", nil
}
