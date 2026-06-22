package util

import (
	"io/fs"
	"os"
	"path/filepath"
)

type DirectoryEntry struct {
	Path string
	Info os.FileInfo
}

func DirectoryWalk(root string) ([]DirectoryEntry, error) {
	var entries []DirectoryEntry

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

		entries = append(entries, DirectoryEntry{Path: path, Info: fileInfo})

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

	var result []DirectoryEntry

	for _, entry := range entries {
		if entry.Info.IsDir() && hasChildren[entry.Path] {
			continue
		}

		result = append(result, entry)
	}

	return result, nil
}
