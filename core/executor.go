package core

import (
	"mimic/core/cli"
	"mimic/core/lang"
	"mimic/core/util"
	"os"
)

type Executor struct {
	reader    *Reader
	FilesRead util.FileMap

	writer       *Writer
	WrittenFiles util.FileMap

	source string
	target string
}

func NewExecutor(source string, target string, comp *lang.Compiler) *Executor {
	return &Executor{
		reader:    NewReader(),
		FilesRead: make(util.FileMap),

		writer:       NewWriter(comp),
		WrittenFiles: make(util.FileMap),

		source: source,
		target: target,
	}
}

func (e *Executor) Read() {
	e.FilesRead = e.reader.Read(e.source)
}

func (e *Executor) Write() {
	e.WrittenFiles = e.writer.Write(e.target, e.FilesRead)
}

func (e *Executor) Init() {
	_, err := os.Getwd()

	if err != nil {
		cli.LogAndExit("Unable to get working directory", cli.LogSeverityError)
	}
}
