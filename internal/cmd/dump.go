package cmd

import (
	"flag"
	"fmt"
	"io"
	"maps"
	"mimic/internal"
	"mimic/internal/lang"
	"mimic/internal/util"
	"os"
)

const (
	DumpCommandDescription = "Dump information about variables declared in files and directories and exit"
)

const (
	DumpCommandSourceFlagUsage = "Set the source path of all files and directories to dump (default .mimic)"

	DumpCommandVarValueFlagUsage  = "Set a variable value by passing a key=value pair"
	DumpCommandVarPromptFlagUsage = "Set a variable prompt message by passing a key=value pair"

	DumpCommandExprOpenFlagUsage  = "Set the open expression syntax (default \"{{\")"
	DumpCommandExprCloseFlagUsage = "Set the close expression syntax (default \"}}\")"

	DumpCommandDebugModeFlagUsage = "Enable debug mode (default false)"
)

func DumpCommandUsage() {
	fmt.Fprintf(os.Stderr, "Usage: mimic dump [OPTION]...\n")
	fmt.Fprintf(os.Stderr, "%s\n", DumpCommandDescription)
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	fmt.Fprintf(os.Stderr, "  -s, --source    %s\n", DumpCommandSourceFlagUsage)
	fmt.Fprintf(os.Stderr, "  -v, --var       %s\n", DumpCommandVarValueFlagUsage)
	fmt.Fprintf(os.Stderr, "  -p, --prompt    %s\n", DumpCommandVarPromptFlagUsage)
	fmt.Fprintf(os.Stderr, "  --expr-open     %s\n", DumpCommandExprOpenFlagUsage)
	fmt.Fprintf(os.Stderr, "  --expr-close    %s\n", DumpCommandExprCloseFlagUsage)
	fmt.Fprintf(os.Stderr, "  --debug     	  %s\n", DumpCommandDebugModeFlagUsage)
	fmt.Fprintln(os.Stderr)
}

type DumpCommandConfig struct {
	SourcePath util.FlagSlice
	TargetPath util.FlagSlice

	Variables util.FlagMap
	Prompts   util.FlagMap

	ExprOpen  string
	ExprClose string

	DebugMode bool
}

type DumpCommand struct {
	FlagSet *flag.FlagSet

	config *DumpCommandConfig
}

func NewDumpCommandConfig() *DumpCommandConfig {
	return &DumpCommandConfig{
		SourcePath: util.NewFlagSlice(".mimic"),
		TargetPath: util.NewFlagSlice("."),

		Variables: make(util.FlagMap),
		Prompts:   make(util.FlagMap),
	}
}

func NewDumpCommand() *DumpCommand {
	config := NewDumpCommandConfig()

	flagSet := flag.NewFlagSet("dump", flag.ExitOnError)
	flagSet.SetOutput(io.Discard)

	flagSet.Usage = DumpCommandUsage

	flagSet.Var(&config.SourcePath, "s", DumpCommandSourceFlagUsage)
	flagSet.Var(&config.SourcePath, "source", DumpCommandSourceFlagUsage)

	flagSet.Var(&config.Variables, "v", DumpCommandVarValueFlagUsage)
	flagSet.Var(&config.Variables, "var", DumpCommandVarValueFlagUsage)

	flagSet.Var(&config.Prompts, "p", DumpCommandVarPromptFlagUsage)
	flagSet.Var(&config.Prompts, "prompt", DumpCommandVarPromptFlagUsage)

	flagSet.StringVar(&config.ExprOpen, "expr-open", lang.DefaultOpenExpr, DumpCommandExprOpenFlagUsage)
	flagSet.StringVar(&config.ExprClose, "expr-close", lang.DefaultCloseExpr, DumpCommandExprCloseFlagUsage)

	flagSet.BoolVar(&config.DebugMode, "debug", false, DumpCommandDebugModeFlagUsage)

	return &DumpCommand{
		FlagSet: flagSet,

		config: config,
	}
}

func (c *DumpCommand) Run(args []string) {
	c.FlagSet.Parse(args)

	scanner := internal.NewScanner(c.config.DebugMode)
	filesRead := scanner.Scan(c.config.SourcePath.Values)

	env := lang.NewEnvironment()
	maps.Copy(env.Vars, c.config.Variables)
	maps.Copy(env.Prompts, c.config.Prompts)

	expr := lang.NewExpressionConfigurable(c.config.ExprOpen, c.config.ExprClose)

	analyzer := lang.NewAnalyzerConfigurable(env, expr)

	dumper := internal.NewDumper(analyzer, c.config.DebugMode)
	dumper.Dump(filesRead)
}
