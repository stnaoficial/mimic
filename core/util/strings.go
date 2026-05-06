package util

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"unicode"
)

func RandDigit(size int) string {
	var max uint64 = 1

	for range size {
		max *= 10
	}

	return fmt.Sprintf("%0*d", size, rand.Uint64()%max)
}

func IsQuoted(s string) bool {
	if len(s) < 2 {
		return false
	}

	first := s[:1]
	last := s[len(s)-1:]

	return (first == "'" && last == "'") || (first == "\"" && last == "\"")
}

func Quote(s string) string {
	return strconv.Quote(s)
}

func Unquote(s string) string {
	value, err := strconv.Unquote(s)

	if err != nil {
		return s
	}

	return value
}

func Capitalize(s string) string {
	if s == "" {
		return ""
	}

	runes := []rune(strings.ToLower(s))

	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
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

func ToCamel(s string) string {
	tokens := strings.Fields(strings.ToLower(s))

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

func ToPascal(s string) string {
	tokens := strings.Fields(strings.ToLower(s))

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
