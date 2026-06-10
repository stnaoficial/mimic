package util

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func Normalize(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

	result, _, err := transform.String(t, s)

	if err != nil {
		return s
	}

	return result
}

func Join(s, sep string) string {
	return strings.Join(strings.Fields(s), sep)
}

func RandDigit(size int) string {
	var max uint64 = 1

	for range size {
		max *= 10
	}

	return fmt.Sprintf("%0*d", size, rand.Uint64()%max)
}
