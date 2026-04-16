package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Ask(question string) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(question)

	answer, err := reader.ReadString('\n')

	return strings.TrimSpace(answer), err
}

func MustAsk(question string) string {
	for {
		if answer, err := Ask(question); err != nil {
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
