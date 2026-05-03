package env

import "mimic/core/util"

func camel(args []string) string {
	if len(args) == 0 {
		return ""
	}

	return util.ToCamel(args[0])
}

func pascal(args []string) string {
	if len(args) == 0 {
		return ""
	}

	return util.ToPascal(args[0])
}

func snake(args []string) string {
	if len(args) == 0 {
		return ""
	}

	return util.ToSnake(args[0])
}

func kebab(args []string) string {
	if len(args) == 0 {
		return ""
	}

	return util.ToKebab(args[0])
}

func dot(args []string) string {
	if len(args) == 0 {
		return ""
	}

	return util.ToDot(args[0])
}

func flat(args []string) string {
	if len(args) == 0 {
		return ""
	}

	return util.ToFlat(args[0])
}

var Cases = map[string]func(args []string) string{
	"camel":  camel,
	"pascal": pascal,
	"snake":  snake,
	"kebab":  kebab,
	"dot":    dot,
	"flat":   flat,
}
