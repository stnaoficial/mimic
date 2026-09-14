package init

import (
	"fmt"
	"mimic/internal/cli"
	"os"
)

type initializer struct {
	debug bool
}

func NewInitializer(debug bool) *initializer {
	return &initializer{
		debug: debug,
	}
}

func (i *initializer) Init() error {
	if startErr := i.start(); startErr != nil {
		if abortErr := i.abort(); abortErr != nil {
			// allow debug
			if i.debug {
				cli.Logf(cli.LogSeverityError, abortErr.Error())
			}

			fmt.Printf("Could not abort changes in the current path\n\n")
		}

		return startErr
	}

	fmt.Printf("Initialized an empty .mimic directory in the current path\n\n")

	return nil
}

func (i *initializer) start() error {
	if err := os.MkdirAll(".mimic/templates", 0755); err != nil {
		return err
	}

	if err := os.WriteFile(".mimic/config", LocalConfig.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

func (i *initializer) abort() error {
	return os.RemoveAll(".mimic")
}
