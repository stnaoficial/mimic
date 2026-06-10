package internal

import (
	"mimic/internal/util"
)

type Config struct {
	SourcePath util.FlagSlice
	TargetPath util.FlagSlice

	Variables util.FlagMap
	Prompts   util.FlagMap

	ExprOpen  string
	ExprClose string

	Init      bool
	DebugMode bool
}

func NewConfig() *Config {
	return &Config{
		SourcePath: util.NewFlagSlice(".mimic"),
		TargetPath: util.NewFlagSlice("."),

		Variables: make(util.FlagMap),
		Prompts:   make(util.FlagMap),
	}
}
