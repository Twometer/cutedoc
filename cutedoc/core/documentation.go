package core

import (
	"cutedoc/diagnostics"
	"cutedoc/manifest"
	"cutedoc/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type navNode struct {
	Name     string
	Url      string
	Children []*navNode
}

type tocEntry struct {
	Id    string
	Name  string
	Level int
}

type generatorInfo struct {
	Name    string
	Version string
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
	Nav       []*navNode
	RootPath  string
	OutPath   string
}

func isIndexFile(filePath string) bool {
	return utils.GetFileName(filePath) == IndexFileName
}

func findDirForPage(page pageInfo, siteManifest manifest.SiteManifest) string {
	// Build the path for the output directory
	relativePath := page.FilePath[len(siteManifest.InputPath)+1:]
	outputDirPath := filepath.Dir(relativePath)

	// If it is not the index file, put it in its own subdirectory so that we get /subdir/index.html
	// which we can open in the browser as just /subdir
	if !isIndexFile(page.FilePath) {
		outputDirPath = filepath.Join(outputDirPath, page.FileName)
	}

	return outputDirPath
}

func createPageContext(mdFile string, rootPath string, siteManifest manifest.SiteManifest) (pageContext, error) {
	page, err := renderMarkdownPage(mdFile)
	if err != nil {
		return pageContext{}, err
	}

	return pageContext{
		Page: page,
		Site: siteManifest,
		Generator: generatorInfo{
			Name:    ProgramName,
			Version: ProgramVersion,
		},
		Now:      time.Now().Format("2006-01-02 15:04:05.000"),
		RootPath: rootPath,
		OutPath:  findDirForPage(page, siteManifest),
	}, nil
}

func openOutputFileForPage(pageContext *pageContext, siteManifest manifest.SiteManifest) (*os.File, error) {
	outPath := filepath.Join(siteManifest.OutputPath, pageContext.OutPath)

	// Create the output directory
	err := os.MkdirAll(outPath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// Create the file in the output directory
	file, err := os.Create(filepath.Join(outPath, "index.html"))
	return file, err
}

func generateThemedHtmlForPage(pageContext *pageContext, siteManifest manifest.SiteManifest, themeTemplate *template.Template) {
	mdFile := pageContext.Page.FileName
	writer, err := openOutputFileForPage(pageContext, siteManifest)
	if err != nil {
		diagnostics.PrintError(err, "failed to open output file for "+mdFile)
		return
	}

	htmlBuf := strings.Builder{}
	err = themeTemplate.ExecuteTemplate(&htmlBuf, RootTemplateName, pageContext)
	if err != nil {
		diagnostics.PrintError(err, "failed to execute template for "+mdFile)
		return
	}

	err = processHtml(strings.NewReader(htmlBuf.String()), writer, pageContext)
	if err != nil {
		diagnostics.PrintError(err, "failed to run HTML postproc for "+mdFile)
		return
	}

	err = writer.Close()
	if err != nil {
		diagnostics.PrintError(err, "failed to close output file for "+mdFile)
		return
	}
}

func prepareDocumentationTree(dirPath string, rootDirPrefix string, parentNode *navNode, siteManifest manifest.SiteManifest, contexts *[]pageContext) error {
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, dirent := range dir {
		childPath := filepath.Join(dirPath, dirent.Name())
		newNode := &navNode{
			Name: dirent.Name(),
			Url:  "",
		}

		if dirent.IsDir() {
			err := prepareDocumentationTree(childPath, rootDirPrefix+"../", newNode, siteManifest, contexts)
			if err != nil {
				return err
			}
		} else if filepath.Ext(dirent.Name()) == ".md" {
			diagnostics.Debug(func() {
				log.Println("processing:", childPath)
			})

			if !isIndexFile(childPath) {
				rootDirPrefix = "../" + rootDirPrefix
			}

			context, err := createPageContext(childPath, rootDirPrefix, siteManifest)
			if err != nil {
				return err
			}

			newNode.Name = context.Page.Title
			newNode.Url = filepath.ToSlash(context.OutPath)
			*contexts = append(*contexts, context)
		}

		parentNode.Children = append(parentNode.Children, newNode)
	}

	return nil
}

func copyMediaFiles(siteManifest manifest.SiteManifest, themeDir string) error {
	predicate := func(ext string) bool {
		return ext != ".md" && ext != ".html" && ext != ".ini"
	}

	err := utils.CopyDirContents(siteManifest.InputPath, siteManifest.OutputPath, predicate)
	if err != nil {
		return err
	}

	err = utils.CopyDirContents(themeDir, siteManifest.OutputPath, predicate)
	if err != nil {
		return err
	}

	return nil
}

func GenerateDocumentation(siteManifest manifest.SiteManifest, themeDir string) error {
	var stopwatch diagnostics.Stopwatch
	stopwatch.Reset()

	// Generate the documentation tree
	var navTreeRoot navNode
	var generatedPageContexts []pageContext
	err := prepareDocumentationTree(siteManifest.InputPath, "", &navTreeRoot, siteManifest, &generatedPageContexts)
	if err != nil {
		return err
	}

	// Generate the template for the theme
	themeTemplate, err := generateTemplate(themeDir)
	if err != nil {
		return err
	}

	// Apply the theme to each file and write it to disk
	for _, pageCtx := range generatedPageContexts {
		pageCtx.Nav = navTreeRoot.Children
		generateThemedHtmlForPage(&pageCtx, siteManifest, themeTemplate)
	}

	// Copy all the media (non-docs) files
	err = copyMediaFiles(siteManifest, themeDir)
	if err != nil {
		return err
	}

	log.Printf("done: generated in %d us", stopwatch.Microseconds())

	return nil
}
