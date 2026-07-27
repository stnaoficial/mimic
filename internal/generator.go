package internal

import (
	"maps"
	"mimic/internal/cli"
	"mimic/internal/lang"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Generator struct {
	config *Config
	comp   *lang.Compiler

	entryMap EntryMap

	debug  bool
	strict bool
}

func NewGenerator(config *Config, debug bool, strict bool) *Generator {
	env := lang.NewEnvironment()
	maps.Copy(env.Vars, config.Variables)
	maps.Copy(env.Prompts, config.Prompts)

	expr := lang.NewExpressionConfigurable(config.ExprOpen, config.ExprClose)

	comp := lang.NewCompilerConfigurable(env, expr, strict)

	return &Generator{
		config: config,
		comp:   comp,

		entryMap: make(EntryMap),

		debug:  debug,
		strict: strict,
	}
}

func (g *Generator) defineGlobalVars() {
	g.comp.Env.Vars["__SOURCE_PATH__"] = g.config.SourcePath.String()
	g.comp.Env.Vars["__TARGET_PATH__"] = g.config.TargetPath.String()
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

func (g *Generator) Generate(entryMap EntryMap) EntryMap {
	g.entryMap = make(EntryMap)

	g.defineGlobalVars()

	for _, targetPath := range g.config.TargetPath.Values {
		if g.debug {
			cli.Logf(cli.LogSeverityWarn, "Generating files for directory %s ...\n", targetPath)
		}

		for pathName, entry := range entryMap {
			g.defineLocalVars(pathName, entry)

			result, err := g.comp.Compile(lang.NewBuffer("<pathname>", pathName))

			if err != nil {
				cli.Logln(cli.LogSeverityError, err.Error())
				os.Exit(0)
			}

			pathName = filepath.Join(targetPath, result)

			if entry.IsDir() {
				g.generateDirectory(pathName, entry)
			} else {
				g.generateFile(pathName, entry)
			}
		}
	}

	return g.entryMap
}

func (g *Generator) generateDirectory(dirName string, entry Entry) {
	if g.debug {
		cli.Logf(cli.LogSeverityWarn, "Generating directory %s ...\n", dirName)
	}

	g.entryMap[dirName] = entry
}

func (g *Generator) generateFile(fileName string, entry Entry) {
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

	g.entryMap[fileName] = entry
}
