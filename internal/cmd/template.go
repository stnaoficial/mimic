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
	TemplateCommandDescription = "Start generating copies of templates"
)

const (
	TemplateCommandSourceNameFlagUsage = "Set the source name of all templates to copy (required)"
	TemplateCommandTargetFlagUsage     = "Set the target path for all copied files and directories (default .)"

	TemplateCommandVarValueFlagUsage  = "Set a variable value by passing a key=value pair"
	TemplateCommandVarPromptFlagUsage = "Set a variable prompt message by passing a key=value pair"

	TemplateCommandExprOpenFlagUsage  = "Set the open expression syntax (default \"{{\")"
	TemplateCommandExprCloseFlagUsage = "Set the close expression syntax (default \"}}\")"

	TemplateCommandUserFlagUsage = "Set the GitHub username of the template repository (default stnaoficial)"
	TemplateCommandRepoFlagUsage = "Set the GitHub repository name (default mimic-template)"
	TemplateCommandListFlagUsage = "List all available templates (default false)"

	TemplateCommandNoAskFlagUsage      = "Disable the \"ask to confirm\" safety feature (default false)"
	TemplateCommandDebugModeFlagUsage  = "Enable debug mode (default false)"
	TemplateCommandStrictModeFlagUsage = "Enable strict mode (default false)"
	TemplateCommandWriteModeFlagUsage  = "Set the write mode 0=override, 1=append) (default 0)"
)

func TemplateCommandUsage() {
	fmt.Fprintf(os.Stderr, "Usage: mimic template [OPTION]...\n")
	fmt.Fprintf(os.Stderr, "%s\n", TemplateCommandDescription)
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	fmt.Fprintf(os.Stderr, "  -n, --name      %s\n", TemplateCommandSourceNameFlagUsage)
	fmt.Fprintf(os.Stderr, "  -t, --target    %s\n", TemplateCommandTargetFlagUsage)
	fmt.Fprintf(os.Stderr, "  -v, --var       %s\n", TemplateCommandVarValueFlagUsage)
	fmt.Fprintf(os.Stderr, "  -p, --prompt    %s\n", TemplateCommandVarPromptFlagUsage)
	fmt.Fprintf(os.Stderr, "  --expr-open     %s\n", TemplateCommandExprOpenFlagUsage)
	fmt.Fprintf(os.Stderr, "  --expr-close    %s\n", TemplateCommandExprCloseFlagUsage)
	fmt.Fprintf(os.Stderr, "  --user     	  %s\n", TemplateCommandUserFlagUsage)
	fmt.Fprintf(os.Stderr, "  --repo     	  %s\n", TemplateCommandRepoFlagUsage)
	fmt.Fprintf(os.Stderr, "  --list     	  %s\n", TemplateCommandListFlagUsage)
	fmt.Fprintf(os.Stderr, "  --no-ask     	  %s\n", TemplateCommandNoAskFlagUsage)
	fmt.Fprintf(os.Stderr, "  --debug     	  %s\n", TemplateCommandDebugModeFlagUsage)
	fmt.Fprintf(os.Stderr, "  --strict     	  %s\n", TemplateCommandStrictModeFlagUsage)
	fmt.Fprintf(os.Stderr, "  --mode     	  %s\n", TemplateCommandWriteModeFlagUsage)
	fmt.Fprintln(os.Stderr)
}

type TemplateCommandConfig struct {
	SourceName util.FlagSlice
	TargetPath util.FlagSlice

	Variables util.FlagMap
	Prompts   util.FlagMap

	ExprOpen  string
	ExprClose string

	User string
	Repo string

	List bool

	NoAsk      bool
	DebugMode  bool
	StrictMode bool
	WriteMode  int
}

type TemplateCommand struct {
	FlagSet *flag.FlagSet

	config *TemplateCommandConfig
}

func NewTemplateCommandConfig() *TemplateCommandConfig {
	return &TemplateCommandConfig{
		TargetPath: util.NewFlagSlice("."),

		Variables: make(util.FlagMap),
		Prompts:   make(util.FlagMap),

		User: "stnaoficial",
		Repo: "mimic-template",
	}
}

func NewTemplateCommand() *TemplateCommand {
	config := NewTemplateCommandConfig()

	flagSet := flag.NewFlagSet("copy", flag.ExitOnError)
	flagSet.SetOutput(io.Discard)

	flagSet.Usage = TemplateCommandUsage

	flagSet.Var(&config.SourceName, "n", TemplateCommandSourceNameFlagUsage)
	flagSet.Var(&config.SourceName, "name", TemplateCommandSourceNameFlagUsage)

	flagSet.Var(&config.TargetPath, "t", TemplateCommandTargetFlagUsage)
	flagSet.Var(&config.TargetPath, "target", TemplateCommandTargetFlagUsage)

	flagSet.Var(&config.Variables, "v", TemplateCommandVarValueFlagUsage)
	flagSet.Var(&config.Variables, "var", TemplateCommandVarValueFlagUsage)

	flagSet.Var(&config.Prompts, "p", TemplateCommandVarPromptFlagUsage)
	flagSet.Var(&config.Prompts, "prompt", TemplateCommandVarPromptFlagUsage)

	flagSet.StringVar(&config.ExprOpen, "expr-open", lang.DefaultOpenExpr, TemplateCommandExprOpenFlagUsage)
	flagSet.StringVar(&config.ExprClose, "expr-close", lang.DefaultCloseExpr, TemplateCommandExprCloseFlagUsage)

	flagSet.StringVar(&config.User, "user", "stnaoficial", TemplateCommandUserFlagUsage)
	flagSet.StringVar(&config.Repo, "repo", "mimic-template", TemplateCommandRepoFlagUsage)

	flagSet.BoolVar(&config.List, "list", false, TemplateCommandListFlagUsage)

	flagSet.BoolVar(&config.NoAsk, "no-ask", false, TemplateCommandNoAskFlagUsage)
	flagSet.BoolVar(&config.DebugMode, "debug", false, TemplateCommandDebugModeFlagUsage)
	flagSet.BoolVar(&config.StrictMode, "strict", false, TemplateCommandStrictModeFlagUsage)
	flagSet.IntVar(&config.WriteMode, "mode", 0, TemplateCommandWriteModeFlagUsage)

	return &TemplateCommand{
		FlagSet: flagSet,

		config: config,
	}
}

func (c *TemplateCommand) Run(args []string) {
	c.FlagSet.Parse(args)

	githubScanner := internal.NewGitHubScanner(c.config.User, c.config.Repo, c.config.DebugMode)

	if c.config.List {
		githubScanner.List()
		os.Exit(0)
	}

	if len(c.config.SourceName.Values) == 0 {
		c.FlagSet.Usage()
		os.Exit(1)
	}

	filesRead := githubScanner.Scan(c.config.SourceName.Values)

	env := lang.NewEnvironment()
	maps.Copy(env.Vars, c.config.Variables)
	maps.Copy(env.Prompts, c.config.Prompts)

	expr := lang.NewExpressionConfigurable(c.config.ExprOpen, c.config.ExprClose)
	comp := lang.NewCompilerConfigurable(env, expr, c.config.StrictMode)

	generator := internal.NewGenerator(comp, c.config.DebugMode)
	filesGenerated := generator.Generate(c.config.TargetPath.Values, filesRead)

	if !c.config.NoAsk && !cli.Confirm("Do you want to continue [Y/n]? ") {
		os.Exit(0)
	}

	writer := internal.NewWriter(c.config.DebugMode, internal.WriteMode(c.config.WriteMode))
	writer.Write(filesGenerated)
}
