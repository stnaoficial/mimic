package core

import "mimic/core/lang"

type Executor struct {
	scanner   *Scanner
	FilesRead EntryMap

	writer       *Writer
	WrittenFiles EntryMap

	sourcePath string
	targetPath string
}

func NewExecutor(sourcePath string, targetPath string, comp *lang.Compiler) *Executor {
	return &Executor{
		scanner:   NewScanner(),
		FilesRead: make(EntryMap),

		writer:       NewWriter(comp),
		WrittenFiles: make(EntryMap),

		sourcePath: sourcePath,
		targetPath: targetPath,
	}
}

func (e *Executor) Scan() {
	e.FilesRead = e.scanner.Scan(e.sourcePath)
}

func (e *Executor) Write() {
	e.WrittenFiles = e.writer.Write(e.targetPath, e.FilesRead)
}
