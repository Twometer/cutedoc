package main

import (
	"cutedoc/generator"
	"cutedoc/manifest"
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"
)

const themeManifestName = "theme.ini"
const sourceManifestName = "cutedoc.ini"

func findThemesDirectory() (string, error) {
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
	config, err := manifest.ParseSourceManifest(sourceManifestName)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(config.OutputPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	themesDir, err := findThemesDirectory()
	if err != nil {
		panic(err)
	}

	theme, themeDir, err := findThemeConfig(themesDir, config.ThemeId)
	if err != nil {
		panic(err)
	}

	themeTemplate, err := generator.GenerateTemplate(themeDir)
	if err != nil {
		panic(err)
	}

	log.Println(themesDir)
	log.Println("Processing", config.InputPath)
	log.Println("Using theme", theme.Name)
	log.Println("in directory", themeDir)
}
