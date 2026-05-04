package core

import (
	"fmt"
	"mimic/core/cli"
	"mimic/core/util"
	"os"
	"path/filepath"
)

type Scanner struct {
	entryMap EntryMap
}

func NewScanner() *Scanner {
	return &Scanner{
		entryMap: make(EntryMap),
	}
}

func (s *Scanner) Scan(sourcePath string) EntryMap {
	cli.Log(fmt.Sprintf("Scanning source path %s ...", sourcePath), cli.LogSeverityInfo)

	s.entryMap = make(EntryMap)

	sourceInfo, err := os.Stat(sourcePath)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to obtain information about path %s", sourcePath), cli.LogSeverityError)
	}

	if sourceInfo.IsDir() {
		s.scanSourceDirectory(sourcePath, sourcePath)
	} else {
		s.scanSourceFile(filepath.Dir(sourcePath), sourcePath, sourceInfo)
	}

	return s.entryMap
}

func (s *Scanner) scanSourceDirectory(basePath string, dirName string) {
	entries, err := util.DirectoryWalk(dirName)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to walk into source directory %s", dirName), cli.LogSeverityError)
	}

	if len(entries) == 0 {
		cli.LogAndExit(fmt.Sprintf("No entries found in source directory %s", dirName), cli.LogSeverityError)
	}

	for _, entry := range entries {
		s.scanEntry(basePath, entry)
	}
}

func (s *Scanner) scanSourceFile(basePath string, fileName string, fileInfo os.FileInfo) {
	s.scanEntry(basePath, util.DirectoryEntry{Path: fileName, Info: fileInfo})
}

func (s *Scanner) scanEntry(basePath string, entry util.DirectoryEntry) {
	if basePath == entry.Path {
		return
	}

	relPath, err := filepath.Rel(basePath, entry.Path)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to determine relative path for %s", entry.Path), cli.LogSeverityError)
	}

	if entry.Info.IsDir() {
		s.scanDirectoryEntry(relPath)
	} else {
		s.scanFileEntry(relPath, entry.Path)
	}
}

func (s *Scanner) scanDirectoryEntry(relPath string) {
	cli.Log(fmt.Sprintf("Scanning directory %s ...", relPath), cli.LogSeverityInfo)

	s.entryMap[relPath] = NewDirectoryEntry(relPath)
}

func (s *Scanner) scanFileEntry(relPath string, fileName string) {
	cli.Log(fmt.Sprintf("Scanning file %s ...", relPath), cli.LogSeverityInfo)

	fileData, err := os.ReadFile(fileName)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to obtain data from file %s", fileName), cli.LogSeverityError)
	}

	s.entryMap[relPath] = NewFileEntry(relPath, fileData)
}
