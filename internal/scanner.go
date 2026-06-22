package internal

import (
	"mimic/internal/cli"
	"mimic/internal/util"
	"os"
	"path/filepath"
)

type Scanner struct {
	config *Config

	entryMap EntryMap

	debug bool
}

func NewScanner(config *Config, debug bool) *Scanner {
	return &Scanner{
		config: config,

		entryMap: make(EntryMap),

		debug: debug,
	}
}

func (s *Scanner) Scan() EntryMap {
	s.entryMap = make(EntryMap)

	for _, sourcePath := range s.config.SourcePath.Values {
		if s.debug {
			cli.Logf(cli.LogSeverityWarn, "Scanning source path %s ...\n", sourcePath)
		}

		sourceInfo, err := os.Stat(sourcePath)

		if err != nil {
			cli.Logf(cli.LogSeverityError, "Unable to obtain information about path %s\n", sourcePath)
			os.Exit(1)
		}

		if sourceInfo.IsDir() {
			s.scanSourceDirectory(sourcePath, sourcePath)
		} else {
			s.scanSourceFile(filepath.Dir(sourcePath), sourcePath, sourceInfo)
		}
	}

	for pathName, entry := range s.entryMap {
		if s.debug {
			cli.Logf(cli.LogSeverityInfo, "Scanned about %d bytes from %s\n", entry.Size, pathName)
		}

		cli.Printf(cli.Normal, cli.Cyan, "@ %s\n", pathName)
	}

	return s.entryMap
}

func (s *Scanner) scanSourceDirectory(basePath string, dirName string) {
	entries, err := util.DirectoryWalk(dirName)

	if err != nil {
		cli.Logf(cli.LogSeverityError, "Unable to walk into source directory %s\n", dirName)
		os.Exit(1)
	}

	if len(entries) == 0 {
		cli.Logf(cli.LogSeverityError, "No entries found in source directory %s\n", dirName)
		os.Exit(1)
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
		cli.Logf(cli.LogSeverityError, "Unable to determine relative path for %s\n", entry.Path)
		os.Exit(1)
	}

	if entry.Info.IsDir() {
		s.scanDirectoryEntry(relPath, entry)
	} else {
		s.scanFileEntry(relPath, entry)
	}
}

func (s *Scanner) scanDirectoryEntry(relPath string, entry util.DirectoryEntry) {
	if s.debug {
		cli.Logf(cli.LogSeverityWarn, "Scanning directory %s ...\n", relPath)
	}

	s.entryMap[relPath] = NewDirectoryEntry(relPath, entry.Info)
}

func (s *Scanner) scanFileEntry(relPath string, entry util.DirectoryEntry) {
	fileName := entry.Path

	if s.debug {
		cli.Logf(cli.LogSeverityWarn, "Scanning file %s ...\n", relPath)
	}

	fileData, err := os.ReadFile(fileName)

	if err != nil {
		cli.Logf(cli.LogSeverityError, "Unable to obtain data from file %s\n", fileName)
		os.Exit(1)
	}

	s.entryMap[relPath] = NewFileEntry(relPath, entry.Info, fileData)
}
