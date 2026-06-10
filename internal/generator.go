package internal

import (
	"maps"
	"mimic/internal/cli"
	"mimic/internal/lang"
	"os"
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

	expr := lang.NewExpressionConfigurable(config.ExprOpen, config.ExprClose)

	comp := lang.NewCompilerConfigurable(env, expr, false)

	return &Generator{
		config: config,
		comp:   comp,

		entryMap: make(EntryMap),
	}
}

func (g *Generator) defineGlobalVars() {
	g.comp.Env.Vars["__SOURCEPATH__"] = g.config.SourcePath.String()
	g.comp.Env.Vars["__TARGETPATH__"] = g.config.TargetPath.String()
}

func (g *Generator) defineLocalVars(pathName string, fileInfo os.FileInfo, data []byte) {
	g.comp.Env.Vars["__BASEPATH__"] = filepath.Dir(pathName)
	g.comp.Env.Vars["__BASENAME__"] = filepath.Base(pathName)

	delete(g.comp.Env.Vars, "__DIRNAME__")
	delete(g.comp.Env.Vars, "__FILENAME__")
	delete(g.comp.Env.Vars, "__FILEDATA__")

	if fileInfo.IsDir() {
		g.comp.Env.Vars["__DIRNAME__"] = pathName
	} else {
		g.comp.Env.Vars["__FILENAME__"] = pathName
		g.comp.Env.Vars["__FILEDATA__"] = string(data)
	}
}

func (g *Generator) Generate(entryMap EntryMap) EntryMap {
	g.entryMap = make(EntryMap)

	g.defineGlobalVars()

	for _, targetPath := range g.config.TargetPath.Values {
		if g.config.DebugMode {
			cli.Logf(cli.LogSeverityWarn, "Generating files for directory %s ...\n", targetPath)
		}

		for pathName, entry := range entryMap {
			g.defineLocalVars(pathName, entry.Info, entry.Data)

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
	if g.config.DebugMode {
		cli.Logf(cli.LogSeverityWarn, "Generating directory %s ...\n", dirName)
	}

	g.entryMap[dirName] = entry
}

func (g *Generator) generateFile(fileName string, entry Entry) {
	before, isCompilable := strings.CutSuffix(fileName, ".mimic")

	if isCompilable {
		fileName = before
	}

	if g.config.DebugMode {
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
