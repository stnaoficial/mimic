package util

import (
	"strings"
	"unicode"
)

func IsQuoted(s string) bool {
	if len(s) < 2 {
		return false
	}

	first := s[:1]
	last := s[len(s)-1:]

	return (first == SingleQuote && last == SingleQuote) || (first == DoubleQuote && last == DoubleQuote)
}

func Unquote(s string) string {
	return strings.Trim(s, SingleQuote+DoubleQuote)
}

func ToSentence(s string) string {
	if s == "" {
		return ""
	}

	runes := []rune(s)

	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

func ToCamel(s string) string {
	tokens := strings.Fields(s)

	if len(tokens) == 0 {
		return ""
	}

	var builder strings.Builder

	builder.WriteString(tokens[0])

	for i := 1; i < len(tokens); i++ {
		builder.WriteString(ToSentence(tokens[i]))
	}

	return builder.String()
}

func ToPascal(s string) string {
	tokens := strings.Fields(s)

	if len(tokens) == 0 {
		return ""
	}

	var builder strings.Builder

	for i := 0; i < len(tokens); i++ {
		builder.WriteString(ToSentence(tokens[i]))
	}

	return builder.String()
}

func ToSnake(s string) string {
	return strings.Join(strings.Fields(s), "_")
}

func ToKebab(s string) string {
	return strings.Join(strings.Fields(s), "-")
}

func ToDot(s string) string {
	return strings.Join(strings.Fields(s), ".")
}

func ToFlat(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), ""))
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func ToUpper(s string) string {
	return strings.ToUpper(s)
}
