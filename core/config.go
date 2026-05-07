package core

import (
	"flag"
	"fmt"
	"mimic/core/util"
	"os"
)

type Config struct {
	SourcePath string
	TargetPath string

	Variables util.FlagMap
	Prompts   util.FlagMap

	ExprOpen  string
	ExprClose string

	PrintVersion bool
	version      string

	DebugMode bool
}

func NewConfig(version string) *Config {
	return &Config{
		SourcePath: "./.mimic",
		TargetPath: ".",

		Variables: make(util.FlagMap),
		Prompts:   make(util.FlagMap),

		ExprOpen:  "{{",
		ExprClose: "}}",

		PrintVersion: false,
		version:      version,

		DebugMode: false,
	}
}

func (c *Config) Parse() {
	c.parseArgs()

	if c.PrintVersion {
		fmt.Printf("Mimic version %s\n", c.version)
		os.Exit(0)
	}
}

func (c *Config) parseArgs() {
	args := flag.Args()

	if len(args) >= 1 {
		c.SourcePath = args[0]
	}

	if len(args) >= 2 {
		c.TargetPath = args[1]
	}

	if len(args) > 2 {
		flag.Usage()
		os.Exit(1)
	}
}
