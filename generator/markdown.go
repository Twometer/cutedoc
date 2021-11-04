package generator

import (
	"bytes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"os"
)

var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
)

func GenerateHtml(mdFile string) (string, error) {
	markdown, err := os.ReadFile(mdFile)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = md.Convert(markdown, &buf)
	if err != nil {
		return "", err
	} else {
		return buf.String(), err
	}
}
