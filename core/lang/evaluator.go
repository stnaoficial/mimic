package lang

import (
	"fmt"
	"mimic/core/cli"
)

type Evaluator struct {
	Env *Environment
}

func NewEvaluator(env *Environment) *Evaluator {
	return &Evaluator{
		Env: env,
	}
}

func (e *Evaluator) abort(cause string) {
	cli.LogAndExit(cause, cli.LogSeverityError)
}

func (e *Evaluator) Eval(node Node) string {
	switch n := node.(type) {

	case Identifier:
		if value, ok := e.Env.Vars[n.Name]; ok {
			return value
		}

		value := cli.MustAsk(fmt.Sprintf("Please enter a value for \"%s\": ", n.Name))

		e.Env.Vars[n.Name] = value

		return value

	case StringLiteral:
		return n.Value

	case CallExpression:
		fn, ok := e.Env.Funcs[n.Name]

		if !ok {
			e.abort(fmt.Sprintf("Unexpected function call \"%s\"", n.Name))
		}

		var args []string

		for _, arg := range n.Args {
			args = append(args, e.Eval(arg))
		}

		return fn(args)
	}

	return ""
}
