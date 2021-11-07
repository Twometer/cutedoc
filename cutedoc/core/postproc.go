package core

import (
	"cutedoc/utils"
	"golang.org/x/net/html"
	"io"
	"strings"
)

func shouldChangeLink(link string) bool {
	if len(link) == 0 {
		return false
	}

	isHeading := link[0] == '#'
	isExternal := strings.Contains(link, "://")
	return !isHeading && !isExternal
}

func processHtmlNode(node *html.Node, pageContext *pageContext) {
	for idx, attr := range node.Attr {
		if (attr.Key == "href" || attr.Key == "src") && shouldChangeLink(attr.Val) {
			path := utils.StripParentDirectories(attr.Val)
			path = pageContext.RootPath + path
			node.Attr[idx].Val = path
		}
	}

	for n := node.FirstChild; n != nil; n = n.NextSibling {
		processHtmlNode(n, pageContext)
	}
}

func processHtml(input io.Reader, output io.Writer, pageContext *pageContext) error {
	htmlRootNode, err := html.Parse(input)
	if err != nil {
		return err
	}

	processHtmlNode(htmlRootNode, pageContext)

	return html.Render(output, htmlRootNode)
}
