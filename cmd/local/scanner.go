package local

import (
	"mimic/internal"
	"mimic/internal/cli"
	"mimic/internal/util"
	"os"
	"path/filepath"
)

type Scanner struct {
	scannedEntries internal.FileSystemMap

	debug bool
}

func NewScanner(debug bool) *Scanner {
	return &Scanner{
		scannedEntries: make(internal.FileSystemMap),

		debug: debug,
	}
}

func (s *Scanner) Scan(sourcePaths []string) (internal.FileSystemMap, error) {
	s.scannedEntries = make(internal.FileSystemMap)

	for _, sourcePath := range sourcePaths {
		// allow debug
		if s.debug {
			cli.Logf(cli.LogSeverityWarn, "Scanning source path %s ...\n", sourcePath)
		}

		sourceInfo, err := os.Stat(sourcePath)

		if err != nil {
			return nil, err
		}

		if sourceInfo.IsDir() {
			if err := s.scanSourceDirectory(sourcePath, sourcePath); err != nil {
				return nil, err
			}
		} else if err := s.scanSourceFile(filepath.Dir(sourcePath), sourcePath, sourceInfo); err != nil {
			return nil, err
		}
	}

	for pathName, entry := range s.scannedEntries {
		// allow debug
		if s.debug {
			cli.Logf(cli.LogSeverityInfo, "Scanned about %d bytes from %s\n", entry.Size, pathName)
		}
	}

	return s.scannedEntries, nil
}

func (s *Scanner) scanSourceDirectory(basePath string, dirName string) error {
	entriesFound, err := util.FileSystemWalk(dirName)

	if err != nil {
		return err
	}

	for _, entry := range entriesFound {
		s.scanEntry(basePath, entry)
	}

	return nil
}

func (s *Scanner) scanSourceFile(basePath string, fileName string, fileInfo os.FileInfo) error {
	return s.scanEntry(basePath, util.FileSystemEntry{Path: fileName, Info: fileInfo})
}

func (s *Scanner) scanEntry(basePath string, entry util.FileSystemEntry) error {
	if basePath == entry.Path {
		return nil
	}

	relPath, err := filepath.Rel(basePath, entry.Path)

	if err != nil {
		return err
	}

	if entry.Info.IsDir() {
		s.scanDirectoryEntry(relPath, entry)
	} else if err := s.scanFileEntry(relPath, entry); err != nil {
		return err
	}

	return nil
}

func (s *Scanner) scanDirectoryEntry(relPath string, entry util.FileSystemEntry) {
	// allow debug
	if s.debug {
		cli.Logf(cli.LogSeverityWarn, "Scanning directory %s ...\n", relPath)
	}

	s.scannedEntries[relPath] = internal.NewDirectoryEntry(relPath, entry.Info)
}

func (s *Scanner) scanFileEntry(relPath string, entry util.FileSystemEntry) error {
	// allow debug
	if s.debug {
		cli.Logf(cli.LogSeverityWarn, "Scanning file %s ...\n", relPath)
	}

	fileName := entry.Path

	fileData, err := os.ReadFile(fileName)

	if err != nil {
		return err
	}

	s.scannedEntries[relPath] = internal.NewFileEntry(relPath, entry.Info, fileData)

	return nil
}
