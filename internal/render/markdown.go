// Package render converts Markdown content to HTML using goldmark.
// This package only handles string→string conversion;
// it does not depend on domain models to avoid circular imports.
package render

import (
	"bytes"
	"fmt"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// MarkdownRenderer converts Markdown to HTML.
type MarkdownRenderer struct {
	md goldmark.Markdown
}

// NewMarkdownRenderer creates a renderer with GFM extensions and code highlighting.
func NewMarkdownRenderer() *MarkdownRenderer {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(false),
					chromahtml.TabWidth(4),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)
	return &MarkdownRenderer{md: md}
}

// Render converts raw Markdown to HTML.
func (r *MarkdownRenderer) Render(markdown string) (string, error) {
	var buf bytes.Buffer
	if err := r.md.Convert([]byte(markdown), &buf); err != nil {
		return "", fmt.Errorf("render markdown: %w", err)
	}
	return buf.String(), nil
}
