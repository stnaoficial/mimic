package core

import (
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

func (w *Writer) Write(source string, sourceInfo os.FileInfo, target string, comp *lang.Compiler, fileMap util.FileMap) util.FileMap {
	cli.Log(fmt.Sprintf("Writing files to directory %s...", target), cli.LogSeverityInfo)

	for k := range w.fileMap {
		delete(w.fileMap, k)
	}

	for filename, filedata := range fileMap {
		var basepath = source

		if !sourceInfo.IsDir() {
			basepath = filepath.Dir(basepath)
		}

		reldir, err := filepath.Rel(basepath, filename)

		if err != nil {
			cli.LogAndExit(fmt.Sprintf("Unable to relate directory %s to file %s", source, filename), cli.LogSeverityError)
		}

		if strings.Contains(reldir, ".mimic") {
			filename = filepath.Join(target, reldir[:len(reldir)-len(".mimic")])
			filedata = w.comp.Compile(lang.NewBuffer(filename, filedata))
		} else {
			filename = filepath.Join(target, reldir)
		}

		filename = w.comp.Compile(lang.NewBuffer("<filename>", filename))

		dirname := filepath.Dir(filename)

		cli.Log(fmt.Sprintf("Writing file %s...", filename), cli.LogSeverityInfo)

		if err := os.MkdirAll(dirname, 0755); err != nil {
			cli.LogAndExit(fmt.Sprintf("Unable to create directory %s", dirname), cli.LogSeverityError)
		}

		if err := os.WriteFile(filename, []byte(filedata), 0644); err != nil {
			cli.LogAndExit(fmt.Sprintf("Unable to write file %s", filename), cli.LogSeverityError)
		}

		w.fileMap[filename] = filedata
	}

	return w.fileMap
}
