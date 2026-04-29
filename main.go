package main

import (
	"flag"
	"fmt"
	"io"
	"maps"
	"mimic/core"
	"mimic/core/cli"
	"mimic/core/lang"
	"mimic/core/util"
	"os"
)

const varFlagUsage = "Set a variable directly by passing a key=value pair"

const exprOpenFlagUsage = "Set the open expression syntax (default \"{{\")"
const exprCloseFlagUsage = "Set the close expression syntax (default \"}}\")"

const helpFlagUsage = "Print Help (this message) and exit"
const versionFlagUsage = "Print version information and exit"

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: mimic [OPTION]... SOURCE TARGET\n")
	fmt.Fprintf(os.Stderr, "Mimic interpret .mimic files in the source path (./.mimic directory by default) and generate copies of them in the target path (the current directory by default).\n\n")
	fmt.Fprintf(os.Stderr, "Provide variables directly\n")
	fmt.Fprintf(os.Stderr, "  -v, --var       %s\n\n", varFlagUsage)
	fmt.Fprintf(os.Stderr, "Configure how to start mimicking files in the source path\n")
	fmt.Fprintf(os.Stderr, "  --expr-open     %s\n", exprOpenFlagUsage)
	fmt.Fprintf(os.Stderr, "  --expr-close    %s\n\n", exprCloseFlagUsage)
	fmt.Fprintf(os.Stderr, "Get more information\n")
	fmt.Fprintf(os.Stderr, "  -h, --help      %s\n", helpFlagUsage)
	fmt.Fprintf(os.Stderr, "  --version       %s\n\n", versionFlagUsage)
}

func versionFlag() *bool {
	return flag.Bool("version", false, versionFlagUsage)
}

func variableFlag() util.FlagMap {
	vars := make(util.FlagMap)

	flag.Var(&vars, "v", varFlagUsage)
	flag.Var(&vars, "var", varFlagUsage)

	return vars
}

func expressionFlag() (*string, *string) {
	exprOpen := flag.String("expr-open", "{{", exprOpenFlagUsage)
	exprClose := flag.String("expr-close", "}}", exprCloseFlagUsage)
	return exprOpen, exprClose
}

func parseFlags() {
	flag.CommandLine.SetOutput(io.Discard)
	flag.Parse()
}

func printVersionAndExit() {
	fmt.Printf("Mimic version 1.0.0-beta\n")
	os.Exit(0)
}

func getSourceAndTarget() (string, string) {
	args := flag.Args()

	source := "./.mimic"
	target := "."

	if len(args) >= 1 {
		source = args[0]
	}

	if len(args) >= 2 {
		target = args[1]
	}

	if len(args) > 2 {
		flag.Usage()
		os.Exit(1)
	}

	return source, target
}

func getEnvironment(vars util.FlagMap) *lang.Environment {
	env := lang.NewEnvironment()

	maps.Copy(env.Vars, vars)

	return env
}

func main() {
	flag.Usage = usage

	version := versionFlag()
	vars := variableFlag()
	exprOpen, exprClose := expressionFlag()

	parseFlags()

	if *version {
		printVersionAndExit()
	}

	source, target := getSourceAndTarget()

	env := lang.NewEnvironment()
	maps.Copy(env.Vars, vars)

	expr := lang.NewExpression(*exprOpen, *exprClose)

	comp := lang.NewCompiler(env, expr)

	executor := core.NewExecutor(source, target, comp)

	executor.Read()

	for filename, _ := range executor.FilesRead {
		cli.LogFileNameAt(filename)
	}

	if !cli.MustConfirmToContinue() {
		os.Exit(0)
	}

	executor.Write()

	for filename, filedata := range executor.WrittenFiles {
		cli.LogFileNameAt(filename)
		cli.LogFileDataAdded(filedata)
	}

	os.Exit(0)
}
