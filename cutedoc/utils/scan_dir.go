package utils

import (
	"io/fs"
	"path/filepath"
)

func ScanDir(dir string, extension string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && filepath.Ext(path) == extension {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
