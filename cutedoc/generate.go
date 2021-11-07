package main

import (
	"cutedoc/core"
	"cutedoc/diagnostics"
	"cutedoc/manifest"
	"log"
	"os"
)

func runGenerator() error {
	siteManifest, err := manifest.ParseSiteManifest(core.SiteManifestName)
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

	diagnostics.Debug(func() {
		log.Printf("using theme: '%s' by %s", themeManifest.Name, themeManifest.Author)
	})
	return core.GenerateDocumentation(siteManifest, themeManifest, themeDir)
}
