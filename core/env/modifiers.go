package env

import "mimic/core/util"

func capitalize(args []string) string {
	if len(args) == 0 {
		return ""
	}

	return util.Capitalize(args[0])
}

func proper(args []string) string {
	if len(args) == 0 {
		return ""
	}

	return util.Proper(args[0])
}

func lower(args []string) string {
	if len(args) == 0 {
		return ""
	}

	return util.ToLower(args[0])
}

func upper(args []string) string {
	if len(args) == 0 {
		return ""
	}

	return util.ToUpper(args[0])
}

var Modifiers = map[string]func(args []string) string{
	"capitalize": capitalize,
	"proper":     proper,
	"lower":      lower,
	"upper":      upper,
}
