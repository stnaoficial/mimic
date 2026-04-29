package core

import (
	"fmt"
	"mimic/core/cli"
	"mimic/core/lang"
	"mimic/core/util"
	"os"
)

type Executor struct {
	comp *lang.Compiler

	reader    *Reader
	FilesRead util.FileMap

	writer       *Writer
	WrittenFiles util.FileMap

	source     string
	sourceInfo os.FileInfo

	target string
}

func NewExecutor(source string, target string, comp *lang.Compiler) *Executor {
	sourceInfo, err := os.Stat(source)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to obtain information about path %s", source), cli.LogSeverityError)
	}

	return &Executor{
		comp: comp,

		reader:    NewReader(),
		FilesRead: make(util.FileMap),

		writer:       NewWriter(comp),
		WrittenFiles: make(util.FileMap),

		source:     source,
		sourceInfo: sourceInfo,

		target: target,
	}
}

func (e *Executor) Read() {
	e.FilesRead = e.reader.Read(e.source, e.sourceInfo)
}

func (e *Executor) Write() {
	e.WrittenFiles = e.writer.Write(e.source, e.sourceInfo, e.target, e.comp, e.FilesRead)
}
