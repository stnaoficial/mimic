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

	return entries, err
}
