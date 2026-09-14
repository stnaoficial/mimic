package config

import (
	"flag"
	"fmt"
	"io"
	"mimic/cmd"
	"mimic/internal/cli"
	"os"
)

var LocalConfig *cmd.Config

const (
	CommandDescription = "Display the current configuration settings defined in .mimic/config"
)

const (
	CommandDebugModeFlagUsage = "Enable debug mode (default false)"
)

func CommandUsage() {
	fmt.Fprintf(os.Stderr, "Usage: mimic config [OPTION]...\n")
	fmt.Fprintf(os.Stderr, "%s\n", CommandDescription)
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	fmt.Fprintf(os.Stderr, "  --debug     	  %s\n", CommandDebugModeFlagUsage)
	fmt.Fprintln(os.Stderr)
}

type CommandConfig struct {
	DebugMode bool
}

type Command struct {
	name string

	FlagSet *flag.FlagSet

	config *CommandConfig
}

func NewCommandConfig() *CommandConfig {
	return &CommandConfig{}
}

func NewCommand(name string) *Command {
	config := NewCommandConfig()

	flagSet := flag.NewFlagSet(name, flag.ExitOnError)
	flagSet.SetOutput(io.Discard)

	flagSet.Usage = CommandUsage

	flagSet.BoolVar(&config.DebugMode, "debug", false, CommandDebugModeFlagUsage)

	return &Command{
		name: name,

		FlagSet: flagSet,

		config: config,
	}
}

func (c *Command) Name() string {
	return c.name
}

func (c *Command) Parse(args []string) {
	c.FlagSet.Parse(args)
}

func (c *Command) Setup() {
	localConfig, err := cmd.NewLocalConfig()

	if err != nil {
		if c.config.DebugMode {
			cli.Logln(cli.LogSeverityError, err.Error())
		}

		fmt.Printf("Unable to load configuration\n\n")

		os.Exit(1)
	}

	LocalConfig = localConfig
}

func (c *Command) Validate() {}

func (c *Command) Run() {
	fmt.Print(string(LocalConfig.Bytes()))
}
