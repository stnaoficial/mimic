package util

import (
	"io/fs"
	"path/filepath"
)

func DirectoryWalk(root string) ([]string, error) {
	filenames := []string{}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		filenames = append(filenames, path)

		return nil
	})

	return filenames, err
}
