package cli

import (
	"strings"
)

func Confirm(message string) bool {
	for {
		if answer, err := Ask(message); err != nil {
			continue

		} else {
			answer = strings.ToUpper(strings.TrimSpace(answer))

			switch answer {
			case "YES", "Y":
				return true
			case "NO", "N":
				return false
			}
		}
	}
}
