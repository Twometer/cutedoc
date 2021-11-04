package generator

import (
	"cutedoc/manifest"
	"cutedoc/utils"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

const rootTemplate = "main.html"

type navNode struct {
	Name        string
	Url         string
	RelativeUrl string
	Children    []*navNode
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

func findDirForPage(page pageInfo, siteManifest manifest.SiteManifest) string {
	// Build the path for the output directory
	relativePath := page.FilePath[len(siteManifest.InputPath)+1:]
	outputDirPath := filepath.Dir(relativePath)

	// If it is not the index file, put it in its own subdirectory so that we get /subdir/index.html
	// which we can open in the browser as just /subdir
	if page.FileName != "index" {
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
			Name:    "Cutedoc",
			Version: "1.0.0",
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
	file, err := os.OpenFile(filepath.Join(outPath, "index.html"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	return file, err
}

func generateThemedHtmlForPage(pageContext *pageContext, siteManifest manifest.SiteManifest, themeTemplate *template.Template) {
	mdFile := pageContext.Page.FileName
	writer, err := openOutputFileForPage(pageContext, siteManifest)
	if err != nil {
		utils.PrintError(err, "failed to open output file for "+mdFile)
		return
	}

	err = themeTemplate.ExecuteTemplate(writer, rootTemplate, pageContext)
	if err != nil {
		utils.PrintError(err, "failed to execute template for "+mdFile)
		return
	}
}

func setRelativeUrl(navNodes []*navNode, prefix string) {
	for _, node := range navNodes {
		node.RelativeUrl = prefix + node.Url
		setRelativeUrl(node.Children, prefix)
	}
}

func GenerateDocumentation(siteManifest manifest.SiteManifest, themeManifest manifest.ThemeManifest, themeDir string) error {
	navTreeRoot := navNode{}

	var pageContexts []pageContext
	var walk func(path string, rootDirPrefix string, parentNode *navNode) error
	walk = func(path string, rootDirPrefix string, parentNode *navNode) error {
		dir, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		for _, dirent := range dir {
			childPath := filepath.Join(path, dirent.Name())
			newNode := &navNode{
				Name: dirent.Name(),
				Url:  "",
			}

			if dirent.IsDir() {
				err := walk(childPath, rootDirPrefix+"../", newNode)
				if err != nil {
					return err
				}
			} else if filepath.Ext(dirent.Name()) == ".md" {
				if dirent.Name() != "index.md" {
					rootDirPrefix = "../" + rootDirPrefix
				}
				context, err := createPageContext(childPath, rootDirPrefix, siteManifest)
				if err != nil {
					return err
				}
				newNode.Name = context.Page.Title
				newNode.Url = filepath.ToSlash(context.OutPath)
				pageContexts = append(pageContexts, context)
			}

			parentNode.Children = append(parentNode.Children, newNode)
		}

		return nil
	}

	err := walk(siteManifest.InputPath, "", &navTreeRoot)
	if err != nil {
		return err
	}

	themeTemplate, err := generateTemplate(themeDir)
	if err != nil {
		return err
	}

	for _, pageCtx := range pageContexts {
		pageCtx.Nav = navTreeRoot.Children
		setRelativeUrl(pageCtx.Nav, pageCtx.RootPath)
		generateThemedHtmlForPage(&pageCtx, siteManifest, themeTemplate)
	}

	return nil
}
