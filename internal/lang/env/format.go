package env

import (
	"mimic/internal/util"
	"strings"
	"unicode"
)

func Upper(s string) string {
	return strings.ToUpper(s)
}

func Lower(s string) string {
	return strings.ToLower(s)
}

func Proper(s string) string {
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

func Title(s string) string {
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

func Capitalize(s string) string {
	if s == "" {
		return ""
	}

	runes := []rune(strings.ToLower(s))

	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

func Pascal(s string) string {
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

func Camel(s string) string {
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

func Flat(s string) string {
	return strings.Join(strings.Fields(strings.ToLower(util.Normalize(s))), "")
}
