package cmd

import (
	"flag"
	"fmt"
	"io"
	"maps"
	"mimic/internal"
	"mimic/internal/cli"
	"mimic/internal/lang"
	"mimic/internal/util"
	"os"
)

const (
	CopyCommandDescription = "Start generating copies of files and directories"
)

const (
	CopyCommandSourceFlagUsage = "Set the source path of all files and directories to copy (default .mimic)"
	CopyCommandTargetFlagUsage = "Set the target path for all copied files and directories (default .)"

	CopyCommandVarValueFlagUsage  = "Set a variable value by passing a key=value pair"
	CopyCommandVarPromptFlagUsage = "Set a variable prompt message by passing a key=value pair"

	CopyCommandExprOpenFlagUsage  = "Set the open expression syntax (default \"{{\")"
	CopyCommandExprCloseFlagUsage = "Set the close expression syntax (default \"}}\")"

	CopyCommandDebugModeFlagUsage  = "Enable debug mode (default false)"
	CopyCommandStrictModeFlagUsage = "Enable strict mode (default false)"
)

func CopyCommandUsage() {
	fmt.Fprintf(os.Stderr, "Usage: mimic copy [OPTION]...\n")
	fmt.Fprintf(os.Stderr, "%s\n", CopyCommandDescription)
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	fmt.Fprintf(os.Stderr, "  -s, --source    %s\n", CopyCommandSourceFlagUsage)
	fmt.Fprintf(os.Stderr, "  -t, --target    %s\n", CopyCommandTargetFlagUsage)
	fmt.Fprintf(os.Stderr, "  -v, --var       %s\n", CopyCommandVarValueFlagUsage)
	fmt.Fprintf(os.Stderr, "  -p, --prompt    %s\n", CopyCommandVarPromptFlagUsage)
	fmt.Fprintf(os.Stderr, "  --expr-open     %s\n", CopyCommandExprOpenFlagUsage)
	fmt.Fprintf(os.Stderr, "  --expr-close    %s\n", CopyCommandExprCloseFlagUsage)
	fmt.Fprintf(os.Stderr, "  --debug     	  %s\n", CopyCommandDebugModeFlagUsage)
	fmt.Fprintf(os.Stderr, "  --strict     	  %s\n", CopyCommandStrictModeFlagUsage)
	fmt.Fprintln(os.Stderr)
}

type CopyCommandConfig struct {
	SourcePath util.FlagSlice
	TargetPath util.FlagSlice

	Variables util.FlagMap
	Prompts   util.FlagMap

	ExprOpen  string
	ExprClose string

	DebugMode  bool
	StrictMode bool
}

type CopyCommand struct {
	FlagSet *flag.FlagSet

	config *CopyCommandConfig
}

func NewCopyCommandConfig() *CopyCommandConfig {
	return &CopyCommandConfig{
		SourcePath: util.NewFlagSlice(".mimic"),
		TargetPath: util.NewFlagSlice("."),

		Variables: make(util.FlagMap),
		Prompts:   make(util.FlagMap),
	}
}

func NewCopyCommand() *CopyCommand {
	config := NewCopyCommandConfig()

	flagSet := flag.NewFlagSet("copy", flag.ExitOnError)
	flagSet.SetOutput(io.Discard)

	flagSet.Usage = CopyCommandUsage

	flagSet.Var(&config.SourcePath, "s", CopyCommandSourceFlagUsage)
	flagSet.Var(&config.SourcePath, "source", CopyCommandSourceFlagUsage)

	flagSet.Var(&config.TargetPath, "t", CopyCommandTargetFlagUsage)
	flagSet.Var(&config.TargetPath, "target", CopyCommandTargetFlagUsage)

	flagSet.Var(&config.Variables, "v", CopyCommandVarValueFlagUsage)
	flagSet.Var(&config.Variables, "var", CopyCommandVarValueFlagUsage)

	flagSet.Var(&config.Prompts, "p", CopyCommandVarPromptFlagUsage)
	flagSet.Var(&config.Prompts, "prompt", CopyCommandVarPromptFlagUsage)

	flagSet.StringVar(&config.ExprOpen, "expr-open", lang.DefaultOpenExpr, CopyCommandExprOpenFlagUsage)
	flagSet.StringVar(&config.ExprClose, "expr-close", lang.DefaultCloseExpr, CopyCommandExprCloseFlagUsage)

	flagSet.BoolVar(&config.DebugMode, "debug", false, CopyCommandDebugModeFlagUsage)
	flagSet.BoolVar(&config.StrictMode, "strict", false, CopyCommandStrictModeFlagUsage)

	return &CopyCommand{
		FlagSet: flagSet,

		config: config,
	}
}

func (c *CopyCommand) Run(args []string) {
	c.FlagSet.Parse(args)

	scanner := internal.NewScanner(c.config.DebugMode)
	filesRead := scanner.Scan(c.config.SourcePath.Values)

	env := lang.NewEnvironment()
	maps.Copy(env.Vars, c.config.Variables)
	maps.Copy(env.Prompts, c.config.Prompts)

	expr := lang.NewExpressionConfigurable(c.config.ExprOpen, c.config.ExprClose)
	comp := lang.NewCompilerConfigurable(env, expr, c.config.StrictMode)

	generator := internal.NewGenerator(comp, c.config.DebugMode)
	filesGenerated := generator.Generate(c.config.TargetPath.Values, filesRead)

	if !cli.Confirm("Do you want to continue [Y/n]? ") {
		os.Exit(0)
	}

	writer := internal.NewWriter(c.config.DebugMode)
	writer.Write(filesGenerated)
}
