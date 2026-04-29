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

func (w *Writer) write(filename string, filedata string) {
	cli.Log(fmt.Sprintf("Writing file %s ...", filename), cli.LogSeverityInfo)

	dirname := filepath.Dir(filename)

	if err := os.MkdirAll(dirname, 0755); err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to create directory %s", dirname), cli.LogSeverityError)
	}

	if err := os.WriteFile(filename, []byte(filedata), 0644); err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to write file %s", filename), cli.LogSeverityError)
	}

	w.fileMap[filename] = filedata
}

func (w *Writer) Write(target string, fileMap util.FileMap) util.FileMap {
	cli.Log(fmt.Sprintf("Writing files to directory %s ...", target), cli.LogSeverityInfo)

	for k := range w.fileMap {
		delete(w.fileMap, k)
	}

	for filename, filedata := range fileMap {
		w.comp.Env.Vars["__DIRNAME__"] = filepath.Dir(filename)
		w.comp.Env.Vars["__FILENAME__"] = filename

		if strings.Contains(filename, ".mimic") {
			filename = filepath.Join(target, strings.TrimRight(filename, ".mimic"))

			w.comp.Env.Vars["__FILENAME__"] = filename

			filedata = w.comp.Compile(lang.NewBuffer(filename, filedata))
		} else {
			filename = filepath.Join(target, filename)
		}

		filename = w.comp.Compile(lang.NewBuffer("<filename>", filename))

		w.write(filename, filedata)
	}

	return w.fileMap
}
