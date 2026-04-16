package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"mimic/cli"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const VariableMatchRegularExpression = `\[\s*(.*?)\s*\]`

const SourceFlagUsage = "Set the source directory path of .mimic files"
const TargetFlagUsage = "Set the target path where all files will be copied"

type Mimic struct {
	source     string
	sourceFile os.FileInfo
	target     string
	targetFile os.FileInfo
	files      map[string]string
	vars       map[string]string
}

func NewMimic() *Mimic {
	var source string
	var target string

	flag.StringVar(&source, "s", "./.mimic", SourceFlagUsage)
	flag.StringVar(&source, "source", "./.mimic", SourceFlagUsage)

	flag.StringVar(&target, "t", ".", TargetFlagUsage)
	flag.StringVar(&target, "target", ".", TargetFlagUsage)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: mimic [-s | --source] [-t | --target]\n")
		fmt.Fprintf(os.Stderr, "Mimic interpret .mimic files in the source path (./.mimic directory by default) and create copies of them in the target path (the current directory by default).\n\n")
		fmt.Fprintf(os.Stderr, "Configure how to start mimicking files across your entire filesystem\n")
		fmt.Fprintf(os.Stderr, "  -s, --source    %s\n", SourceFlagUsage)
		fmt.Fprintf(os.Stderr, "  -t, --target    %s\n\n", TargetFlagUsage)
	}

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

	return &Mimic{
		source:     source,
		sourceFile: sourceFile,
		target:     target,
		targetFile: targetFile,
		files:      make(map[string]string),
		vars:       make(map[string]string),
	}
}

func (m *Mimic) Scan() {
	cli.Log(fmt.Sprintf("Scanning files from the source directory %s...", m.source), cli.LogSeverityWarn)

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

	re := regexp.MustCompile(VariableMatchRegularExpression)

	matches := re.FindAllStringSubmatch(name+string(data), -1)

	for _, match := range matches {
		if _, exists := m.vars[match[0]]; !exists {
			m.vars[match[0]] = cli.MustAsk(fmt.Sprintf("Please enter a value for %s: ", match[0]))
		}
	}

	m.files[name] = string(data)
}

func (m *Mimic) Copy() {
	cli.Log(fmt.Sprintf("Copying files to the target directory %s...", m.target), cli.LogSeverityWarn)

	for name, data := range m.files {
		name = name[:len(name)-len(".mimic")]
		name = name[len(m.source)-1:]
		name = m.target + "/" + name

		name = m.fill(name)
		data = m.fill(data)

		m.write(name, data)

		cli.LogFileNameAt(name)
		cli.LogFileDataAdded(data)
	}
}

func (m *Mimic) fill(text string) string {
	result := text

	for name, value := range m.vars {
		result = strings.ReplaceAll(result, name, value)
	}

	return result
}

func (m *Mimic) write(name string, data string) {
	dirname := filepath.Dir(name)

	if err := os.MkdirAll(dirname, 0755); err != nil {
		cli.Log("Could not create %s", cli.LogSeverityWarn)
	}

	if err := os.WriteFile(name, []byte(data), 0755); err != nil {
		cli.Log("Could not create %s", cli.LogSeverityWarn)
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
