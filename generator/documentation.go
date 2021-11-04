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

type pageInfo struct {
	FilePath string
	FileName string
	Title    string
	Body     string
}

type pageContext struct {
	Page      pageInfo
	Generator generatorInfo
	Now       string
}

func newPageContext(info pageInfo) pageContext {
	return pageContext{
		Page: info,
		Generator: generatorInfo{
			Name:    "Cutedoc",
			Version: "1.0.0",
		},
		Now: time.Now().Format("2006-01-02 15:04:05.000"),
	}
}

func openFile(sourceManifest manifest.SourceManifest, info *pageInfo) (*os.File, error) {
	// Build the path for the output directory
	relativePath := info.FilePath[len(sourceManifest.InputPath)+1:]
	outputDirPath := filepath.Join(sourceManifest.OutputPath, filepath.Dir(relativePath))

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

func generateFile(wg *sync.WaitGroup, file string, sourceManifest manifest.SourceManifest, themeTemplate *template.Template) {
	defer wg.Done()

	pageInfo, err := createPage(file)
	if err != nil {
		utils.PrintError(err, "failed to generate html for "+file)
		return
	}

	writer, err := openFile(sourceManifest, &pageInfo)
	if err != nil {
		utils.PrintError(err, "failed to open output file for "+file)
		return
	}

	pageContext := newPageContext(pageInfo)
	err = themeTemplate.ExecuteTemplate(writer, rootTemplate, pageContext)
	if err != nil {
		utils.PrintError(err, "failed to run template for "+file)
		return
	}
}

func GenerateDocumentation(sourceManifest manifest.SourceManifest, themeManifest manifest.ThemeManifest, themeDir string) {
	log.Println("Using theme", themeManifest.Name)
	themeTemplate, err := generateTemplate(themeDir)
	utils.HandleError(err)

	files, err := utils.ScanDir(sourceManifest.InputPath, ".md")
	utils.HandleError(err)

	var stopwatch utils.Stopwatch
	var wg sync.WaitGroup

	stopwatch.Reset()
	for _, file := range files {
		wg.Add(1)
		go generateFile(&wg, file, sourceManifest, themeTemplate)
	}
	wg.Wait()
	log.Printf("Generator completed in %d us.\n", stopwatch.Microseconds())
}
