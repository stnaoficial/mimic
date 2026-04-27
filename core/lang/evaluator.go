package lang

import (
	"fmt"
	"mimic/core/cli"
)

type Evaluator struct {
	env *Environment
}

func NewEvaluator(env *Environment) *Evaluator {
	return &Evaluator{
		env: env,
	}
}

func (e *Evaluator) Eval(node Node) string {
	switch n := node.(type) {

	case Identifier:
		if value, ok := e.env.Vars[n.Name]; ok {
			return value
		}

		value := cli.MustAsk(fmt.Sprintf("Please enter a value for \"%s\": ", n.Name))

		e.env.Vars[n.Name] = value

		return value

	case StringLiteral:
		return n.Value

	case CallExpression:
		fn, ok := e.env.Funcs[n.Name]

		if !ok {
			cli.LogAndExit(fmt.Sprintf("Unexpected function call \"%s\"", n.Name), cli.LogSeverityError)
		}

		var args []string

		for _, arg := range n.Args {
			args = append(args, e.Eval(arg))
		}

		return fn(args)
	}

	return ""
}
