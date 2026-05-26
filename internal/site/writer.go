package site

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/galanhao/blog-forge/internal/content"
	"github.com/galanhao/blog-forge/internal/theme"
)

// writePosts generates individual HTML pages for each post.
func (b *Builder) writePosts(distDir string, siteCtx theme.SiteCtx, themeCtx theme.ThemeCtx, posts []*content.Post) error {
	for i, p := range posts {
		postCtx := buildPostCtx(siteCtx, themeCtx, p)

		if i > 0 {
			prevCtx := buildPostCtx(siteCtx, themeCtx, posts[i-1])
			postCtx.Prev = &prevCtx
		}
		if i < len(posts)-1 {
			nextCtx := buildPostCtx(siteCtx, themeCtx, posts[i+1])
			postCtx.Next = &nextCtx
		}

		html, err := b.engine.Execute("post.html", postCtx)
		if err != nil {
			return fmt.Errorf("render post %s: %w", p.Slug, err)
		}

		outPath := filepath.Join(distDir, permalinkToPath(p.Permalink), "index.html")
		if err := writeFile(outPath, html); err != nil {
			return err
		}
	}
	return nil
}

// writePages generates individual HTML pages for standalone pages.
func (b *Builder) writePages(distDir string, siteCtx theme.SiteCtx, themeCtx theme.ThemeCtx, pages []*content.Post) error {
	for _, p := range pages {
		pageCtx := buildPostCtx(siteCtx, themeCtx, p)
		pageCtx.Layout = content.LayoutPage

		tplName := "page.html"
		if p.Layout != "" && p.Layout != content.LayoutPage {
			tplName = p.Layout + ".html"
		}

		html, err := b.engine.Execute(tplName, pageCtx)
		if err != nil {
			return fmt.Errorf("render page %s: %w", p.Slug, err)
		}

		outPath := filepath.Join(distDir, permalinkToPath(p.Permalink), "index.html")
		if err := writeFile(outPath, html); err != nil {
			return err
		}
	}
	return nil
}

// writeIndex generates paginated index pages.
func (b *Builder) writeIndex(distDir string, siteCtx theme.SiteCtx, themeCtx theme.ThemeCtx, posts []*content.Post) error {
	perPage := b.cfg.PerPage
	if perPage <= 0 {
		perPage = 10
	}

	totalPages := (len(posts) + perPage - 1) / perPage
	if totalPages == 0 {
		totalPages = 1
	}

	for page := 1; page <= totalPages; page++ {
		start := (page - 1) * perPage
		end := start + perPage
		if end > len(posts) {
			end = len(posts)
		}

		pagePosts := posts[start:end]
		postCtxs := make([]theme.PostCtx, 0, len(pagePosts))
		for _, p := range pagePosts {
			postCtxs = append(postCtxs, buildPostCtx(siteCtx, themeCtx, p))
		}

		pagCtx := buildPaginationCtx(page, totalPages, "")

		idxCtx := theme.IndexCtx{
			Site:       siteCtx,
			Theme:      themeCtx,
			Layout:     "index",
			Permalink:  pageURL("", page),
			Posts:      postCtxs,
			Pagination: pagCtx,
		}

		html, err := b.engine.Execute("index.html", idxCtx)
		if err != nil {
			return fmt.Errorf("render index page %d: %w", page, err)
		}

		outPath := pageIndexPath(distDir, page)
		if err := writeFile(outPath, html); err != nil {
			return err
		}
	}

	return nil
}

// writeFile creates the directory and writes the HTML content.
func writeFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(path), err)
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

// permalinkToPath converts a permalink to a file path segment.
func permalinkToPath(p string) string {
	return filepath.Clean(p)
}

// pageIndexPath returns the output path for a paginated index page.
func pageIndexPath(distDir string, page int) string {
	if page == 1 {
		return filepath.Join(distDir, "index.html")
	}
	return filepath.Join(distDir, "page", fmt.Sprintf("%d", page), "index.html")
}
