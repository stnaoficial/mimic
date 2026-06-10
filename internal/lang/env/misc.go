package env

import (
	"mimic/internal/util"
	"strings"
)

func Replace(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

func Normalize(s string) string {
	return util.Normalize(s)
}

func Join(s, sep string) string {
	return util.Join(s, sep)
}
