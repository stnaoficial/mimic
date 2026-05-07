package core

import (
	"fmt"
	"maps"
	"mimic/core/cli"
	"mimic/core/lang"
	"path/filepath"
	"strings"
)

type Generator struct {
	config *Config
	comp   *lang.Compiler

	entryMap EntryMap
}

func NewGenerator(config *Config) *Generator {
	env := lang.NewEnvironment()
	maps.Copy(env.Vars, config.Variables)
	maps.Copy(env.Prompts, config.Prompts)

	comp := lang.NewCompilerConfigurable(env, lang.NewExpressionConfigurable(config.ExprOpen, config.ExprClose))

	return &Generator{
		config: config,
		comp:   comp,

		entryMap: make(EntryMap),
	}
}

func (g *Generator) Generate(entryMap EntryMap) EntryMap {
	if g.config.DebugMode {
		cli.LogWithPrefix(fmt.Sprintf("Generating files for directory %s ...", g.config.TargetPath), cli.LogSeverityInfo)
	}

	g.entryMap = make(EntryMap)

	g.comp.Env.DefineGlobalVars()

	for pathName, entry := range entryMap {
		g.comp.Env.DefineLocalVars(pathName, entry.Info, entry.Data)

		result, err := g.comp.Compile(lang.NewBuffer("<pathname>", pathName))

		if err != nil {
			cli.LogAndExit(err.Error(), cli.LogSeverityError)
		}

		pathName = filepath.Join(g.config.TargetPath, result)

		if entry.IsDir() {
			g.generateDirectory(pathName, entry)
		} else {
			g.generateFile(pathName, entry)
		}
	}

	return g.entryMap
}

func (g *Generator) generateDirectory(dirName string, entry Entry) {
	if g.config.DebugMode {
		cli.LogWithPrefix(fmt.Sprintf("Generating directory %s ...", dirName), cli.LogSeverityInfo)
	}

	g.entryMap[dirName] = entry
}

func (g *Generator) generateFile(fileName string, entry Entry) {
	before, isCompilable := strings.CutSuffix(fileName, ".mimic")

	if isCompilable {
		fileName = before
	}

	if g.config.DebugMode {
		cli.LogWithPrefix(fmt.Sprintf("Generating file %s ...", fileName), cli.LogSeverityInfo)
	}

	if isCompilable {
		result, err := g.comp.Compile(lang.NewBuffer(fileName, string(entry.Data)))

		if err != nil {
			cli.LogAndExit(err.Error(), cli.LogSeverityError)
		}

		entry.Data = []byte(result)
	}

	g.entryMap[fileName] = entry
}
