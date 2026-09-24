package update

import (
	"flag"
	"fmt"
	"io"
	"mimic/cmd"
	"mimic/internal/cli"
	"mimic/tools/github"
	"os"
)

var GitHubApiClient *github.ApiClient

const (
	CommandDescription = "Update the library to a specific version (the latest version by default)"
)

const (
	CommandDebugModeFlagUsage = "Enable debug mode (default false)"
)

type CommandConfig struct {
	DebugMode bool
	Version   string
}

type Command struct {
	name string

	FlagSet *flag.FlagSet

	config *CommandConfig
}

func NewCommandConfig() *CommandConfig {
	return &CommandConfig{}
}

func CommandUsage() {
	fmt.Fprintf(os.Stderr, "Usage: mimic update [VERSION] [OPTION]...\n")
	fmt.Fprintf(os.Stderr, "%s\n", CommandDescription)
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	fmt.Fprintf(os.Stderr, "  --debug    %s\n", CommandDebugModeFlagUsage)
	fmt.Fprintln(os.Stderr)
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

	if c.FlagSet.NArg() > 0 {
		c.config.Version = c.FlagSet.Arg(0)
	}
}

func (c *Command) Setup() {
	localConfig, err := cmd.NewLocalConfig()

	if err != nil {
		if c.config.DebugMode {
			cli.Logln(cli.LogSeverityError, err.Error())
		}

		fmt.Printf("Unable to start operation\n\n")

		os.Exit(1)
	}

	localConfig.RegisterDefaultSettings()

	GitHubApiClient = github.NewApiClient(
		localConfig.Require("update.github.user.username"),
		localConfig.Require("update.github.repository.name"),
		localConfig.Require("update.github.repository.branch.name"),
		localConfig.Require("update.github.user.access-token"),
	)
}

func (c *Command) Validate() {}

func (c *Command) Run() {
	updateManager := NewUpdateManager()

	if err := updateManager.Update(c.config.Version); err != nil {
		if c.config.DebugMode {
			cli.Logln(cli.LogSeverityError, err.Error())
		}

		fmt.Printf("Unable to install update\n\n")

		os.Exit(1)
	}
}
