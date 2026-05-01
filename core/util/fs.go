package util

import (
	"io/fs"
	"path/filepath"
)

type FileMap = map[string]string

func DirectoryWalk(root string) ([]string, error) {
	fileNames := []string{}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		fileNames = append(fileNames, path)

		return nil
	})

	return fileNames, err
}
