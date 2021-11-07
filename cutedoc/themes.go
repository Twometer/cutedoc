package main

import (
	"cutedoc/core"
	"cutedoc/manifest"
	"errors"
	"os"
	"path"
	"path/filepath"
)

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

	themeManifest, err := manifest.ParseThemeManifest(path.Join(themeDir, core.ThemeManifestName))
	return themeManifest, themeDir, err
}
