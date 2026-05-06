package core

import (
	"fmt"
	"mimic/core/cli"
	"os"
	"path/filepath"
	"strings"
)

type Writer struct {
	config *Config

	entryMap EntryMap
}

func NewWriter(config *Config) *Writer {
	return &Writer{
		config: config,

		entryMap: make(EntryMap),
	}
}

func (w *Writer) Write(entryMap EntryMap) EntryMap {
	if w.config.DebugMode {
		cli.LogWithPrefix(fmt.Sprintf("Writing files to directory %s ...", w.config.TargetPath), cli.LogSeverityInfo)
	}

	for pathName, entry := range entryMap {
		if entry.IsDir() {
			w.writeDirectory(pathName, entry)
		} else {
			w.writeFile(pathName, entry)
		}
	}

	for pathName, entry := range w.entryMap {
		if w.config.DebugMode {
			cli.LogWithPrefix(fmt.Sprintf("Wrote about %d bytes at %s", entry.Size, pathName), cli.LogSeveritySuccess)
		}

		cli.Log(fmt.Sprintf("@ %s", pathName), cli.LogSeverityInfo)

		for line := range strings.SplitSeq(string(entry.Data), "\n") {
			cli.Log(fmt.Sprintf("+ %s", line), cli.LogSeveritySuccess)
		}

		fmt.Println()
	}

	return w.entryMap
}

func (w *Writer) writeDirectory(dirName string, entry Entry) {
	if w.config.DebugMode {
		cli.LogWithPrefix(fmt.Sprintf("Writing directory %s ...", dirName), cli.LogSeverityInfo)
	}

	if err := os.MkdirAll(dirName, 0755); err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to create directory %s", dirName), cli.LogSeverityError)
	}

	w.entryMap[dirName] = entry
}

func (w *Writer) writeFile(fileName string, entry Entry) {
	if w.config.DebugMode {
		cli.LogWithPrefix(fmt.Sprintf("Writing file %s ...", fileName), cli.LogSeverityInfo)
	}

	dirName := filepath.Dir(fileName)

	if err := os.MkdirAll(dirName, 0755); err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to create directory %s", dirName), cli.LogSeverityError)
	}

	if err := os.WriteFile(fileName, entry.Data, 0644); err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to write file %s", fileName), cli.LogSeverityError)
	}

	w.entryMap[fileName] = entry
}
