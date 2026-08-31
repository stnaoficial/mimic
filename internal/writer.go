package internal

import (
	"fmt"
	"mimic/internal/cli"
	"os"
	"path/filepath"
	"strings"
)

type Writer struct {
	entryMap EntryMap

	debug bool
}

func NewWriter(debug bool) *Writer {
	return &Writer{
		entryMap: make(EntryMap),

		debug: debug,
	}
}

func (w *Writer) Write(entryMap EntryMap) EntryMap {
	for pathName, entry := range entryMap {
		if entry.IsDir() {
			w.writeDirectory(pathName, entry)
		} else {
			w.writeFile(pathName, entry)
		}
	}

	for pathName, entry := range w.entryMap {
		if w.debug {
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
	if w.debug {
		cli.Logf(cli.LogSeverityWarn, "Writing directory %s ...\n", dirName)
	}

	if err := os.MkdirAll(dirName, 0755); err != nil {
		cli.Logf(cli.LogSeverityError, "Unable to create directory %s\n", dirName)
		os.Exit(1)
	}

	w.entryMap[dirName] = entry
}

func (w *Writer) writeFile(fileName string, entry Entry) {
	if w.debug {
		cli.Logf(cli.LogSeverityWarn, "Writing file %s ...\n", fileName)
	}

	dirName := filepath.Dir(fileName)

	if err := os.MkdirAll(dirName, 0755); err != nil {
		cli.Logf(cli.LogSeverityError, "Unable to create directory %s\n", dirName)
		os.Exit(1)
	}

	var data []byte

	if fileData, err := os.ReadFile(fileName); err == nil {
		data = fileData
	}

	data = append(data, entry.Data...)

	if err := os.WriteFile(fileName, data, 0644); err != nil {
		cli.Logf(cli.LogSeverityError, "Unable to write file %s\n", fileName)
		os.Exit(1)
	}

	w.entryMap[fileName] = entry
}
