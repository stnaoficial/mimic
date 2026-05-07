package lang

import (
	"strings"
)

type Compiler struct {
	Env       *Environment
	evaluator *Evaluator
	expr      *Expression
}

func NewCompiler() *Compiler {
	env := NewEnvironment()

	return &Compiler{
		Env:       env,
		evaluator: NewEvaluator(env),
		expr:      NewExpressionConfigurable(DefaultOpenExpr, DefaultCloseExpr),
	}
}

func NewCompilerConfigurable(env *Environment, expr *Expression) *Compiler {
	return &Compiler{
		Env:       env,
		evaluator: NewEvaluator(env),
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
			node, err := NewParser(lexer).Parse()

			if err != nil {
				return "", err
			}

			value, err := c.evaluator.Eval(node)

			if err != nil {
				return "", err
			}

			result.WriteString(value)
		}
	}

	return result.String(), nil
}
