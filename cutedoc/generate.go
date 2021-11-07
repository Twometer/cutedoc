package main

import (
	"cutedoc/generator"
	"cutedoc/manifest"
	"cutedoc/utils"
	"log"
	"os"
)

func runGenerator() error {
	siteManifest, err := manifest.ParseSiteManifest(utils.SiteManifestName)
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

	log.Printf("using theme: '%s' by %s", themeManifest.Name, themeManifest.Author)
	return generator.GenerateDocumentation(siteManifest, themeDir)
}
