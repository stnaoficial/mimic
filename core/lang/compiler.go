package lang

import (
	"strings"
)

type Compiler struct {
	Evaluator *Evaluator
	Env       *Environment
	Expr      *Expression
}

func NewCompiler(env *Environment, expr *Expression) *Compiler {
	eval := NewEvaluator(env)

	return &Compiler{
		Evaluator: eval,
		Env:       env,
		Expr:      expr,
	}
}

func (i *Compiler) Compile(buffer *Buffer) string {
	lexer := NewLexer(buffer, i.Expr)

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

			value := i.Evaluator.Eval(ast)

			result.WriteString(value)
		}
	}

	return result.String()
}
