package core

import (
	"mimic/core/lang"
	"mimic/core/util"
)

type Executor struct {
	reader    *Reader
	FilesRead util.FileMap

	writer       *Writer
	WrittenFiles util.FileMap

	sourcePath string
	targetPath string
}

func NewExecutor(sourcePath string, targetPath string, comp *lang.Compiler) *Executor {
	return &Executor{
		reader:    NewReader(),
		FilesRead: make(util.FileMap),

		writer:       NewWriter(comp),
		WrittenFiles: make(util.FileMap),

		sourcePath: sourcePath,
		targetPath: targetPath,
	}
}

func (e *Executor) Read() {
	e.FilesRead = e.reader.Read(e.sourcePath)
}

func (e *Executor) Write() {
	e.WrittenFiles = e.writer.Write(e.targetPath, e.FilesRead)
}
