package env

import (
	"mimic/internal/util"
	"strings"
)

func Before(value string, target string) string {
	before, _, ok := strings.Cut(value, target)

	if !ok {
		return value
	}

	return before
}

func After(value string, target string) string {
	_, after, ok := strings.Cut(value, target)

	if !ok {
		return value
	}

	return after
}

func Between(value string, first string, last string) string {
	start := strings.Index(value, first)

	if start == -1 {
		return value
	}

	start += len(first)

	end := strings.Index(value[start:], last)

	if end == -1 {
		return value
	}

	return value[start : start+end]
}

func Replace(str, old, new string) string {
	return strings.ReplaceAll(str, old, new)
}

func Normalize(str string) string {
	return util.Normalize(str)
}

func Delimit(str, del string) string {
	return util.Delimit(str, del)
}

func PadLeft(value any, length int, pad string) (string, error) {
	str, err := util.StringValue(value)
	return util.Pad(str, length, pad, util.PadLeft), err
}

func PadRight(value any, length int, pad string) (string, error) {
	str, err := util.StringValue(value)
	return util.Pad(str, length, pad, util.PadRight), err
}

func PadBoth(value any, length int, pad string) (string, error) {
	str, err := util.StringValue(value)
	return util.Pad(str, length, pad, util.PadBoth), err
}

func ZeroFill(value any, length int) (string, error) {
	str, err := util.StringValue(value)
	return util.ZeroFill(str, length), err
}
