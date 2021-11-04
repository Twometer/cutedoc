package generator

import (
	"html/template"
	"io/fs"
	"path/filepath"
)

func GenerateTemplate(themeDir string) (*template.Template, error) {
	var files []string
	err := filepath.WalkDir(themeDir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && filepath.Ext(path) == ".html" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
