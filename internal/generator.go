package internal

import (
	"mimic/internal/cli"
	"mimic/internal/lang"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Generator struct {
	comp *lang.Compiler

	generatedEntries FileSystemMap

	debug bool
}

func NewGenerator(comp *lang.Compiler, debug bool) *Generator {
	return &Generator{
		comp: comp,

		generatedEntries: make(FileSystemMap),

		debug: debug,
	}
}

func (g *Generator) defineGlobalVars() {
	// TODO
}

func (g *Generator) defineLocalVars(pathName string, _ FileSystemEntry) {
	dirName := filepath.Dir(pathName)

	count := 0

	if entries, err := os.ReadDir(dirName); err == nil {
		count = len(entries)
	}

	prevCount := count - 1

	if prevCount <= 0 {
		prevCount = 0
	}

	nextCount := count + 1

	g.comp.Env.Vars["__COUNT__"] = strconv.Itoa(count)
	g.comp.Env.Vars["__PREV_COUNT__"] = strconv.Itoa(prevCount)
	g.comp.Env.Vars["__NEXT_COUNT__"] = strconv.Itoa(nextCount)

	g.comp.Env.Vars["__PATHNAME__"] = pathName
	g.comp.Env.Vars["__DIRNAME__"] = filepath.Dir(pathName)
	g.comp.Env.Vars["__BASENAME__"] = filepath.Base(pathName)
}

func (g *Generator) Generate(targetPaths []string, entries FileSystemMap) (FileSystemMap, error) {
	g.generatedEntries = make(FileSystemMap)

	g.defineGlobalVars()

	for _, targetPath := range targetPaths {
		// allow debug
		if g.debug {
			cli.Logf(cli.LogSeverityWarn, "Generating files for directory %s ...\n", targetPath)
		}

		for pathName, entry := range entries {
			g.defineLocalVars(pathName, entry)

			result, err := g.comp.Compile(lang.NewBuffer("<pathname>", pathName))

			if err != nil {
				return nil, err
			}

			pathName = filepath.Join(targetPath, result)

			if entry.IsDir() {
				g.generateDirectory(pathName, entry)
			} else if err := g.generateFile(pathName, entry); err != nil {
				return nil, err
			}
		}
	}

	return g.generatedEntries, nil
}

func (g *Generator) generateDirectory(dirName string, entry FileSystemEntry) {
	// allow debug
	if g.debug {
		cli.Logf(cli.LogSeverityWarn, "Generating directory %s ...\n", dirName)
	}

	g.generatedEntries[dirName] = entry
}

func (g *Generator) generateFile(fileName string, entry FileSystemEntry) error {
	before, isCompilable := strings.CutSuffix(fileName, ".mimic")

	if isCompilable {
		fileName = before
	}

	// allow debug
	if g.debug {
		cli.Logf(cli.LogSeverityWarn, "Generating file %s ...\n", fileName)
	}

	if isCompilable {
		result, err := g.comp.Compile(lang.NewBuffer(fileName, string(entry.Data)))

		if err != nil {
			return err
		}

		entry.Data = []byte(result)
	}

	g.generatedEntries[fileName] = entry

	return nil
}
