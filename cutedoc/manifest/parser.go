package manifest

import (
	"errors"
	"gopkg.in/ini.v1"
	"path/filepath"
)

func ParseSiteManifest(path string) (SiteManifest, error) {
	manifest, err := ini.Load(path)
	if err != nil {
		return SiteManifest{}, err
	}

	result := SiteManifest{}

	pageSection := manifest.Section("Page")
	if pageSection != nil {
		result.Name = pageSection.Key("Name").String()
		result.ThemeId = pageSection.Key("Theme").MustString("default")
	}

	filesSection := manifest.Section("Files")
	if filesSection != nil {
		result.InputPath = pageSection.Key("Input").MustString("docs")
		result.OutputPath = pageSection.Key("Output").MustString("docs_gen")
	}

	if !result.IsValid() {
		return result, errors.New("missing required parameters")
	}

	result.InputPath = filepath.Clean(result.InputPath)
	result.OutputPath = filepath.Clean(result.OutputPath)

	return result, nil
}

func ParseThemeManifest(path string) (ThemeManifest, error) {
	manifest, err := ini.Load(path)
	if err != nil {
		return ThemeManifest{}, err
	}

	result := ThemeManifest{}

	rootSection := manifest.Section("Theme")
	if rootSection != nil {
		result.Name = rootSection.Key("Name").String()
		result.Description = rootSection.Key("Description").String()
		result.Repository = rootSection.Key("Repository").String()
		result.Version = rootSection.Key("Version").String()
		result.Author = rootSection.Key("Author").String()
		result.License = rootSection.Key("License").String()
	}

	highlightingSection := manifest.Section("Highlighting")
	if highlightingSection != nil {
		result.Highlighting.Style = highlightingSection.Key("Style").MustString("bw")
		result.Highlighting.LineNumbers = highlightingSection.Key("LineNumbers").MustBool(false)
	}

	if !result.IsValid() {
		return result, errors.New("missing required parameters")
	}

	return result, nil
}
