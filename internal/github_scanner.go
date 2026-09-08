package internal

import (
	"fmt"
	"mimic/internal/cli"
	"mimic/internal/util"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type GitHubScanner struct {
	util util.GitHubUtility

	debug bool
}

func NewGitHubScanner(userName string, repoName string, debug bool) *GitHubScanner {
	return &GitHubScanner{
		util: *util.NewGitHubUtility(userName, repoName),

		debug: debug,
	}
}

func (g *GitHubScanner) Scan(sourceNames []string) EntryMap {
	entryMap := make(EntryMap)

	for _, sourceName := range sourceNames {
		sourcePath := ".mimic/" + sourceName

		if g.debug {
			cli.Logf(cli.LogSeverityWarn, "Scanning source path %s ...\n", sourcePath)
		}

		g.scanSourceDirectory(entryMap, sourcePath, sourcePath)
	}

	if g.debug {
		for pathName, entry := range entryMap {
			cli.Logf(cli.LogSeverityInfo, "Scanned about %d bytes from %s\n", entry.Size, pathName)
		}
	}

	return entryMap
}

func (g *GitHubScanner) scanSourceDirectory(entryMap EntryMap, basePath string, dirName string) {
	entries, err := g.util.Walk(dirName)

	if err != nil {
		cli.Logf(cli.LogSeverityError, "Unable to walk into source directory %s\n", dirName)
		os.Exit(1)
	}

	if len(entries) == 0 {
		cli.Logf(cli.LogSeverityError, "No entries found in source directory %s\n", dirName)
		os.Exit(1)
	}

	for _, entry := range entries {
		g.scanEntry(entryMap, basePath, entry)
	}
}

func (g *GitHubScanner) scanSourceFile(entryMap EntryMap, basePath string, fileName string, fileInfo os.FileInfo) {
	g.scanEntry(entryMap, basePath, util.GitHubRemoteEntry{Path: fileName, Info: fileInfo})
}

func (g *GitHubScanner) scanEntry(entryMap EntryMap, basePath string, entry util.GitHubRemoteEntry) {
	if basePath == entry.Path {
		return
	}

	relPath, err := filepath.Rel(basePath, entry.Path)

	if err != nil {
		cli.Logf(cli.LogSeverityError, "Unable to determine relative path for %s\n", entry.Path)
		os.Exit(1)
	}

	if entry.Info.IsDir() {
		g.scanDirectoryEntry(entryMap, relPath, entry)
	} else {
		g.scanFileEntry(entryMap, relPath, entry)
	}
}

func (g *GitHubScanner) scanDirectoryEntry(entryMap EntryMap, relPath string, entry util.GitHubRemoteEntry) {
	if g.debug {
		cli.Logf(cli.LogSeverityWarn, "Scanning directory %s ...\n", relPath)
	}

	entryMap[relPath] = NewDirectoryEntry(relPath, entry.Info)
}

func (g *GitHubScanner) scanFileEntry(entryMap EntryMap, relPath string, entry util.GitHubRemoteEntry) {
	fileName := entry.Path

	if g.debug {
		cli.Logf(cli.LogSeverityWarn, "Scanning file %s ...\n", relPath)
	}

	fileData, err := g.util.FetchRawContent(fileName)

	if err != nil {
		cli.Logf(cli.LogSeverityError, "Unable to obtain data from file %s\n", fileName)
		os.Exit(1)
	}

	entryMap[relPath] = NewFileEntry(relPath, entry.Info, fileData)
}

func (g *GitHubScanner) List() {
	entries, err := g.util.Walk(".mimic")

	if err != nil {
		cli.Logln(cli.LogSeverityError, "Unable to walk into source directory")
		os.Exit(1)
	}

	dirNames := []string{}

	for _, entry := range entries {
		parts := strings.Split(entry.Path, "/")
		dirName := parts[1]

		if dirName == "mimic" {
			dirName = parts[1] + "/" + parts[2]
		}

		if slices.Contains(dirNames, dirName) {
			continue
		}

		fmt.Printf("%s\n", dirName)

		dirNames = append(dirNames, dirName)
	}

	fmt.Println()
}
