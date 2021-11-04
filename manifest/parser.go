package manifest

import (
	"errors"
	"gopkg.in/ini.v1"
	"path/filepath"
)

func ParseSourceManifest(path string) (SourceManifest, error) {
	manifest, err := ini.Load(path)
	if err != nil {
		return SourceManifest{}, err
	}

	result := NewDefaultSourceManifest()

	pageSection := manifest.Section("Page")
	if pageSection != nil {
		result.Name = pageSection.Key("Name").String()
		result.ThemeId = pageSection.Key("Theme").MustString(result.ThemeId)
	}

	filesSection := manifest.Section("Files")
	if filesSection != nil {
		result.InputPath = pageSection.Key("Input").MustString(result.InputPath)
		result.OutputPath = pageSection.Key("Output").MustString(result.OutputPath)
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

	result := NewDefaultThemeManifest()

	themeSection := manifest.Section("Theme")
	if themeSection != nil {
		result.Name = themeSection.Key("Name").String()
		result.Description = themeSection.Key("Description").String()
		result.Repository = themeSection.Key("Repository").String()
		result.Version = themeSection.Key("Version").String()
		result.Author = themeSection.Key("Author").String()
		result.License = themeSection.Key("License").String()
	}

	if !result.IsValid() {
		return result, errors.New("missing required parameters")
	}

	return result, nil
}
