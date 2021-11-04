package main

import (
	"cutedoc/generator"
	"cutedoc/manifest"
	"os"
)

func runGenerator() error {
	sourceManifest, err := manifest.ParseSourceManifest(sourceManifestName)
	if err != nil {
		return err
	}

	err = os.MkdirAll(sourceManifest.OutputPath, os.ModePerm)
	if err != nil {
		return err
	}

	themeBaseDir, err := findThemesBaseDir()
	if err != nil {
		return err
	}

	themeManifest, themeDir, err := findThemeConfig(themeBaseDir, sourceManifest.ThemeId)
	if err != nil {
		return err
	}

	return generator.GenerateDocumentation(sourceManifest, themeManifest, themeDir)
}
