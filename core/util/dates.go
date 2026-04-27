package util

import (
	"fmt"
	"strings"
	"time"
)

func FormatDateToken(currentTime time.Time, token string) string {
	switch token {
	case "YYYY":
		return fmt.Sprintf("%04d", currentTime.Year())
	case "YY":
		return fmt.Sprintf("%02d", currentTime.Year()%100)
	case "M":
		return fmt.Sprintf("%d", int(currentTime.Month()))
	case "MM":
		return fmt.Sprintf("%02d", int(currentTime.Month()))
	case "D":
		return fmt.Sprintf("%d", currentTime.Day())
	case "DD":
		return fmt.Sprintf("%02d", currentTime.Day())
	case "H":
		return fmt.Sprintf("%d", currentTime.Hour())
	case "HH":
		return fmt.Sprintf("%02d", currentTime.Hour())
	case "m":
		return fmt.Sprintf("%d", currentTime.Minute())
	case "mm":
		return fmt.Sprintf("%02d", currentTime.Minute())
	case "s":
		return fmt.Sprintf("%d", currentTime.Second())
	case "ss":
		return fmt.Sprintf("%02d", currentTime.Second())
	default:
		return token
	}
}

func FormatDate(pattern string) string {
	currentTime := time.Now()

	tokens := []string{
		"YYYY", "YY",
		"MM", "M",
		"DD", "D",
		"HH", "H",
		"mm", "m",
		"ss", "s",
	}

	result := pattern

	for _, token := range tokens {
		result = strings.ReplaceAll(result, token, FormatDateToken(currentTime, token))
	}

	return result
}
