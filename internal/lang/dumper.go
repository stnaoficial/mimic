package lang

type Dumper struct {
	Env  *Environment
	expr *Expression
}

type Dump any

type VariableDump struct {
	Name  string
	Value string
	Ok    bool
}

type DumpMap = map[string]Dump

func NewVariableDump(name string, value string, ok bool) VariableDump {
	return VariableDump{
		Name:  name,
		Value: value,
		Ok:    ok,
	}
}

func NewDumper() *Dumper {
	env := NewEnvironment()

	return &Dumper{
		Env:  env,
		expr: NewExpressionConfigurable(DefaultOpenExpr, DefaultCloseExpr),
	}
}

func NewDumperConfigurable(env *Environment, expr *Expression) *Dumper {
	return &Dumper{
		Env:  env,
		expr: expr,
	}
}

func (d *Dumper) Dump(buffer *Buffer) (DumpMap, error) {
	lexer := NewLexer(buffer, d.expr)

	dumpMap := make(DumpMap)

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

			d.dump(dumpMap, node)
		}
	}

	return dumpMap, nil
}

func (d *Dumper) dump(dumpMap DumpMap, node Node) {
	switch n := node.(type) {
	case Identifier:
		name := n.Name
		value, ok := d.Env.Vars[name]

		dumpMap[name] = NewVariableDump(name, value, ok)

	case Unary:
		d.dump(dumpMap, n.Value)

	case Binary:
		d.dump(dumpMap, n.Left)
		d.dump(dumpMap, n.Right)

	case Callable:
		for _, arg := range n.Args {
			d.dump(dumpMap, arg)
		}
	}
}
