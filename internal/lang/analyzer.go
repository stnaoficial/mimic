package lang

type Analyzer struct {
	Env  *Environment
	expr *Expression
}

type Analisys any

type VariableAnalisys struct {
	Name  string
	Value string
	Ok    bool
}

type AnalisysMap = map[string]Analisys

func NewVariableAnalisys(name string, value string, ok bool) VariableAnalisys {
	return VariableAnalisys{
		Name:  name,
		Value: value,
		Ok:    ok,
	}
}

func NewAnalyzer() *Analyzer {
	env := NewEnvironment()

	return &Analyzer{
		Env:  env,
		expr: NewExpressionConfigurable(DefaultOpenExpr, DefaultCloseExpr),
	}
}

func NewAnalyzerConfigurable(env *Environment, expr *Expression) *Analyzer {
	return &Analyzer{
		Env:  env,
		expr: expr,
	}
}

func (d *Analyzer) Analyze(buffer *Buffer) (AnalisysMap, error) {
	lexer := NewLexer(buffer, d.expr)

	analisysMap := make(AnalisysMap)

	for {
		token, err := lexer.Next()

		if err != nil {
			return nil, err
		}

		if token.Type == TokenEOF {
			break
		}

		if token.Type == TokenRaw {
			continue
		}

		if token.Type == TokenOpenExpr {
			parser := NewParser(lexer)

			node, err := parser.Parse()

			if err != nil {
				return nil, err
			}

			d.analyze(analisysMap, node)
		}
	}

	return analisysMap, nil
}

func (d *Analyzer) analyze(analisysMap AnalisysMap, node Node) {
	switch n := node.(type) {
	case Identifier:
		name := n.Name
		value, ok := d.Env.Vars[name]

		analisysMap[name] = NewVariableAnalisys(name, value, ok)

	case Unary:
		d.analyze(analisysMap, n.Value)

	case Binary:
		d.analyze(analisysMap, n.Left)
		d.analyze(analisysMap, n.Right)

	case Callable:
		for _, arg := range n.Args {
			d.analyze(analisysMap, arg)
		}
	}
}
