package core

import (
	_ "embed"
	"fmt"
	"mimic/core/cli"
	"mimic/core/lang"
	"mimic/core/util"
	"os"
	"path/filepath"
	"strings"
)

type Writer struct {
	comp *lang.Compiler

	fileMap util.FileMap
}

func NewWriter(comp *lang.Compiler) *Writer {
	return &Writer{
		comp: comp,

		fileMap: make(util.FileMap),
	}
}

func (w *Writer) writeFile(fileName string, fileData string) {
	cli.Log(fmt.Sprintf("Writing file %s ...", fileName), cli.LogSeverityInfo)

	dirName := filepath.Dir(fileName)

	if err := os.MkdirAll(dirName, 0755); err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to create directory %s", dirName), cli.LogSeverityError)
	}

	if err := os.WriteFile(fileName, []byte(fileData), 0644); err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to write file %s", fileName), cli.LogSeverityError)
	}

	w.fileMap[fileName] = fileData
}

func (w *Writer) Write(targetPath string, fileMap util.FileMap) util.FileMap {
	for k := range w.fileMap {
		delete(w.fileMap, k)
	}

	cli.Log(fmt.Sprintf("Writing files to directory %s ...", targetPath), cli.LogSeverityInfo)

	for fileName, fileData := range fileMap {
		w.comp.Env.DefineScopeVars(fileName, fileData)

		if strings.Contains(fileName, ".mimic") {
			fileName = filepath.Join(targetPath, strings.TrimRight(fileName, ".mimic"))
			fileData = w.comp.Compile(lang.NewBuffer(fileName, fileData))
		} else {
			fileName = filepath.Join(targetPath, fileName)
		}

		fileName = w.comp.Compile(lang.NewBuffer("<fileName>", fileName))

		w.writeFile(fileName, fileData)
	}

	return w.fileMap
}
