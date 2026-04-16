package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"mimic/cli"
	"mimic/utils"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const VariablePrefix = "{{"
const VariableSufix = "}}"

var variableRegex = regexp.MustCompile(regexp.QuoteMeta(VariablePrefix) + `\s*(.*?)\s*` + regexp.QuoteMeta(VariableSufix))

const SourceFlagUsage = "Set the source directory path of .mimic files"
const TargetFlagUsage = "Set the target path where all files will be copied"
const VarFlagUsage = "Set a var directly by passing as a key=value pair"

type Mimic struct {
	source     string
	sourceFile os.FileInfo
	target     string
	targetFile os.FileInfo
	fileMap    map[string]string
	varMap     map[string]string
}

func NewMimic() *Mimic {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: mimic [-s | --source] [-t | --target] [-v | --var]\n")
		fmt.Fprintf(os.Stderr, "Mimic interpret .mimic files in the source path (./.mimic directory by default) and create copies of them in the target path (the current directory by default).\n\n")
		fmt.Fprintf(os.Stderr, "Configure how to start mimicking files across your entire filesystem\n")
		fmt.Fprintf(os.Stderr, "  -s, --source    %s\n", SourceFlagUsage)
		fmt.Fprintf(os.Stderr, "  -t, --target    %s\n", TargetFlagUsage)
		fmt.Fprintf(os.Stderr, "  -v, --var       %s\n\n", VarFlagUsage)
	}

	var source string
	var target string

	flag.StringVar(&source, "s", "./.mimic", SourceFlagUsage)
	flag.StringVar(&source, "source", "./.mimic", SourceFlagUsage)

	flag.StringVar(&target, "t", ".", TargetFlagUsage)
	flag.StringVar(&target, "target", ".", TargetFlagUsage)

	vars := make(utils.FlagMap)

	flag.Var(&vars, "v", VarFlagUsage)
	flag.Var(&vars, "var", VarFlagUsage)

	flag.CommandLine.SetOutput(io.Discard)

	flag.Parse()

	if flag.NArg() > 0 {
		flag.Usage()
		os.Exit(1)
	}

	sourceFile, err := os.Stat(source)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to get information about %s", source), cli.LogSeverityError)
	}

	targetFile, err := os.Stat(target)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to get information about %s", target), cli.LogSeverityError)
	}

	varMap := make(map[string]string)

	for key, value := range vars {
		varMap[fmt.Sprintf("%s%s%s", VariablePrefix, key, VariableSufix)] = value
	}

	return &Mimic{
		source:     source,
		sourceFile: sourceFile,
		target:     target,
		targetFile: targetFile,
		fileMap:    make(map[string]string),
		varMap:     varMap,
	}
}

func (m *Mimic) Scan() {
	cli.Log(fmt.Sprintf("Scanning files from the source directory %s...", m.source), cli.LogSeverityInfo)

	names, err := m.walk(m.source)

	if err != nil {
		cli.LogAndExit("Unable to walk into the given source directory path", cli.LogSeverityError)
	}

	if len(names) == 0 {
		cli.LogAndExit("No .mimic files found in the specified source directory", cli.LogSeverityWarn)
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
		cli.LogAndExit(fmt.Sprintf("Unable to read %s", name), cli.LogSeverityError)
	}

	matches := variableRegex.FindAllStringSubmatch(name+string(data), -1)

	for _, match := range matches {
		if _, exists := m.varMap[match[0]]; !exists {
			m.varMap[match[0]] = cli.MustAsk(fmt.Sprintf("Please enter a value for %s: ", match[0]))
		}
	}

	m.fileMap[name] = string(data)
}

func (m *Mimic) Copy() {
	cli.Log(fmt.Sprintf("Copying files to the target directory %s...", m.target), cli.LogSeverityInfo)

	for name, data := range m.fileMap {
		rel, err := filepath.Rel(m.source, name)

		if err != nil {
			cli.LogAndExit("Unable to resolve relative path", cli.LogSeverityError)
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
	return variableRegex.ReplaceAllStringFunc(text, func(match string) string {
		if value, exists := m.varMap[match]; exists {
			return value
		}

		return match
	})
}

func (m *Mimic) write(name string, data string) {
	dirname := filepath.Dir(name)

	if err := os.MkdirAll(dirname, 0755); err != nil {
		cli.LogAndExit(fmt.Sprintf("Could not create %s", name), cli.LogSeverityError)
	}

	if err := os.WriteFile(name, []byte(data), 0644); err != nil {
		cli.LogAndExit(fmt.Sprintf("Could not create %s", name), cli.LogSeverityError)
	}

	cli.LogFileNameAdded(name)
}

func main() {
	mimic := NewMimic()

	mimic.Scan()

	if !cli.ConfirmToContinue() {
		os.Exit(0)
	}

	mimic.Copy()

	os.Exit(0)
}
