package util

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	lowerRegex  = regexp.MustCompile(`^[^A-Z]*[a-z][^A-Z]*$`)
	camelRegex  = regexp.MustCompile(`^[a-z]+(?:[A-Z][a-z0-9]*)+$`)
	pascalRegex = regexp.MustCompile(`^(?:[A-Z][a-z0-9]*)+$`)
)

func capitalize(w string) string {
	if w == "" {
		return ""
	}

	runes := []rune(w)

	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

func isSeparator(r rune) bool {
	if unicode.IsSpace(r) {
		return true
	}

	if r == '-' || r == '_' || r == '.' {
		return true
	}

	return false
}

func tokenize(s string) []string {
	var tokens []string
	var current []rune

	for i, r := range s {
		if isSeparator(r) {
			if len(current) > 0 {
				tokens = append(tokens, strings.ToLower(string(current)))
				current = []rune{}
			}
			continue
		}

		if i > 0 && unicode.IsUpper(r) && (len(current) > 0 && !unicode.IsUpper(rune(s[i-1]))) {
			tokens = append(tokens, strings.ToLower(string(current)))
			current = []rune{}
		}

		current = append(current, r)
	}

	if len(current) > 0 {
		tokens = append(tokens, strings.ToLower(string(current)))
	}

	return tokens
}

func IsLower(s string) bool {
	return lowerRegex.MatchString(s)
}

func IsCamel(s string) bool {
	return camelRegex.MatchString(s)
}

func IsPascal(s string) bool {
	return pascalRegex.MatchString(s)
}

func ToKebab(s string) string {
	return strings.Join(tokenize(s), "-")
}

func ToSnake(s string) string {
	return strings.Join(tokenize(s), "_")
}

func ToDot(s string) string {
	return strings.Join(tokenize(s), ".")
}

func ToLower(s string) string {
	return strings.Join(tokenize(s), "")
}

func ToUpper(s string) string {
	return strings.ToUpper(strings.Join(tokenize(s), ""))
}

func ToCamel(s string) string {
	tokens := tokenize(s)

	if len(tokens) == 0 {
		return ""
	}

	var builder strings.Builder

	builder.WriteString(tokens[0])

	for i := 1; i < len(tokens); i++ {
		builder.WriteString(capitalize(tokens[i]))
	}

	return builder.String()
}

func ToPascal(s string) string {
	tokens := tokenize(s)

	if len(tokens) == 0 {
		return ""
	}

	var builder strings.Builder

	for i := 0; i < len(tokens); i++ {
		builder.WriteString(capitalize(tokens[i]))
	}

	return builder.String()
}
