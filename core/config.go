package core

import (
	"mimic/core/util"
)

type Config struct {
	SourcePath string
	TargetPath string

	Variables util.FlagMap
	Prompts   util.FlagMap

	ExprOpen  string
	ExprClose string

	Init      bool
	DebugMode bool
}

func NewConfig() *Config {
	return &Config{
		SourcePath: "./.mimic",
		TargetPath: ".",

		Variables: make(util.FlagMap),
		Prompts:   make(util.FlagMap),

		ExprOpen:  "{{",
		ExprClose: "}}",

		Init:      false,
		DebugMode: false,
	}
}
