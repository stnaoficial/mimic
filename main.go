package main

import (
	"flag"
	"fmt"
	"mimic/internal/cmd"
	"os"
)

var Version = "development"

const (
	helpFlagUsage         = "Print help (this message) and exit"
	printVersionFlagUsage = "Print version information and exit"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: mimic [OPTION]... [COMMAND] [ARG]...\n")
	fmt.Fprintf(os.Stderr, "Mimic interprets files and directories from a source path (.mimic directory by default) and generates copies in a target path (the current directory by default).\n")
	fmt.Fprintf(os.Stderr, "\nCommands:\n")
	fmt.Fprintf(os.Stderr, "  init    %s\n", cmd.InitCommandDescription)
	fmt.Fprintf(os.Stderr, "  copy    %s\n", cmd.CopyCommandDescription)
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	fmt.Fprintf(os.Stderr, "  -h, --help    %s\n", helpFlagUsage)
	fmt.Fprintf(os.Stderr, "  --version     %s\n", printVersionFlagUsage)
	fmt.Fprintln(os.Stderr)
}

func run(args []string) {
	var printVersion bool

	flagSet := flag.NewFlagSet("mimic", flag.ExitOnError)

	flagSet.Usage = usage
	flagSet.BoolVar(&printVersion, "version", false, printVersionFlagUsage)

	flagSet.Parse(args)

	if printVersion {
		fmt.Printf("Mimic version %s\n", Version)
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

	switch os.Args[1] {
	case "init":
		cmd.NewInitCommand().Run(os.Args[2:])
	case "copy":
		cmd.NewCopyCommand().Run(os.Args[2:])
	default:
		run(os.Args[1:])
	}

	os.Exit(0)
}
