package env

import (
	"strings"
)

func replace(args []string) string {
	if len(args) != 3 {
		return ""
	}

	return strings.ReplaceAll(args[0], args[1], args[2])
}

var Utils = map[string]func(args []string) string{
	"replace": replace,
}
