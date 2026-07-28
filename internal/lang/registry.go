package lang

import (
	"mimic/internal/util"
	"strings"
	"unicode"
)

var Registry = make(map[string]any)

func init() {
	Registry["upper"] = func(s string) string {
		return strings.ToUpper(s)
	}

	Registry["lower"] = func(s string) string {
		return strings.ToLower(s)
	}

	Registry["proper"] = func(s string) string {
		tokens := strings.Fields(strings.ToLower(s))

		if len(tokens) == 0 {
			return ""
		}

		for i := range tokens {
			runes := []rune(tokens[i])

			runes[0] = unicode.ToUpper(runes[0])

			tokens[i] = string(runes)
		}

		return strings.Join(tokens, " ")
	}

	Registry["title"] = func(s string) string {
		tokens := strings.Fields(strings.ToLower(s))

		if len(tokens) == 0 {
			return ""
		}

		for i := range tokens {
			runes := []rune(tokens[i])

			if len(runes) > 2 || (i == 0 || i == len(tokens)-1) {
				runes[0] = unicode.ToUpper(runes[0])
			}

			tokens[i] = string(runes)
		}

		return strings.Join(tokens, " ")
	}

	Registry["capitalize"] = func(s string) string {
		if s == "" {
			return ""
		}

		runes := []rune(strings.ToLower(s))

		runes[0] = unicode.ToUpper(runes[0])

		return string(runes)
	}

	Registry["pascal"] = func(s string) string {
		tokens := strings.Fields(strings.ToLower(util.Normalize(s)))

		if len(tokens) == 0 {
			return ""
		}

		var builder strings.Builder

		for i := range tokens {
			runes := []rune(tokens[i])

			runes[0] = unicode.ToUpper(runes[0])

			builder.WriteString(string(runes))
		}

		return builder.String()
	}

	Registry["camel"] = func(s string) string {
		tokens := strings.Fields(strings.ToLower(util.Normalize(s)))

		if len(tokens) == 0 {
			return ""
		}

		var builder strings.Builder

		builder.WriteString(tokens[0])

		for i := 1; i < len(tokens); i++ {
			runes := []rune(tokens[i])

			runes[0] = unicode.ToUpper(runes[0])

			builder.WriteString(string(runes))
		}

		return builder.String()
	}

	Registry["flat"] = func(s string) string {
		return strings.Join(strings.Fields(strings.ToLower(util.Normalize(s))), "")
	}

	Registry["snake"] = func(s string) string {
		return util.Delimit(strings.ToLower(util.Normalize(s)), "_")
	}

	Registry["kebab"] = func(s string) string {
		return util.Delimit(strings.ToLower(util.Normalize(s)), "-")
	}

	Registry["before"] = func(value string, target string) string {
		before, _, ok := strings.Cut(value, target)

		if !ok {
			return value
		}

		return before
	}

	Registry["after"] = func(value string, target string) string {
		_, after, ok := strings.Cut(value, target)

		if !ok {
			return value
		}

		return after
	}

	Registry["between"] = func(value string, first string, last string) string {
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

	Registry["replace"] = func(str, old, new string) string {
		return strings.ReplaceAll(str, old, new)
	}

	Registry["normalize"] = func(str string) string {
		return util.Normalize(str)
	}

	Registry["delimit"] = func(str, del string) string {
		return util.Delimit(str, del)
	}

	Registry["pad_left"] = func(value any, length int, pad string) (string, error) {
		str, err := util.StringValue(value)
		return util.Pad(str, length, pad, util.PadLeft), err
	}

	Registry["pad_right"] = func(value any, length int, pad string) (string, error) {
		str, err := util.StringValue(value)
		return util.Pad(str, length, pad, util.PadRight), err
	}

	Registry["pad_both"] = func(value any, length int, pad string) (string, error) {
		str, err := util.StringValue(value)
		return util.Pad(str, length, pad, util.PadBoth), err
	}

	Registry["zero_fill"] = func(value any, length int) (string, error) {
		str, err := util.StringValue(value)
		return util.ZeroFill(str, length), err
	}
}
