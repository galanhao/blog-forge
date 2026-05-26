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
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
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

// RenderWithTOC converts Markdown and extracts the TOC separately.
func (r *MarkdownRenderer) RenderWithTOC(markdown string) (contentHTML string, tocHTML string, err error) {
	source := []byte(markdown)

	// Parse the document
	ctx := parser.NewContext()
	doc := r.md.Parser().Parse(text.NewReader(source), parser.WithContext(ctx))

	// Render body
	var bodyBuf bytes.Buffer
	if err := r.md.Renderer().Render(&bodyBuf, source, doc); err != nil {
		return "", "", fmt.Errorf("render body: %w", err)
	}

	// Walk AST to build TOC
	toc := buildTOC(doc, source)

	return bodyBuf.String(), toc, nil
}

// tocEntry is a single TOC item.
type tocEntry struct {
	Level int
	ID    string
	Text  string
}

// buildTOC walks the AST and builds an HTML table of contents.
func buildTOC(doc ast.Node, source []byte) string {
	var entries []tocEntry
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		heading, ok := n.(*ast.Heading)
		if !ok {
			return ast.WalkContinue, nil
		}
		if heading.Level < 1 || heading.Level > 4 {
			return ast.WalkContinue, nil
		}
		id, _ := heading.AttributeString("id")
		idStr := ""
		if id != nil {
			idStr = string(id.([]byte))
		}
		text := extractText(heading, source)
		entries = append(entries, tocEntry{
			Level: heading.Level,
			ID:    idStr,
			Text:  text,
		})
		return ast.WalkContinue, nil
	})

	if len(entries) == 0 {
		return ""
	}

	var buf bytes.Buffer
	buf.WriteString("<ul>\n")
	for _, e := range entries {
		buf.WriteString(fmt.Sprintf(`<li><a href="#%s">%s</a></li>`+"\n", e.ID, e.Text))
	}
	buf.WriteString("</ul>")
	return buf.String()
}

// extractText gets the plain text content of a node.
func extractText(n ast.Node, source []byte) string {
	var buf bytes.Buffer
	ast.Walk(n, func(child ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if t, ok := child.(*ast.Text); ok {
			seg := t.Segment
			buf.Write(seg.Value(source))
		}
		return ast.WalkContinue, nil
	})
	return buf.String()
}

// Ensure util package is referenced (needed by goldmark extension).
var _ = util.IsSpace
