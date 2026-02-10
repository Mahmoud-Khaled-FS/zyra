package utils

import (
	"io/fs"
	"path/filepath"
)

func ReadDirR(path string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return files, err
	}
	return files, nil
}
