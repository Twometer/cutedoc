package generator

import (
	"bytes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"os"
	"path/filepath"
	"strings"
)

var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
)

func scanTree(node ast.Node, consumer func(node ast.Node)) {
	consumer(node)
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		scanTree(child, consumer)
	}
}

func analyzeDocument(astRoot ast.Node, source []byte, pageInfo *pageInfo) {
	scanTree(astRoot, func(node ast.Node) {
		switch node.Kind() {
		case ast.KindHeading:
			heading := node.(*ast.Heading)
			if heading.Level == 1 && pageInfo.Title == "" {
				pageInfo.Title = string(heading.Text(source))
			}
		}
	})
}

func getFileName(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

func createPage(mdFile string) (pageInfo, error) {
	result := pageInfo{
		FilePath: filepath.Clean(mdFile),
		FileName: getFileName(mdFile),
	}

	source, err := os.ReadFile(mdFile)
	if err != nil {
		return result, err
	}

	// Parse Markdown
	reader := text.NewReader(source)
	astRoot := md.Parser().Parse(reader)
	analyzeDocument(astRoot, source, &result)

	// Render to HTML
	var buf bytes.Buffer
	err = md.Renderer().Render(&buf, source, astRoot)
	if err == nil {
		result.Body = buf.String()
	}

	return result, err
}
