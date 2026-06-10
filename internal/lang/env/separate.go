package env

import (
	"mimic/internal/util"
	"strings"
)

func Snake(s string) string {
	return util.Join(strings.ToLower(util.Normalize(s)), "_")
}

func Kebab(s string) string {
	return util.Join(strings.ToLower(util.Normalize(s)), "-")
}
