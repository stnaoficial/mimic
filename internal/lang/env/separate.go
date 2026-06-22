package env

import (
	"mimic/internal/util"
	"strings"
)

func Snake(s string) string {
	return util.Delimit(strings.ToLower(util.Normalize(s)), "_")
}

func Kebab(s string) string {
	return util.Delimit(strings.ToLower(util.Normalize(s)), "-")
}
