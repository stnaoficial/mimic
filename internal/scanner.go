package internal

import (
	"mimic/internal/cli"
	"mimic/internal/util"
	"os"
	"path/filepath"
)

type Scanner struct {
	debug bool
}

func NewScanner(debug bool) *Scanner {
	return &Scanner{
		debug: debug,
	}
}

func (s *Scanner) Scan(sourcePaths []string) EntryMap {
	entryMap := make(EntryMap)

	for _, sourcePath := range sourcePaths {
		if s.debug {
			cli.Logf(cli.LogSeverityWarn, "Scanning source path %s ...\n", sourcePath)
		}

		sourceInfo, err := os.Stat(sourcePath)

		if err != nil {
			cli.Logf(cli.LogSeverityError, "Unable to obtain information about path %s\n", sourcePath)
			os.Exit(1)
		}

		if sourceInfo.IsDir() {
			s.scanSourceDirectory(entryMap, sourcePath, sourcePath)
		} else {
			s.scanSourceFile(entryMap, filepath.Dir(sourcePath), sourcePath, sourceInfo)
		}
	}

	if s.debug {
		for pathName, entry := range entryMap {
			cli.Logf(cli.LogSeverityInfo, "Scanned about %d bytes from %s\n", entry.Size, pathName)
		}
	}

	return entryMap
}

func (s *Scanner) scanSourceDirectory(entryMap EntryMap, basePath string, dirName string) {
	entries, err := util.FileSystemWalk(dirName)

	if err != nil {
		cli.Logf(cli.LogSeverityError, "Unable to walk into source directory %s\n", dirName)
		os.Exit(1)
	}

	if len(entries) == 0 {
		cli.Logf(cli.LogSeverityError, "No entries found in source directory %s\n", dirName)
		os.Exit(1)
	}

	for _, entry := range entries {
		s.scanEntry(entryMap, basePath, entry)
	}
}

func (s *Scanner) scanSourceFile(entryMap EntryMap, basePath string, fileName string, fileInfo os.FileInfo) {
	s.scanEntry(entryMap, basePath, util.FileSystemEntry{Path: fileName, Info: fileInfo})
}

func (s *Scanner) scanEntry(entryMap EntryMap, basePath string, entry util.FileSystemEntry) {
	if basePath == entry.Path {
		return
	}

	relPath, err := filepath.Rel(basePath, entry.Path)

	if err != nil {
		cli.Logf(cli.LogSeverityError, "Unable to determine relative path for %s\n", entry.Path)
		os.Exit(1)
	}

	if entry.Info.IsDir() {
		s.scanDirectoryEntry(entryMap, relPath, entry)
	} else {
		s.scanFileEntry(entryMap, relPath, entry)
	}
}

func (s *Scanner) scanDirectoryEntry(entryMap EntryMap, relPath string, entry util.FileSystemEntry) {
	if s.debug {
		cli.Logf(cli.LogSeverityWarn, "Scanning directory %s ...\n", relPath)
	}

	entryMap[relPath] = NewDirectoryEntry(relPath, entry.Info)
}

func (s *Scanner) scanFileEntry(entryMap EntryMap, relPath string, entry util.FileSystemEntry) {
	fileName := entry.Path

	if s.debug {
		cli.Logf(cli.LogSeverityWarn, "Scanning file %s ...\n", relPath)
	}

	fileData, err := os.ReadFile(fileName)

	if err != nil {
		cli.Logf(cli.LogSeverityError, "Unable to obtain data from file %s\n", fileName)
		os.Exit(1)
	}

	entryMap[relPath] = NewFileEntry(relPath, entry.Info, fileData)
}
