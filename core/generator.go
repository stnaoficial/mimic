package core

import (
	"fmt"
	"mimic/core/cli"
	"mimic/core/lang"
	"mimic/core/util"
	"os"
	"path/filepath"
	"slices"
)

type Generator struct {
	comp *lang.Compiler

	source     string
	sourceFile os.FileInfo
	target     string
	targetFile os.FileInfo
	fileMap    map[string][]byte

	GeneratedFiles []string
}

func NewGenerator(source string, target string, env *lang.Environment, expr *lang.Expression) *Generator {
	sourceFile, err := os.Stat(source)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to obtain information about path %s", source), cli.LogSeverityError)
	}

	targetFile, err := os.Stat(target)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to obtain information about path %s", target), cli.LogSeverityError)
	}

	return &Generator{
		comp: lang.NewCompiler(env, expr),

		source:     source,
		sourceFile: sourceFile,
		target:     target,
		targetFile: targetFile,
		fileMap:    make(map[string][]byte),

		GeneratedFiles: []string{},
	}
}

func (g *Generator) scanDirectory(dirname string) {
	cli.Log(fmt.Sprintf("Scanning directory %s...", g.source), cli.LogSeverityInfo)

	GeneratedFiles, err := util.DirectoryWalk(g.source)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to walk into directory %s", g.source), cli.LogSeverityError)
	}

	if len(GeneratedFiles) == 0 {
		cli.LogAndExit(fmt.Sprintf("No .mimic files found in directory %s", g.source), cli.LogSeverityWarn)
	}

	g.GeneratedFiles = GeneratedFiles

	for _, filename := range GeneratedFiles {
		g.scanFile(filename)
	}
}

func (g *Generator) scanFile(filename string) {
	cli.Log(fmt.Sprintf("Scanning file %s...", g.source), cli.LogSeverityInfo)

	if !slices.Contains(g.GeneratedFiles, filename) {
		g.GeneratedFiles = append(g.GeneratedFiles, filename)
	}

	data, err := os.ReadFile(filename)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to obtain data from file %s", filename), cli.LogSeverityError)
	}

	g.fileMap[filename] = data
}

func (g *Generator) Scan() {
	if g.sourceFile.IsDir() {
		g.scanDirectory(g.source)
	} else {
		g.scanFile(g.source)
	}
}

func (g *Generator) Copy() {
	cli.Log(fmt.Sprintf("Copying files to directory %s...", g.target), cli.LogSeverityInfo)

	for filename, data := range g.fileMap {
		var basepath = g.source

		if !g.sourceFile.IsDir() {
			basepath = filepath.Dir(basepath)
		}

		reldir, err := filepath.Rel(basepath, filename)

		if err != nil {
			cli.LogAndExit(fmt.Sprintf("Unable to relate directory %s to file %s", g.source, filename), cli.LogSeverityError)
		}

		filename = filepath.Join(g.target, reldir[:len(reldir)-len(".mimic")])
		filename = g.comp.Compile(lang.NewBuffer("<filename>", []byte(filename)))

		dirname := filepath.Dir(filename)

		if err := os.MkdirAll(dirname, 0755); err != nil {
			cli.LogAndExit(fmt.Sprintf("Unable to create file %s", filename), cli.LogSeverityError)
		}

		result := g.comp.Compile(lang.NewBuffer(filename, data))

		if err := os.WriteFile(filename, []byte(result), 0644); err != nil {
			cli.LogAndExit(fmt.Sprintf("Unable to create file %s", filename), cli.LogSeverityError)
		}

		cli.LogFileNameAt(filename)
		cli.LogFileDataAdded(result)
	}
}
