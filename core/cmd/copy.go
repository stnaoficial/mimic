package cmd

import (
	"flag"
	"fmt"
	"mimic/core"
	"mimic/core/cli"
	"mimic/core/lang"
	"os"
)

const (
	CopyCommandDescription = "Start generating copies of files and directories"
)

const (
	CopyCommandVarFlagUsage       = "Set a variable value by passing a key=value pair"
	CopyCommandVarPromptFlagUsage = "Set a variable prompt message by passing a key=value pair"

	CopyCommandExprOpenFlagUsage  = "Set the open expression syntax (default \"{{\")"
	CopyCommandExprCloseFlagUsage = "Set the close expression syntax (default \"}}\")"

	CopyCommandDebugModeFlagUsage = "Enable debug mode (default false)"
)

func CopyCommandUsage() {
	fmt.Fprintf(os.Stderr, "Usage: mimic copy [OPTION]... SOURCE TARGET\n")
	fmt.Fprintf(os.Stderr, "%s\n", CopyCommandDescription)
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	fmt.Fprintf(os.Stderr, "  -v, --var       %s\n", CopyCommandVarFlagUsage)
	fmt.Fprintf(os.Stderr, "  -p, --prompt    %s\n", CopyCommandVarPromptFlagUsage)
	fmt.Fprintf(os.Stderr, "  --expr-open     %s\n", CopyCommandExprOpenFlagUsage)
	fmt.Fprintf(os.Stderr, "  --expr-close    %s\n", CopyCommandExprCloseFlagUsage)
	fmt.Fprintf(os.Stderr, "  --debug     	  %s\n", CopyCommandDebugModeFlagUsage)
	fmt.Fprintln(os.Stderr)
}

type CopyCommand struct {
	FlagSet *flag.FlagSet

	config *core.Config
}

func NewCopyCommand() *CopyCommand {
	config := core.NewConfig()
	flagSet := flag.NewFlagSet("copy", flag.ExitOnError)

	flagSet.SetOutput(os.Stderr)

	flagSet.Usage = CopyCommandUsage

	flagSet.Var(&config.Variables, "v", CopyCommandVarFlagUsage)
	flagSet.Var(&config.Variables, "var", CopyCommandVarFlagUsage)
	flagSet.Var(&config.Prompts, "p", CopyCommandVarPromptFlagUsage)
	flagSet.Var(&config.Prompts, "prompt", CopyCommandVarPromptFlagUsage)
	flagSet.StringVar(&config.ExprOpen, "expr-open", lang.DefaultOpenExpr, CopyCommandExprOpenFlagUsage)
	flagSet.StringVar(&config.ExprClose, "expr-close", lang.DefaultCloseExpr, CopyCommandExprCloseFlagUsage)
	flagSet.BoolVar(&config.DebugMode, "debug", false, CopyCommandDebugModeFlagUsage)

	return &CopyCommand{
		FlagSet: flagSet,

		config: config,
	}
}

func (c *CopyCommand) Run(args []string) {
	c.FlagSet.Parse(args)

	parsedArgs := c.FlagSet.Args()

	if len(parsedArgs) >= 1 {
		c.config.SourcePath = parsedArgs[0]
	}

	if len(parsedArgs) >= 2 {
		c.config.TargetPath = parsedArgs[1]
	}

	filesRead := core.NewScanner(c.config).Scan()

	filesGenerated := core.NewGenerator(c.config).Generate(filesRead)

	if !cli.MustConfirmToContinue() {
		os.Exit(0)
	}

	core.NewWriter(c.config).Write(filesGenerated)
}
