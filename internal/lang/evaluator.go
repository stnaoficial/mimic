package lang

import (
	"fmt"
	"mimic/internal/cli"
	"mimic/internal/util"
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

func (e *Evaluator) prompt(name string) string {
	prompt, ok := e.env.Prompts[name]

	if !ok {
		prompt = fmt.Sprintf("Please enter a value for \"%s\": ", name)
	}

	return cli.MustPrompt(prompt)
}

func (e *Evaluator) Eval(node Node) (any, error) {
	switch n := node.(type) {

	case Null:
		return "", nil

	case String:
		return n.Value, nil

	case Number:
		return n.Value, nil

	case Identifier:
		if value, ok := e.env.Vars[n.Name]; ok {
			return value, nil
		}

		if e.strict {
			return nil, fmt.Errorf("Undefined variable \"%s\"", n.Name)
		}

		value := e.prompt(n.Name)

		e.env.Vars[n.Name] = value

		return value, nil

	case Unary:
		value, err := e.Eval(n.Value)

		if err != nil {
			return nil, err
		}

		number, err := util.NumberValue(value)

		if err != nil {
			return nil, err
		}

		switch n.Operator {
		case TokenPlus:
			return number, nil

		case TokenMinus:
			return -number, nil
		}

		return nil, fmt.Errorf("Unsupported unary operator")

	case Binary:
		left, err := e.Eval(n.Left)

		if err != nil {
			return nil, err
		}

		right, err := e.Eval(n.Right)

		if err != nil {
			return nil, err
		}

		leftNumber, err := util.NumberValue(left)

		if err != nil {
			return nil, err
		}

		rightNumber, err := util.NumberValue(right)

		if err != nil {
			return nil, err
		}

		switch n.Operator {
		case TokenPlus:
			return leftNumber + rightNumber, nil

		case TokenMinus:
			return leftNumber - rightNumber, nil

		case TokenMultiply:
			return leftNumber * rightNumber, nil

		case TokenDivide:
			return leftNumber / rightNumber, nil
		}

		return nil, fmt.Errorf("Unsupported binary operator")

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

	return util.StringValue(value)
}
