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

const (
	varFlagUsage = "Set a variable directly by passing a key=value pair"

	exprOpenFlagUsage  = "Set the open expression syntax (default \"{{\")"
	exprCloseFlagUsage = "Set the close expression syntax (default \"}}\")"

	helpFlagUsage    = "Print Help (this message) and exit"
	versionFlagUsage = "Print version information and exit"
)

var Version = "development"

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: mimic [OPTION]... SOURCE TARGET\n")
	fmt.Fprintf(os.Stderr, "Mimic interpret files and directories in the source path (./.mimic directory by default) and generate copies of them in the target path (the current directory by default).\n\n")
	fmt.Fprintf(os.Stderr, "Provide variables directly\n")
	fmt.Fprintf(os.Stderr, "  -v, --var       %s\n\n", varFlagUsage)
	fmt.Fprintf(os.Stderr, "Configure how to start mimicking\n")
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
	exprOpen := flag.String("expr-open", lang.DefaultOpenExpr, exprOpenFlagUsage)
	exprClose := flag.String("expr-close", lang.DefaultCloseExpr, exprCloseFlagUsage)
	return exprOpen, exprClose
}

func parseFlags() {
	flag.CommandLine.SetOutput(io.Discard)
	flag.Parse()
}

func printVersionAndExit() {
	fmt.Printf("Mimic version %s\n", Version)
	os.Exit(0)
}

func sourceAndtargetPath() (string, string) {
	args := flag.Args()

	sourcePath := "./.mimic"
	targetPath := "."

	if len(args) >= 1 {
		sourcePath = args[0]
	}

	if len(args) >= 2 {
		targetPath = args[1]
	}

	if len(args) > 2 {
		flag.Usage()
		os.Exit(1)
	}

	return sourcePath, targetPath
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

	sourcePath, targetPath := sourceAndtargetPath()

	env := lang.NewEnvironment()
	maps.Copy(env.Vars, vars)

	comp := lang.NewCompilerConfigurable(env, lang.NewExpressionConfigurable(*exprOpen, *exprClose))

	executor := core.NewExecutor(sourcePath, targetPath, comp)

	executor.Read()

	for fileName := range executor.FilesRead {
		cli.LogFileNameAt(fileName)
	}

	if !cli.MustConfirmToContinue() {
		os.Exit(0)
	}

	executor.Write()

	for fileName, fileData := range executor.WrittenFiles {
		cli.LogFileNameAt(fileName)
		cli.LogFileDataAdded(fileData)
	}

	os.Exit(0)
}
