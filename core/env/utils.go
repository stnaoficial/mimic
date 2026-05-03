package env

import (
	"mimic/core/cli"
	"strings"
)

func shell(args []string) string {
	if len(args) != 1 {
		return ""
	}

	if output, err := cli.Shell(args[0]); err == nil {
		return output
	} else {
		return ""
	}
}

func replace(args []string) string {
	if len(args) != 3 {
		return ""
	}

	return strings.ReplaceAll(args[0], args[1], args[2])
}

var Utils = map[string]func(args []string) string{
	"$":       shell,
	"replace": replace,
}
