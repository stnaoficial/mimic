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
		cli.Logf(cli.LogSeverityWarn, "Writing files to directory %s ...\n", w.config.TargetPath)
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
			cli.Logf(cli.LogSeveritySuccess, "Wrote about %d bytes at %s\n", entry.Size, pathName)
		}

		cli.Printf(cli.Normal, cli.Cyan, "@ %s\n", pathName)

		for line := range strings.SplitSeq(string(entry.Data), "\n") {
			cli.Printf(cli.Normal, cli.Green, "+ %s\n", line)
		}

		fmt.Println()
	}

	return w.entryMap
}

func (w *Writer) writeDirectory(dirName string, entry Entry) {
	if w.config.DebugMode {
		cli.Logf(cli.LogSeverityWarn, "Writing directory %s ...\n", dirName)
	}

	if err := os.MkdirAll(dirName, 0755); err != nil {
		cli.Logf(cli.LogSeverityError, "Unable to create directory %s\n", dirName)
		os.Exit(1)
	}

	w.entryMap[dirName] = entry
}

func (w *Writer) writeFile(fileName string, entry Entry) {
	if w.config.DebugMode {
		cli.Logf(cli.LogSeverityWarn, "Writing file %s ...\n", fileName)
	}

	dirName := filepath.Dir(fileName)

	if err := os.MkdirAll(dirName, 0755); err != nil {
		cli.Logf(cli.LogSeverityError, "Unable to create directory %s\n", dirName)
		os.Exit(1)
	}

	if err := os.WriteFile(fileName, entry.Data, 0644); err != nil {
		cli.Logf(cli.LogSeverityError, "Unable to write file %s\n", fileName)
		os.Exit(1)
	}

	w.entryMap[fileName] = entry
}
