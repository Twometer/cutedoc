package generator

import (
	"cutedoc/manifest"
	"cutedoc/utils"
	"log"
	"os"
	"path/filepath"
	"sync"
	"text/template"
	"time"
)

const rootTemplate = "main.html"

type generatorInfo struct {
	Name    string
	Version string
}

type tocEntry struct {
	Id    string
	Name  string
	Level int
}

type pageInfo struct {
	FilePath string
	FileName string
	Title    string
	Body     string
	Toc      []tocEntry
}

type pageContext struct {
	Page      pageInfo
	Generator generatorInfo
	Now       string
	Site      manifest.SiteManifest
}

func newPageContext(page pageInfo, site manifest.SiteManifest) pageContext {
	return pageContext{
		Page: page,
		Site: site,
		Generator: generatorInfo{
			Name:    "Cutedoc",
			Version: "1.0.0",
		},
		Now: time.Now().Format("2006-01-02 15:04:05.000"),
	}
}

func openFile(siteManifest manifest.SiteManifest, info *pageInfo) (*os.File, error) {
	// Build the path for the output directory
	relativePath := info.FilePath[len(siteManifest.InputPath)+1:]
	outputDirPath := filepath.Join(siteManifest.OutputPath, filepath.Dir(relativePath))

	// If it is not the index file, put it in its own subdirectory so that we get /subdir/index.html
	// which we can open in the browser as just /subdir
	if info.FileName != "index" {
		outputDirPath = filepath.Join(outputDirPath, info.FileName)
	}

	// Create the output directory
	err := os.MkdirAll(outputDirPath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// Create the file in the output directory
	file, err := os.OpenFile(filepath.Join(outputDirPath, "index.html"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	return file, err
}

func generateFile(wg *sync.WaitGroup, file string, siteManifest manifest.SiteManifest, themeTemplate *template.Template) {
	defer wg.Done()

	pageInfo, err := createPage(file)
	if err != nil {
		utils.PrintError(err, "failed to generate html for "+file)
		return
	}

	writer, err := openFile(siteManifest, &pageInfo)
	if err != nil {
		utils.PrintError(err, "failed to open output file for "+file)
		return
	}

	pageContext := newPageContext(pageInfo, siteManifest)
	err = themeTemplate.ExecuteTemplate(writer, rootTemplate, pageContext)
	if err != nil {
		utils.PrintError(err, "failed to run template for "+file)
		return
	}
}

func GenerateDocumentation(siteManifest manifest.SiteManifest, themeManifest manifest.ThemeManifest, themeDir string) error {
	log.Println("Using theme", themeManifest.Name)
	themeTemplate, err := generateTemplate(themeDir)
	if err != nil {
		return err
	}

	files, err := utils.ScanDir(siteManifest.InputPath, ".md")
	if err != nil {
		return err
	}

	var stopwatch utils.Stopwatch
	var wg sync.WaitGroup

	stopwatch.Reset()
	for _, file := range files {
		wg.Add(1)
		go generateFile(&wg, file, siteManifest, themeTemplate)
	}
	wg.Wait()
	log.Printf("Generator completed in %d us.\n", stopwatch.Microseconds())

	return nil
}
