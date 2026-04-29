package util

import (
	"fmt"
	"strings"
)

type FlagMap map[string]string

func (f *FlagMap) String() string {
	return fmt.Sprint(*f)
}

func (f *FlagMap) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)

	if len(parts) != 2 {
		return fmt.Errorf("invalid format, expected key=value")
	}

	key := strings.TrimSpace(parts[0])
	val := strings.TrimSpace(parts[1])

	(*f)[key] = val

	return nil
}
