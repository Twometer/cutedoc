package model

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