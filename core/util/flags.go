package util

import (
	"fmt"
	"strings"
)

type FlagMap map[string]string

func (mf *FlagMap) String() string {
	return fmt.Sprint(*mf)
}

func (mf *FlagMap) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)

	if len(parts) != 2 {
		return fmt.Errorf("invalid format, expected key=value")
	}

	key := strings.TrimSpace(parts[0])
	val := strings.TrimSpace(parts[1])

	(*mf)[key] = val

	return nil
}
