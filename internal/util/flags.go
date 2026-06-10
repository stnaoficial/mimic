package util

import (
	"fmt"
	"strings"
)

type FlagMap map[string]string

func (m *FlagMap) String() string {
	return fmt.Sprint(*m)
}

func (m *FlagMap) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)

	if len(parts) != 2 {
		return fmt.Errorf("invalid format, expected key=value")
	}

	key := strings.TrimSpace(parts[0])
	val := strings.TrimSpace(parts[1])

	(*m)[key] = val

	return nil
}

type FlagSlice struct {
	Values     []string
	hasDefault bool
}

func NewFlagSlice(defaults ...string) FlagSlice {
	return FlagSlice{
		Values:     defaults,
		hasDefault: true,
	}
}

func (s *FlagSlice) String() string {
	return fmt.Sprint(s.Values)
}

func (s *FlagSlice) Set(value string) error {
	if s.hasDefault {
		s.Values = nil
		s.hasDefault = false
	}

	s.Values = append(s.Values, value)

	return nil
}
