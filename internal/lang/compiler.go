package lang

import (
	"strings"
)

type Compiler struct {
	Env       *Environment
	evaluator *Evaluator
	expr      *Expression
}

func NewCompiler(strict bool) *Compiler {
	env := NewEnvironment()

	return &Compiler{
		Env:       env,
		evaluator: NewEvaluator(env, strict),
		expr:      NewExpressionConfigurable(DefaultOpenExpr, DefaultCloseExpr),
	}
}

func NewCompilerConfigurable(env *Environment, expr *Expression, strict bool) *Compiler {
	return &Compiler{
		Env:       env,
		evaluator: NewEvaluator(env, strict),
		expr:      expr,
	}
}

func (c *Compiler) Compile(buffer *Buffer) (string, error) {
	lexer := NewLexer(buffer, c.expr)

	var result strings.Builder

	for {
		token, err := lexer.Next()

		if err != nil {
			return "", err
		}

		if token.Type == TokenEOF {
			break
		}

		if token.Type == TokenRaw {
			result.WriteString(token.Value)
			continue
		}

		if token.Type == TokenOpenExpr {
			parser := NewParser(lexer)

			node, err := parser.Parse()

			if err != nil {
				return "", err
			}

			value, err := c.evaluator.EvalAsString(node)

			if err != nil {
				return "", err
			}

			result.WriteString(value)
		}
	}

	return result.String(), nil
}
