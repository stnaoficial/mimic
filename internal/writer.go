package internal

import (
	"fmt"
	"mimic/internal/cli"
	"os"
	"path/filepath"
	"strings"
)

type WriteMode int

const (
	WriteModeOverride WriteMode = iota
	WriteModeAppend
)

type Writer struct {
	writtenEntries FileSystemMap

	debug bool
	mode  WriteMode
}

func NewWriter(debug bool, mode WriteMode) *Writer {
	return &Writer{
		writtenEntries: make(FileSystemMap),

		debug: debug,
		mode:  mode,
	}
}

func (w *Writer) Write(entries FileSystemMap) (FileSystemMap, error) {
	w.writtenEntries = make(FileSystemMap)

	for pathName, entry := range entries {
		if entry.IsDir() {
			if err := w.writeDirectory(pathName, entry); err != nil {
				return nil, err
			}
		} else if err := w.writeFile(pathName, entry); err != nil {
			return nil, err
		}
	}

	for pathName, entry := range w.writtenEntries {
		// allow debug
		if w.debug {
			cli.Logf(cli.LogSeveritySuccess, "Wrote about %d bytes at %s\n", entry.Size, pathName)
		}

		cli.Printf(cli.Normal, cli.Cyan, "@ %s\n", pathName)

		for line := range strings.SplitSeq(string(entry.Data), "\n") {
			cli.Printf(cli.Normal, cli.Green, "+ %s\n", line)
		}

		fmt.Println()
	}

	return w.writtenEntries, nil
}

func (w *Writer) writeDirectory(dirName string, entry FileSystemEntry) error {
	// allow debug
	if w.debug {
		cli.Logf(cli.LogSeverityWarn, "Writing directory %s ...\n", dirName)
	}

	if err := os.MkdirAll(dirName, 0755); err != nil {
		return err
	}

	w.writtenEntries[dirName] = entry

	return nil
}

func (w *Writer) writeFile(fileName string, entry FileSystemEntry) error {
	// allow debug
	if w.debug {
		cli.Logf(cli.LogSeverityWarn, "Writing file %s ...\n", fileName)
	}

	dirName := filepath.Dir(fileName)

	if err := os.MkdirAll(dirName, 0755); err != nil {
		return err
	}

	var data []byte

	if fileData, err := os.ReadFile(fileName); err == nil {
		data = fileData
	}

	switch w.mode {
	case WriteModeAppend:
		data = append(data, entry.Data...)
	case WriteModeOverride:
		data = entry.Data
	}

	if err := os.WriteFile(fileName, data, 0644); err != nil {
		return err
	}

	w.writtenEntries[fileName] = entry

	return nil
}
