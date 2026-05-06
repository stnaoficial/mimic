package core

import (
	"fmt"
	"maps"
	"mimic/core/cli"
	"mimic/core/lang"
	"mimic/core/util"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Generator struct {
	config *Config
	comp   *lang.Compiler

	entryMap EntryMap
}

func NewGenerator(config *Config) *Generator {
	env := lang.NewEnvironment()
	maps.Copy(env.Vars, config.Variables)

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

	g.defineGlobalVars()

	for pathName, entry := range entryMap {
		g.defineLocalVars(pathName, entry)

		pathName = g.comp.Compile(lang.NewBuffer("<pathname>", pathName))
		pathName = filepath.Join(g.config.TargetPath, pathName)

		if entry.IsDir() {
			g.generateDirectory(pathName, entry)
		} else {
			g.generateFile(pathName, entry)
		}
	}

	return g.entryMap
}

func (g *Generator) defineGlobalVars() {
	now := time.Now()

	g.comp.Env.Vars["__TIMESTAMP__"] = fmt.Sprintf("%d", now.Unix())

	g.comp.Env.Vars["__DATE__"] = now.Format("2006-01-02")
	g.comp.Env.Vars["__TIME__"] = now.Format("15:04:05")
	g.comp.Env.Vars["__DATETIME__"] = now.Format("2006-01-02T15:04:05Z")

	g.comp.Env.Vars["__YEAR__"] = now.Format("2006")
	g.comp.Env.Vars["__MONTH__"] = now.Format("01")
	g.comp.Env.Vars["__DAY__"] = now.Format("02")

	g.comp.Env.Vars["__HOUR__"] = now.Format("15")
	g.comp.Env.Vars["__MINUTE__"] = now.Format("04")
	g.comp.Env.Vars["__SECOND__"] = now.Format("05")

	ns := now.Nanosecond()

	g.comp.Env.Vars["__MILLISECOND__"] = fmt.Sprintf("%03d", ns/1_000_000)
	g.comp.Env.Vars["__MICROSECOND__"] = fmt.Sprintf("%06d", ns/1_000)
	g.comp.Env.Vars["__NANOSECOND__"] = fmt.Sprintf("%09d", ns)
}

func (g *Generator) defineLocalVars(pathName string, entry Entry) {
	g.comp.Env.Vars["__UID__"] = uuid.NewString()

	g.comp.Env.Vars["__16_DIGIT__"] = util.RandDigit(16)
	g.comp.Env.Vars["__8_DIGIT__"] = util.RandDigit(8)
	g.comp.Env.Vars["__4_DIGIT__"] = util.RandDigit(4)
	g.comp.Env.Vars["__2_DIGIT__"] = util.RandDigit(2)

	g.comp.Env.Vars["__BASEPATH__"] = filepath.Dir(pathName)
	g.comp.Env.Vars["__BASENAME__"] = filepath.Base(pathName)

	delete(g.comp.Env.Vars, "__DIRNAME__")
	delete(g.comp.Env.Vars, "__FILENAME__")
	delete(g.comp.Env.Vars, "__FILEDATA__")

	if entry.IsDir() {
		g.comp.Env.Vars["__DIRNAME__"] = pathName
	} else {
		g.comp.Env.Vars["__FILENAME__"] = pathName
		g.comp.Env.Vars["__FILEDATA__"] = string(entry.Data)
	}
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
		entry.Data = []byte(g.comp.Compile(lang.NewBuffer(fileName, string(entry.Data))))
	}

	g.entryMap[fileName] = entry
}
