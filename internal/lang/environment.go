package lang

type Environment struct {
	Vars    map[string]string
	Prompts map[string]string
	Funcs   map[string]func(args []any) (any, error)
}

func NewEnvironment() *Environment {
	env := &Environment{
		Vars:    make(map[string]string),
		Prompts: make(map[string]string),
		Funcs:   make(map[string]func(args []any) (any, error)),
	}

	Registry.Assign(env)

	return env
}
