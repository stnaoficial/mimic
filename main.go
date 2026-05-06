package main

import (
	"flag"
	"fmt"
	"io"
	"mimic/core"
	"mimic/core/cli"
	"mimic/core/lang"
	"os"
)

var Version = "development"

const (
	varFlagUsage = "Set a variable directly by passing a key=value pair"

	exprOpenFlagUsage  = "Set the open expression syntax (default \"{{\")"
	exprCloseFlagUsage = "Set the close expression syntax (default \"}}\")"
	debugModeFlagUsage = "Enable debug mode (default false)"

	helpFlagUsage         = "Print Help (this message) and exit"
	printVersionFlagUsage = "Print version information and exit"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: mimic [OPTION]... SOURCE TARGET\n")
	fmt.Fprintf(os.Stderr, "Mimic interpret files and directories in the source path (./.mimic directory by default) and generate copies of them in the target path (the current directory by default).\n\n")
	fmt.Fprintf(os.Stderr, "Provide variables directly\n")
	fmt.Fprintf(os.Stderr, "  -v, --var       %s\n\n", varFlagUsage)
	fmt.Fprintf(os.Stderr, "Configure how to start mimicking\n")
	fmt.Fprintf(os.Stderr, "  --expr-open     %s\n", exprOpenFlagUsage)
	fmt.Fprintf(os.Stderr, "  --expr-close    %s\n", exprCloseFlagUsage)
	fmt.Fprintf(os.Stderr, "  --debug     	  %s\n\n", debugModeFlagUsage)
	fmt.Fprintf(os.Stderr, "Get more information\n")
	fmt.Fprintf(os.Stderr, "  -h, --help      %s\n", helpFlagUsage)
	fmt.Fprintf(os.Stderr, "  --version       %s\n\n", printVersionFlagUsage)
}

func main() {
	config := core.NewConfig(Version)

	flag.Usage = usage

	flag.Var(&config.Variables, "v", varFlagUsage)
	flag.Var(&config.Variables, "var", varFlagUsage)

	flag.StringVar(&config.ExprOpen, "expr-open", lang.DefaultOpenExpr, exprOpenFlagUsage)
	flag.StringVar(&config.ExprClose, "expr-close", lang.DefaultCloseExpr, exprCloseFlagUsage)

	flag.BoolVar(&config.DebugMode, "debug", false, debugModeFlagUsage)

	flag.BoolVar(&config.PrintVersion, "version", false, printVersionFlagUsage)

	flag.CommandLine.SetOutput(io.Discard)

	flag.Parse()
	config.Parse()

	executor := core.NewExecutor(config)

	executor.Scan()
	executor.Generate()

	if !cli.MustConfirmToContinue() {
		os.Exit(0)
	}

	executor.Write()

	os.Exit(0)
}
