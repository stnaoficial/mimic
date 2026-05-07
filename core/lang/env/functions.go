package env

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func normalize(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

	result, _, err := transform.String(t, s)

	if err != nil {
		return s
	}

	return result
}

func join(s, sep string) string {
	return strings.Join(strings.Fields(s), sep)
}

// case formatters

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
	tokens := strings.Fields(strings.ToLower(normalize(s)))

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
	tokens := strings.Fields(strings.ToLower(normalize(s)))

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
	return strings.Join(strings.Fields(strings.ToLower(normalize(s))), "")
}

// case separators

func Snake(s string) string {
	return join(strings.ToLower(normalize(s)), "_")
}

func Kebab(s string) string {
	return join(strings.ToLower(normalize(s)), "-")
}

// miscellaneous

func Replace(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

func Normalize(s string) string {
	return normalize(s)
}

func Join(s, sep string) string {
	return join(s, sep)
}
