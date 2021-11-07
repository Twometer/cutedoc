package utils

import (
	"github.com/iancoleman/strcase"
	"path/filepath"
	"regexp"
	"strings"
)

var pathRegex *regexp.Regexp

func init() {
	regex, err := regexp.Compile("^(\\.\\./)*")
	HandleError(err)

	pathRegex = regex
}

func GetFileName(path string) string {
	cleanPath := filepath.Clean(path)
	return strings.TrimSuffix(filepath.Base(cleanPath), filepath.Ext(cleanPath))
}

func PrettifyTitle(path string) string {
	fileName := GetFileName(path)
	return strcase.ToDelimited(fileName, ' ')
}

func StripParentDirectories(path string) string {
	cleanPath := filepath.ToSlash(filepath.Clean(path))
	return pathRegex.ReplaceAllString(cleanPath, "")
}
