package main

import (
	"cutedoc/generator"
	"cutedoc/manifest"
	"cutedoc/utils"
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

type Page struct {
	Title string
	Body  string
}

func main() {
	config, err := manifest.ParseSourceManifest(sourceManifestName)
	utils.HandleError(err)

	err = os.MkdirAll(config.OutputPath, os.ModePerm)
	utils.HandleError(err)

	themesDir, err := findThemesDirectory()
	utils.HandleError(err)

	theme, themeDir, err := findThemeConfig(themesDir, config.ThemeId)
	utils.HandleError(err)

	log.Println("Using theme", theme.Name)
	themeTemplate, err := generator.GenerateTemplate(themeDir)
	utils.HandleError(err)

	files, err := utils.ScanDir(config.InputPath, ".md")
	utils.HandleError(err)

	for _, file := range files {
		html, err := generator.GenerateHtml(file)
		if err != nil {
			utils.PrintError(err, "failed to generate html for "+file)
			continue
		}

		writer, err := os.OpenFile(path.Join(config.OutputPath, "output.html"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			utils.PrintError(err, "failed to open output file for "+file)
			continue
		}

		err = themeTemplate.ExecuteTemplate(writer, "main.html", Page{
			Title: "Test",
			Body:  html,
		})
		if err != nil {
			utils.PrintError(err, "failed to run template for "+file)
			continue
		}
	}
}
