package manifest

type SiteManifest struct {
	Name       string
	ThemeId    string
	InputPath  string
	OutputPath string
}

type ThemeManifest struct {
	Name        string
	Description string
	Repository  string
	Version     string
	Author      string
	License     string
}

func NewDefaultSiteManifest() SiteManifest {
	return SiteManifest{
		ThemeId:    "default",
		InputPath:  "docs",
		OutputPath: "docs_gen",
	}
}

func NewDefaultThemeManifest() ThemeManifest {
	return ThemeManifest{}
}

func (manifest *SiteManifest) IsValid() bool {
	return manifest.Name != "" && manifest.ThemeId != "" && manifest.InputPath != "" && manifest.OutputPath != ""
}

func (manifest *ThemeManifest) IsValid() bool {
	return manifest.Name != "" && manifest.Version != ""
}
