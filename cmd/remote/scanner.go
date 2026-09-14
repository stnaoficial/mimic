package remote

import (
	"mimic/internal"
	"mimic/internal/cli"
	"mimic/tools/github"
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

		if err := s.scanSourceDirectory(sourcePath, sourcePath); err != nil {
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
	entriesFound, err := GithubApiClient.Walk(dirName)

	if err != nil {
		return err
	}

	for _, entry := range entriesFound {
		if err := s.scanEntry(basePath, entry); err != nil {
			return err
		}
	}

	return nil
}

func (s *Scanner) scanEntry(basePath string, entry github.ApiTreeEntry) error {
	if basePath == entry.Path {
		return nil
	}

	relPath, err := filepath.Rel(basePath, entry.Path)

	if err != nil {
		return err
	}

	if entry.Type == github.ApiTreeEntryTypeTree {
		s.scanDirectoryEntry(relPath, entry)
	} else if err := s.scanFileEntry(relPath, entry); err != nil {
		return err
	}

	return nil
}

func (s *Scanner) scanDirectoryEntry(relPath string, entry github.ApiTreeEntry) {
	// allow debug
	if s.debug {
		cli.Logf(cli.LogSeverityWarn, "Scanning directory %s ...\n", relPath)
	}

	remoteEntry := GithubApiClient.ParseRemoteEntry(entry)

	s.scannedEntries[relPath] = internal.NewDirectoryEntry(relPath, remoteEntry.Info)
}

func (s *Scanner) scanFileEntry(relPath string, entry github.ApiTreeEntry) error {
	// allow debug
	if s.debug {
		cli.Logf(cli.LogSeverityWarn, "Scanning file %s ...\n", relPath)
	}

	blobEntry, err := GithubApiClient.FetchApiBlob(entry.Sha)

	if err != nil {
		return err
	}

	entryData, err := blobEntry.DecodeContent()

	if err != nil {
		return err
	}

	remoteEntry := GithubApiClient.ParseRemoteEntry(entry)

	s.scannedEntries[relPath] = internal.NewFileEntry(relPath, remoteEntry.Info, entryData)

	return nil
}
