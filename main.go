package main

import (
	"cutedoc/generator"
	"cutedoc/manifest"
	"cutedoc/utils"
	"errors"
	"os"
	"path"
	"path/filepath"
)

const themeManifestName = "theme.ini"
const sourceManifestName = "cutedoc.ini"

func findThemesBaseDir() (string, error) {
	envConfig := os.Getenv("CUTEDOC_THEME_DIR")
	if envConfig != "" {
		return envConfig, nil
	}

	executable, err := os.Executable()
	if err != nil {
		return "", err
	}

	baseDir := filepath.Dir(executable)
	return path.Join(baseDir, "themes"), nil
}

func findThemeConfig(themesDir string, themeId string) (manifest.ThemeManifest, string, error) {
	themeDir := path.Join(themesDir, themeId)
	themeDirStat, err := os.Stat(themeDir)
	if err != nil {
		return manifest.ThemeManifest{}, "", err
	}

	if !themeDirStat.IsDir() {
		return manifest.ThemeManifest{}, "", errors.New("theme is not a directory")
	}

	themeManifest, err := manifest.ParseThemeManifest(path.Join(themeDir, themeManifestName))
	return themeManifest, themeDir, err
}

func main() {
	sourceManifest, err := manifest.ParseSourceManifest(sourceManifestName)
	utils.HandleError(err)

	err = os.MkdirAll(sourceManifest.OutputPath, os.ModePerm)
	utils.HandleError(err)

	themeBaseDir, err := findThemesBaseDir()
	utils.HandleError(err)

	themeManifest, themeDir, err := findThemeConfig(themeBaseDir, sourceManifest.ThemeId)
	utils.HandleError(err)

	generator.GenerateDocumentation(sourceManifest, themeManifest, themeDir)
}
