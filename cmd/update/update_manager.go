package update

import (
	"fmt"
	"os"
	"os/exec"
)

type UpdateManager struct{}

func NewUpdateManager() *UpdateManager {
	return &UpdateManager{}
}

func (um *UpdateManager) Update(version string) error {
	command := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf(
			"curl -fsSL %s | sh -s %s",
			GitHubApiClient.ParseRawApiEntryUrl("install.sh"),
			version,
		),
	)

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}
