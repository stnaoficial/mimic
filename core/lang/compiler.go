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
		evaluator: NewEvaluator(env),
		expr:      expr,
	}
}

func (c *Compiler) Compile(buffer *Buffer) string {
	lexer := NewLexer(buffer, c.expr)

	var result strings.Builder

	for {
		token := lexer.Next()

		if token.Type == TokenEOF {
			break
		}

		if token.Type == TokenRaw {
			result.WriteString(token.Value)
			continue
		}

		if token.Type == TokenOpenExpr {
			parser := NewParser(lexer)

			ast := parser.Parse()

			value := c.evaluator.Eval(ast)

			result.WriteString(value)
		}
	}

	return result.String()
}
