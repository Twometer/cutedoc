package generator

import (
	"github.com/iancoleman/strcase"
	"path/filepath"
	"strings"
)

func getFileName(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

func prettyTitle(filePath string) string {
	fileName := getFileName(filePath)
	return strcase.ToDelimited(fileName, ' ')
}
