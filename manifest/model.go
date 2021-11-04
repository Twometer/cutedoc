package manifest

type SourceManifest struct {
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

func NewDefaultSourceManifest() SourceManifest {
	return SourceManifest{
		ThemeId:    "default",
		InputPath:  "docs",
		OutputPath: "docs_gen",
	}
}

func NewDefaultThemeManifest() ThemeManifest {
	return ThemeManifest{}
}

func (manifest *SourceManifest) IsValid() bool {
	return manifest.Name != "" && manifest.ThemeId != "" && manifest.InputPath != "" && manifest.OutputPath != ""
}

func (manifest *ThemeManifest) IsValid() bool {
	return manifest.Name != "" && manifest.Version != ""
}
