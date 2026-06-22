package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	InitCommandDescription = "Initialize an empty .mimic directory in the current path and exit"
)

type InitCommand struct {
	FlagSet *flag.FlagSet
}

func InitCommandUsage() {
	fmt.Fprintf(os.Stderr, "Usage: mimic init\n")
	fmt.Fprintf(os.Stderr, "%s\n", InitCommandDescription)
	fmt.Fprintln(os.Stderr)
}

func NewInitCommand() *InitCommand {
	flagSet := flag.NewFlagSet("init", flag.ExitOnError)
	flagSet.SetOutput(io.Discard)

	flagSet.Usage = InitCommandUsage

	return &InitCommand{
		FlagSet: flagSet,
	}
}

func (c *InitCommand) Run(args []string) error {
	c.FlagSet.Parse(args)

	if err := os.MkdirAll(".mimic", 0755); err != nil {
		fmt.Printf("Could not initialize an empty .mimic directory in the current path\n\n")
	} else {
		fmt.Printf("Initialized an empty .mimic directory in the current path\n\n")
	}

	return nil
}
