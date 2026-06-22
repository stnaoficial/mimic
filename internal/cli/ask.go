package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Prompt(question string) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(question)

	answer, err := reader.ReadString('\n')

	return strings.TrimSpace(answer), err
}

func MustPrompt(question string) string {
	for {
		if answer, err := Prompt(question); err != nil {
			continue
		} else {
			answer = strings.TrimSpace(answer)

			if len(answer) == 0 {
				continue
			}

			return answer
		}
	}
}
