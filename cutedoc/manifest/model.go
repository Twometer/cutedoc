package manifest

type SiteManifest struct {
	Name       string
	ThemeId    string
	InputPath  string
	OutputPath string
}

type HighlightingConfig struct {
	Style       string
	LineNumbers bool
}

type ThemeManifest struct {
	Name         string
	Description  string
	Repository   string
	Version      string
	Author       string
	License      string
	Highlighting HighlightingConfig
}

func (manifest *SiteManifest) IsValid() bool {
	return manifest.Name != "" && manifest.ThemeId != "" && manifest.InputPath != "" && manifest.OutputPath != ""
}

func (manifest *ThemeManifest) IsValid() bool {
	return manifest.Name != "" && manifest.Version != "" && manifest.Highlighting.Style != ""
}
