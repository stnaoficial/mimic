package main

import (
	"flag"
	"fmt"
	"io"
	configCmd "mimic/cmd/config"
	initCmd "mimic/cmd/init"
	localCmd "mimic/cmd/local"
	remoteCmd "mimic/cmd/remote"
	updateCmd "mimic/cmd/update"
	"os"
)

var Version = "development"

const (
	helpFlagUsage         = "Print help (this message) and exit"
	printVersionFlagUsage = "Print version information and exit"
)

type Command interface {
	Name() string
	Parse(args []string)
	Setup()
	Validate()
	Run()
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: mimic [OPTION]... [COMMAND] [ARG]...\n")
	fmt.Fprintf(os.Stderr, "Mimic is a template management and interpretation library.\n")
	fmt.Fprintf(os.Stderr, "\nCommands:\n")
	fmt.Fprintf(os.Stderr, "  config    %s\n", configCmd.CommandDescription)
	fmt.Fprintf(os.Stderr, "  init      %s\n", initCmd.CommandDescription)
	fmt.Fprintf(os.Stderr, "  local     %s\n", localCmd.CommandDescription)
	fmt.Fprintf(os.Stderr, "  remote    %s\n", remoteCmd.CommandDescription)
	fmt.Fprintf(os.Stderr, "  update    %s\n", updateCmd.CommandDescription)
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	fmt.Fprintf(os.Stderr, "  -h, --help    %s\n", helpFlagUsage)
	fmt.Fprintf(os.Stderr, "  --version     %s\n", printVersionFlagUsage)
	fmt.Fprintln(os.Stderr)
}

func run(args []string) {
	var printVersion bool

	flagSet := flag.NewFlagSet("mimic", flag.ExitOnError)
	flagSet.SetOutput(io.Discard)

	flagSet.Usage = usage
	flagSet.BoolVar(&printVersion, "version", false, printVersionFlagUsage)

	flagSet.Parse(args)

	if printVersion {
		fmt.Printf("Mimic version %s\n\n", Version)
		os.Exit(0)
	}

	usage()

	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	commands := []Command{
		configCmd.NewCommand("config"),
		initCmd.NewCommand("init"),
		localCmd.NewCommand("local"),
		remoteCmd.NewCommand("remote"),
		updateCmd.NewCommand("update"),
	}

	for _, command := range commands {
		if command.Name() == os.Args[1] {
			command.Parse(os.Args[2:])
			command.Setup()
			command.Validate()
			command.Run()
			os.Exit(0)
		}
	}

	run(os.Args[1:])
	os.Exit(0)
}
