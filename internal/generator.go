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

	debug bool
}

func NewGenerator(comp *lang.Compiler, debug bool) *Generator {
	return &Generator{
		comp: comp,

		debug: debug,
	}
}

func (g *Generator) defineGlobalVars() {
	// TODO
}

func (g *Generator) defineLocalVars(pathName string, entry Entry) {
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

func (g *Generator) Generate(targetPaths []string, inputEntryMap EntryMap) EntryMap {
	outputEntryMap := make(EntryMap)

	g.defineGlobalVars()

	for _, targetPath := range targetPaths {
		if g.debug {
			cli.Logf(cli.LogSeverityWarn, "Generating files for directory %s ...\n", targetPath)
		}

		for pathName, entry := range inputEntryMap {
			g.defineLocalVars(pathName, entry)

			result, err := g.comp.Compile(lang.NewBuffer("<pathname>", pathName))

			if err != nil {
				cli.Logln(cli.LogSeverityError, err.Error())
				os.Exit(0)
			}

			pathName = filepath.Join(targetPath, result)

			if entry.IsDir() {
				g.generateDirectory(outputEntryMap, pathName, entry)
			} else {
				g.generateFile(outputEntryMap, pathName, entry)
			}
		}
	}

	return outputEntryMap
}

func (g *Generator) generateDirectory(outputEntryMap EntryMap, dirName string, entry Entry) {
	if g.debug {
		cli.Logf(cli.LogSeverityWarn, "Generating directory %s ...\n", dirName)
	}

	outputEntryMap[dirName] = entry
}

func (g *Generator) generateFile(outputEntryMap EntryMap, fileName string, entry Entry) {
	before, isCompilable := strings.CutSuffix(fileName, ".mimic")

	if isCompilable {
		fileName = before
	}

	if g.debug {
		cli.Logf(cli.LogSeverityWarn, "Generating file %s ...\n", fileName)
	}

	if isCompilable {
		result, err := g.comp.Compile(lang.NewBuffer(fileName, string(entry.Data)))

		if err != nil {
			cli.Logln(cli.LogSeverityError, err.Error())
			os.Exit(0)
		}

		entry.Data = []byte(result)
	}

	outputEntryMap[fileName] = entry
}
