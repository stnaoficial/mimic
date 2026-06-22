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

const (
	PadLeft int = iota
	PadRight
	PadBoth
)

func Normalize(str string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

	result, _, err := transform.String(t, str)

	if err != nil {
		return str
	}

	return result
}

func Delimit(str, del string) string {
	return strings.Join(strings.Fields(str), del)
}

func RandDigit(size int) string {
	var max uint64 = 1

	for range size {
		max *= 10
	}

	return fmt.Sprintf("%0*d", size, rand.Uint64()%max)
}

func Pad(str string, length int, pad string, tp int) string {
	if length <= len(str) || pad == "" {
		return str
	}

	missing := length - len(str)

	repeat := func(n int) string {
		if n <= 0 {
			return ""
		}

		repeated := strings.Repeat(pad, (n+len(pad)-1)/len(pad))
		return repeated[:n]
	}

	switch tp {
	case PadLeft:
		return repeat(missing) + str

	case PadRight:
		return str + repeat(missing)

	case PadBoth:
		left := missing / 2
		right := missing - left

		return repeat(left) + str + repeat(right)

	default:
		return str
	}
}

func ZeroFill(str string, length int) string {
	return Pad(str, length, "0", PadLeft)
}
