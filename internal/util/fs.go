package util

import (
	"io/fs"
	"os"
	"path/filepath"
)

type FileSystemEntry struct {
	Path string
	Info os.FileInfo
}

func FileSystemWalk(root string) ([]FileSystemEntry, error) {
	var entries []FileSystemEntry

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == root {
			return nil
		}

		fileInfo, err := d.Info()

		if err != nil {
			return err
		}

		entries = append(entries, FileSystemEntry{Path: path, Info: fileInfo})

		return nil
	})

	if err != nil {
		return nil, err
	}

	hasChildren := make(map[string]bool)

	for _, entry := range entries {
		parent := filepath.Dir(entry.Path)

		for parent != root && parent != "." {
			hasChildren[parent] = true
			parent = filepath.Dir(parent)
		}

		hasChildren[parent] = true
	}

	var result []FileSystemEntry

	for _, entry := range entries {
		if entry.Info.IsDir() && hasChildren[entry.Path] {
			continue
		}

		result = append(result, entry)
	}

	return result, nil
}
