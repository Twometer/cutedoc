package main

import (
	"cutedoc/generator"
	"cutedoc/manifest"
	"os"
)

func runGenerator() error {
	siteManifest, err := manifest.ParseSiteManifest(siteManifestName)
	if err != nil {
		return err
	}

	err = os.MkdirAll(siteManifest.OutputPath, os.ModePerm)
	if err != nil {
		return err
	}

	themeBaseDir, err := findThemesBaseDir()
	if err != nil {
		return err
	}

	themeManifest, themeDir, err := findThemeConfig(themeBaseDir, siteManifest.ThemeId)
	if err != nil {
		return err
	}

	return generator.GenerateDocumentation(siteManifest, themeManifest, themeDir)
}
