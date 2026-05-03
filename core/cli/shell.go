package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/google/shlex"
)

func Shell(input string) (string, error) {
	args, err := shlex.Split(input)

	if err != nil {
		return "", err
	}

	if len(args) == 0 {
		return "", fmt.Errorf("Empty argument")
	}

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Dir = os.TempDir()
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return string(output), nil
}
