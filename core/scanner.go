package core

import (
	"fmt"
	"mimic/core/cli"
	"mimic/core/util"
	"os"
	"path/filepath"
)

type Scanner struct {
	config *Config

	entryMap EntryMap
}

func NewScanner(config *Config) *Scanner {
	return &Scanner{
		config: config,

		entryMap: make(EntryMap),
	}
}

func (s *Scanner) Scan() EntryMap {
	if s.config.DebugMode {
		cli.LogWithPrefix(fmt.Sprintf("Scanning source path %s ...", s.config.SourcePath), cli.LogSeverityInfo)
	}

	s.entryMap = make(EntryMap)

	sourceInfo, err := os.Stat(s.config.SourcePath)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to obtain information about path %s", s.config.SourcePath), cli.LogSeverityError)
	}

	if sourceInfo.IsDir() {
		s.scanSourceDirectory(s.config.SourcePath, s.config.SourcePath)
	} else {
		s.scanSourceFile(filepath.Dir(s.config.SourcePath), s.config.SourcePath, sourceInfo)
	}

	for pathName, entry := range s.entryMap {
		if s.config.DebugMode {
			cli.LogWithPrefix(fmt.Sprintf("Scanned about %d bytes from %s", entry.Size, pathName), cli.LogSeverityWarn)
		}

		cli.Log(fmt.Sprintf("~ %s", pathName), cli.LogSeverityWarn)
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
		s.scanDirectoryEntry(relPath, entry)
	} else {
		s.scanFileEntry(relPath, entry)
	}
}

func (s *Scanner) scanDirectoryEntry(relPath string, entry util.DirectoryEntry) {
	if s.config.DebugMode {
		cli.LogWithPrefix(fmt.Sprintf("Scanning directory %s ...", relPath), cli.LogSeverityInfo)
	}

	s.entryMap[relPath] = NewDirectoryEntry(relPath, entry.Info)
}

func (s *Scanner) scanFileEntry(relPath string, entry util.DirectoryEntry) {
	fileName := entry.Path

	if s.config.DebugMode {
		cli.LogWithPrefix(fmt.Sprintf("Scanning file %s ...", relPath), cli.LogSeverityInfo)
	}

	fileData, err := os.ReadFile(fileName)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to obtain data from file %s", fileName), cli.LogSeverityError)
	}

	s.entryMap[relPath] = NewFileEntry(relPath, entry.Info, fileData)
}
