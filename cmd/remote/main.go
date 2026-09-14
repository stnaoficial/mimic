package remote

import (
	"flag"
	"fmt"
	"io"
	"maps"
	"mimic/cmd"
	"mimic/internal"
	"mimic/internal/cli"
	"mimic/internal/lang"
	"mimic/internal/util"
	"mimic/tools/github"
	"os"
	"path/filepath"
	"strings"
)

var GithubApiClient *github.ApiClient

const (
	CommandDescription = "Start generating templates from remote repositories"
)

const (
	CommandTemplateNameFlagUsage = "Set the template name to copy"

	CommandSourceFlagUsage = "Set the source path for templates to copy (default .mimic/templates)"
	CommandTargetFlagUsage = "Set the target path for generated templates (default .)"

	CommandVarValueFlagUsage  = "Set a variable value by passing a key=value pair"
	CommandVarPromptFlagUsage = "Set a variable prompt message by passing a key=value pair"

	CommandExprOpenFlagUsage  = "Set the open expression syntax (default \"{{\")"
	CommandExprCloseFlagUsage = "Set the close expression syntax (default \"}}\")"

	CommandListFlagUsage = "List all templates available"
	CommandDumpFlagUsage = "Dump all variables found in the current operation"

	CommandNoAskFlagUsage      = "Disable the \"ask to confirm\" safety feature (default false)"
	CommandNoCacheFlagUsage    = "Disable the caching feature (default false)"
	CommandDebugModeFlagUsage  = "Enable debug mode (default false)"
	CommandStrictModeFlagUsage = "Enable strict mode (default false)"
	CommandWriteModeFlagUsage  = "Set the write mode 0=override, 1=append) (default 0)"
)

func CommandUsage() {
	fmt.Fprintf(os.Stderr, "Usage: mimic remote [OPTION]...\n")
	fmt.Fprintf(os.Stderr, "%s\n", CommandDescription)
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	fmt.Fprintf(os.Stderr, "  -n, --name      %s\n", CommandTemplateNameFlagUsage)
	fmt.Fprintf(os.Stderr, "  -s, --source    %s\n", CommandSourceFlagUsage)
	fmt.Fprintf(os.Stderr, "  -t, --target    %s\n", CommandTargetFlagUsage)
	fmt.Fprintf(os.Stderr, "  -v, --var       %s\n", CommandVarValueFlagUsage)
	fmt.Fprintf(os.Stderr, "  -p, --prompt    %s\n", CommandVarPromptFlagUsage)
	fmt.Fprintf(os.Stderr, "  --expr-open     %s\n", CommandExprOpenFlagUsage)
	fmt.Fprintf(os.Stderr, "  --expr-close    %s\n", CommandExprCloseFlagUsage)
	fmt.Fprintf(os.Stderr, "  --list     	  %s\n", CommandListFlagUsage)
	fmt.Fprintf(os.Stderr, "  --dump     	  %s\n", CommandDumpFlagUsage)
	fmt.Fprintf(os.Stderr, "  --no-ask     	  %s\n", CommandNoAskFlagUsage)
	fmt.Fprintf(os.Stderr, "  --no-cache      %s\n", CommandNoCacheFlagUsage)
	fmt.Fprintf(os.Stderr, "  --debug     	  %s\n", CommandDebugModeFlagUsage)
	fmt.Fprintf(os.Stderr, "  --strict     	  %s\n", CommandStrictModeFlagUsage)
	fmt.Fprintf(os.Stderr, "  --mode     	  %s\n", CommandWriteModeFlagUsage)
	fmt.Fprintln(os.Stderr)
}

type CommandConfig struct {
	TemplateName util.FlagSlice

	SourcePath util.FlagSlice
	TargetPath util.FlagSlice

	Variables util.FlagMap
	Prompts   util.FlagMap

	ExprOpen  string
	ExprClose string

	List bool
	Dump bool

	NoAsk      bool
	NoCache    bool
	DebugMode  bool
	StrictMode bool
	WriteMode  int
}

type Command struct {
	name string

	FlagSet *flag.FlagSet

	config *CommandConfig
}

func NewCommandConfig() *CommandConfig {
	return &CommandConfig{
		SourcePath: util.NewFlagSlice(),
		TargetPath: util.NewFlagSlice("."),

		Variables: make(util.FlagMap),
		Prompts:   make(util.FlagMap),
	}
}

func NewCommand(name string) *Command {
	config := NewCommandConfig()

	flagSet := flag.NewFlagSet(name, flag.ExitOnError)
	flagSet.SetOutput(io.Discard)

	flagSet.Usage = CommandUsage

	flagSet.Var(&config.TemplateName, "n", CommandTemplateNameFlagUsage)
	flagSet.Var(&config.TemplateName, "name", CommandTemplateNameFlagUsage)

	flagSet.Var(&config.SourcePath, "s", CommandSourceFlagUsage)
	flagSet.Var(&config.SourcePath, "source", CommandSourceFlagUsage)

	flagSet.Var(&config.TargetPath, "t", CommandTargetFlagUsage)
	flagSet.Var(&config.TargetPath, "target", CommandTargetFlagUsage)

	flagSet.Var(&config.Variables, "v", CommandVarValueFlagUsage)
	flagSet.Var(&config.Variables, "var", CommandVarValueFlagUsage)

	flagSet.Var(&config.Prompts, "p", CommandVarPromptFlagUsage)
	flagSet.Var(&config.Prompts, "prompt", CommandVarPromptFlagUsage)

	flagSet.StringVar(&config.ExprOpen, "expr-open", lang.DefaultOpenExpr, CommandExprOpenFlagUsage)
	flagSet.StringVar(&config.ExprClose, "expr-close", lang.DefaultCloseExpr, CommandExprCloseFlagUsage)

	flagSet.BoolVar(&config.List, "list", false, CommandListFlagUsage)
	flagSet.BoolVar(&config.Dump, "dump", false, CommandDumpFlagUsage)

	flagSet.BoolVar(&config.NoAsk, "no-ask", false, CommandNoAskFlagUsage)
	flagSet.BoolVar(&config.NoCache, "no-cache", false, CommandNoCacheFlagUsage)
	flagSet.BoolVar(&config.DebugMode, "debug", false, CommandDebugModeFlagUsage)
	flagSet.BoolVar(&config.StrictMode, "strict", false, CommandStrictModeFlagUsage)
	flagSet.IntVar(&config.WriteMode, "mode", 0, CommandWriteModeFlagUsage)

	return &Command{
		name: name,

		FlagSet: flagSet,

		config: config,
	}
}

func (c *Command) Name() string { return c.name }

func (c *Command) Parse(args []string) { c.FlagSet.Parse(args) }

func (c *Command) Setup() {
	var localConfig, err = cmd.NewLocalConfig()

	if err != nil {
		// allow debug
		if c.config.DebugMode {
			cli.Logln(cli.LogSeverityError, err.Error())
		}

		fmt.Printf("Unable to start operation\n\n")

		os.Exit(1)
	}

	GithubApiClient = github.NewApiClient(
		localConfig.Get("remote.github.user.username"),
		localConfig.Get("remote.github.repository.name"),
		localConfig.Get("remote.github.repository.branch.name"),
		localConfig.Get("remote.github.user.access-token"),
	)

	if c.config.NoCache {
		GithubApiClient.Cache = false
	}
}

func (c *Command) Validate() {
	templateNames := c.config.TemplateName.Values
	sourcePaths := c.config.SourcePath.Values

	if len(templateNames) > 0 && len(sourcePaths) > 0 {
		fmt.Printf("You cannot use both a template name and a source path at the same time\n\n")
		os.Exit(1)
	}
}

func (c *Command) Run() {
	templateNames := c.config.TemplateName.Values
	sourcePaths := c.config.SourcePath.Values
	targetPaths := c.config.TargetPath.Values

	if len(templateNames) > 0 {
		for _, templateName := range templateNames {
			sourcePaths = append(sourcePaths, filepath.Join(cmd.DefaultConfigTemplatesDirectoryPath, templateName))
		}
	}

	if len(sourcePaths) == 0 {
		sourcePaths = append(sourcePaths, cmd.DefaultConfigTemplatesDirectoryPath)
	}

	reader := NewReader(c.config.DebugMode)

	if c.config.List {
		if err := reader.List(sourcePaths); err != nil {
			// allow debug
			if c.config.DebugMode {
				cli.Logln(cli.LogSeverityError, err.Error())
			}

			fmt.Printf("No templates available\n\n")

			os.Exit(1)
		}

		os.Exit(0)
	}

	lastSourcePath := cmd.DefaultConfigTemplatesDirectoryPath

	if value, err := c.config.SourcePath.Last(); err == nil {
		lastSourcePath = value
	}

	cli.Printf(cli.Normal, cli.Cyan, "Connected to GitHub on branch %s of %s/%s\n\n",
		GithubApiClient.BranchName(),
		GithubApiClient.Username(),
		GithubApiClient.RepositoryName(),
	)

	if len(templateNames) == 0 && strings.HasSuffix(lastSourcePath, cmd.DefaultConfigTemplatesDirectoryPath) {
		fmt.Printf("No template selected (use --list to list all templates available then -n or --name to select one)\n")
		cli.Printf(cli.Normal, cli.Yellow, "Considering all templates defined in %s\n", lastSourcePath)
		cli.Printf(cli.Normal, cli.Yellow, "This action may produce unexpected behavior\n\n")

		fmt.Printf("You can change some settings in %s\n\n", cmd.DefaultConfigFilePath)

		if !c.config.NoAsk && !cli.Confirm("Do you want to continue [Y/n]? ") {
			os.Exit(0)
		}
	}

	scanner := NewScanner(c.config.DebugMode)

	scannedEntries, err := scanner.Scan(sourcePaths)

	if err != nil {
		// allow debug
		if c.config.DebugMode {
			cli.Logln(cli.LogSeverityError, err.Error())
		}

		fmt.Printf("No such file or directory\n\n")

		os.Exit(1)
	}

	env := lang.NewEnvironment()
	maps.Copy(env.Vars, c.config.Variables)
	maps.Copy(env.Prompts, c.config.Prompts)

	expr := lang.NewExpressionConfigurable(c.config.ExprOpen, c.config.ExprClose)

	if c.config.Dump {
		analyzer := lang.NewAnalyzerConfigurable(env, expr)

		dumper := internal.NewDumper(analyzer, c.config.DebugMode)

		if _, err := dumper.Dump(scannedEntries); err != nil {
			// allow debug
			if c.config.DebugMode {
				cli.Logln(cli.LogSeverityError, err.Error())
			}

			fmt.Printf("Unable to dump variables\n\n")

			os.Exit(1)
		}

		os.Exit(0)
	}

	comp := lang.NewCompilerConfigurable(env, expr, c.config.StrictMode)

	generator := internal.NewGenerator(comp, c.config.DebugMode)
	generatedEntries, err := generator.Generate(targetPaths, scannedEntries)

	if err != nil {
		// allow debug
		if c.config.DebugMode {
			cli.Logln(cli.LogSeverityError, err.Error())
		}

		fmt.Printf("Unable to start compiling\n\n")

		os.Exit(1)
	}

	if len(generatedEntries) > 0 && !c.config.NoAsk && !cli.Confirm("Do you want to continue [Y/n]? ") {
		os.Exit(0)
	}

	writer := internal.NewWriter(c.config.DebugMode, internal.WriteMode(c.config.WriteMode))
	_, err = writer.Write(generatedEntries)

	if err != nil {
		// allow debug
		if c.config.DebugMode {
			cli.Logln(cli.LogSeverityError, err.Error())
		}

		fmt.Printf("Unable to write files\n\n")

		os.Exit(1)
	}
}
