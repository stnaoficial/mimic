package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"mimic/cli"
	"mimic/util"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var VariableModifiers = []string{
	"camel",
	"pascal",
	"snake",
	"kebab",
	"dot",
	"flat",
	"lower",
	"upper",
}

type Mimic struct {
	source     string
	sourceFile os.FileInfo
	target     string
	targetFile os.FileInfo
	fileMap    map[string]string
	varMap     map[string]string
	varRegex   *regexp.Regexp
}

func NewMimic(source string, target string, varMap map[string]string, varRegex *regexp.Regexp) *Mimic {
	sourceFile, err := os.Stat(source)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to obtain information about path %s", source), cli.LogSeverityError)
	}

	targetFile, err := os.Stat(target)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to obtain information about path %s", target), cli.LogSeverityError)
	}

	return &Mimic{
		source:     source,
		sourceFile: sourceFile,
		target:     target,
		targetFile: targetFile,
		fileMap:    make(map[string]string),
		varMap:     varMap,
		varRegex:   varRegex,
	}
}

func (m *Mimic) Scan() {
	cli.Log(fmt.Sprintf("Scanning files from the source path %s...", m.source), cli.LogSeverityInfo)

	names, err := m.walk(m.source)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to walk into the source path %s", m.source), cli.LogSeverityError)
	}

	if len(names) == 0 {
		cli.LogAndExit(fmt.Sprintf("No .mimic files found in the source path %s", m.source), cli.LogSeverityWarn)
	}

	for _, name := range names {
		cli.LogFileNameAt(name)
	}

	for _, name := range names {
		m.collect(name)
	}
}

func (m *Mimic) walk(root string) ([]string, error) {
	names := []string{}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".mimic") {
			names = append(names, path)
		}

		return nil
	})

	return names, err
}

func (m *Mimic) collect(name string) {
	data, err := os.ReadFile(name)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to obtain data from file %s", name), cli.LogSeverityError)
	}

	submatches := m.varRegex.FindAllStringSubmatch(name+string(data), -1)

	for _, submatch := range submatches {
		_, name := m.parse(submatch[1])

		if _, exists := m.varMap[name]; !exists {
			m.varMap[name] = cli.MustAsk(fmt.Sprintf("Please enter a value for \"%s\": ", name))
		}
	}

	m.fileMap[name] = string(data)
}

func (m *Mimic) parse(value string) ([]string, string) {
	current := strings.TrimSpace(value)

	var modifiers []string

	for {
		beginParen := strings.Index(current, "(")
		endParen := strings.LastIndex(current, ")")

		// If no more balanced parentheses are found, we are at the core key
		if beginParen == -1 || endParen == -1 || endParen < beginParen {
			break
		}

		// Extracts the modifier name
		// e.g. "modifier(name)" -> extracts "modifier"
		modifier := strings.TrimSpace(current[:beginParen])

		if slices.Contains(VariableModifiers, modifier) {
			modifiers = append(modifiers, modifier)

			// Updates 'current' to be the content INSIDE the parentheses
			// current becomes "name"
			current = strings.TrimSpace(current[beginParen+1 : endParen])
		} else {
			// If it has parens but the prefix isn't a modifier,
			// stop and treat the remaining string as the key
			break
		}
	}

	name := current

	if name == "" {
		cli.LogAndExit(fmt.Sprintf("Unable to parse variable %s", value), cli.LogSeverityError)
	}

	// Reverse so innermost is first
	// e.g. {{ modifier1(modifier0(name)) }} -> [modifier0, modifier1]
	slices.Reverse(modifiers)

	return modifiers, name
}

func (m *Mimic) modify(modifiers []string, value string) string {
	result := value

	for _, mod := range modifiers {
		switch mod {
		case "camel":
			result = util.ToCamel(result)
		case "pascal":
			result = util.ToPascal(result)
		case "snake":
			result = util.ToSnake(result)
		case "kebab":
			result = util.ToKebab(result)
		case "dot":
			result = util.ToDot(result)
		case "flat":
			result = util.ToFlat(result)
		case "lower":
			result = util.ToLower(result)
		case "upper":
			result = util.ToUpper(result)
		default:
			continue
		}
	}

	return result
}

func (m *Mimic) Start() {
	cli.Log(fmt.Sprintf("Starting to mimic files to the target path %s...", m.target), cli.LogSeverityInfo)

	for name, data := range m.fileMap {
		rel, err := filepath.Rel(m.source, name)

		if err != nil {
			cli.LogAndExit(fmt.Sprintf("Unable to relate path %s to %s", m.source, name), cli.LogSeverityError)
		}

		name = filepath.Join(m.target, rel[:len(rel)-len(".mimic")])

		name = m.fill(name)
		data = m.fill(data)

		m.write(name, data)

		cli.LogFileNameAt(name)
		cli.LogFileDataAdded(data)
	}
}

func (m *Mimic) fill(text string) string {
	return m.varRegex.ReplaceAllStringFunc(text, func(match string) string {
		submatch := m.varRegex.FindStringSubmatch(match)

		modifiers, name := m.parse(submatch[1])

		if value, exists := m.varMap[name]; exists {
			return m.modify(modifiers, value)
		}

		return match
	})
}

func (m *Mimic) write(name string, data string) {
	dirname := filepath.Dir(name)

	if err := os.MkdirAll(dirname, 0755); err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to create file %s", name), cli.LogSeverityError)
	}

	if err := os.WriteFile(name, []byte(data), 0644); err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to create file %s", name), cli.LogSeverityError)
	}
}

const VarFlagUsage = "Set a variable directly by passing a key=value pair"
const VarPrefixFlagUsage = "Set the variable pattern prefix (default \"{{\")"
const VarSufixFlagUsage = "Set the variable pattern sufix (default \"}}\")"

const HelpFlagUsage = "Print Help (this message) and exit"
const VersionFlagUsage = "Print version information and exit"

func PrintVersionAndExit(_ string) {
	fmt.Printf("Mimic version 1.0.0-beta\n")
	os.Exit(0)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: mimic [OPTION]... SOURCE TARGET\n")
		fmt.Fprintf(os.Stderr, "Mimic interpret .mimic files in the source path (./.mimic directory by default) and create copies of them in the target path (the current directory by default).\n\n")
		fmt.Fprintf(os.Stderr, "Provide variables directly\n")
		fmt.Fprintf(os.Stderr, "  -v, --var       %s\n\n", VarFlagUsage)
		fmt.Fprintf(os.Stderr, "Configure how to start mimicking files in the source path\n")
		fmt.Fprintf(os.Stderr, "  --var-prefix    %s\n", VarPrefixFlagUsage)
		fmt.Fprintf(os.Stderr, "  --var-sufix     %s\n\n", VarSufixFlagUsage)
		fmt.Fprintf(os.Stderr, "Get more information\n")
		fmt.Fprintf(os.Stderr, "  -h, --help    %s\n", HelpFlagUsage)
		fmt.Fprintf(os.Stderr, "  --version     %s\n\n", VersionFlagUsage)
	}

	source := "./.mimic"
	target := "."
	vars := make(util.FlagMap)
	flag.Var(&vars, "v", VarFlagUsage)
	flag.Var(&vars, "var", VarFlagUsage)

	varPrefix := flag.String("var-prefix", "{{", VarPrefixFlagUsage)
	varSufix := flag.String("var-sufix", "}}", VarSufixFlagUsage)
	varRegex := regexp.MustCompile(regexp.QuoteMeta(*varPrefix) + `\s*(.*?)\s*` + regexp.QuoteMeta(*varSufix))

	version := flag.Bool("version", false, VersionFlagUsage)

	flag.CommandLine.SetOutput(io.Discard)
	flag.Parse()

	if *version {
		fmt.Printf("Mimic version 1.0.0-beta\n")
		os.Exit(0)
	}

	args := flag.Args()

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

	mimic := NewMimic(source, target, vars, varRegex)

	mimic.Scan()

	if !cli.MustConfirmToContinue() {
		os.Exit(0)
	}

	mimic.Start()

	os.Exit(0)
}
