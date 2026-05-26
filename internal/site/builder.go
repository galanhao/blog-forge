// Package site orchestrates the full build pipeline:
// load content → render markdown → compute permalinks → apply templates → write output.
package site

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/galanhao/blog-forge/internal/config"
	"github.com/galanhao/blog-forge/internal/content"
	"github.com/galanhao/blog-forge/internal/permalink"
	"github.com/galanhao/blog-forge/internal/render"
	"github.com/galanhao/blog-forge/internal/theme"
)

// Builder orchestrates the build pipeline.
type Builder struct {
	cfg      *config.SiteConfig
	loader   *content.Loader
	renderer *render.MarkdownRenderer
	theme    *theme.Theme
	engine   *theme.Engine
}

// New creates a Builder from site config.
func New(cfg *config.SiteConfig, themesDir string) (*Builder, error) {
	loader := content.NewLoader("content")
	renderer := render.NewMarkdownRenderer()

	themeDir := filepath.Join(themesDir, cfg.Theme)
	t, err := theme.Load(themeDir)
	if err != nil {
		return nil, fmt.Errorf("load theme: %w", err)
	}

	engine, err := theme.NewEngine(t)
	if err != nil {
		return nil, fmt.Errorf("init template engine: %w", err)
	}

	return &Builder{
		cfg:      cfg,
		loader:   loader,
		renderer: renderer,
		theme:    t,
		engine:   engine,
	}, nil
}

// Build runs the full pipeline and writes output to distDir.
func (b *Builder) Build(distDir string) error {
	posts, err := b.loader.LoadPosts()
	if err != nil {
		return fmt.Errorf("load posts: %w", err)
	}

	pages, err := b.loader.LoadPages()
	if err != nil {
		return fmt.Errorf("load pages: %w", err)
	}

	posts = filterPublished(posts)
	pages = filterPublished(pages)

	if err := renderContent(b.renderer, posts); err != nil {
		return err
	}
	if err := renderContent(b.renderer, pages); err != nil {
		return err
	}

	computePermalinks(b.cfg, posts)
	computePagePermalinks(b.cfg, pages)

	if err := os.RemoveAll(distDir); err != nil {
		return fmt.Errorf("clean dist: %w", err)
	}
	if err := os.MkdirAll(distDir, 0o755); err != nil {
		return fmt.Errorf("create dist: %w", err)
	}

	siteCtx := buildSiteCtx(b.cfg, b.theme)
	themeCtx := buildThemeCtx(b.theme)

	if err := b.writePosts(distDir, siteCtx, themeCtx, posts); err != nil {
		return err
	}

	if err := b.writePages(distDir, siteCtx, themeCtx, pages); err != nil {
		return err
	}

	if err := b.writeIndex(distDir, siteCtx, themeCtx, posts); err != nil {
		return err
	}

	if err := b.writeArchive(distDir, siteCtx, themeCtx, posts); err != nil {
		return err
	}

	if err := b.writeTags(distDir, siteCtx, themeCtx, posts); err != nil {
		return err
	}

	if err := b.writeCategories(distDir, siteCtx, themeCtx, posts); err != nil {
		return err
	}

	if err := copyAssets(b.theme, distDir); err != nil {
		return fmt.Errorf("copy assets: %w", err)
	}

	if err := b.writeRSS(distDir, siteCtx, posts); err != nil {
		return fmt.Errorf("rss: %w", err)
	}

	if err := b.writeSitemap(distDir, siteCtx, posts, pages); err != nil {
		return fmt.Errorf("sitemap: %w", err)
	}

	return nil
}

// renderContent renders markdown content to HTML for each post.
func renderContent(r *render.MarkdownRenderer, posts []*content.Post) error {
	for _, p := range posts {
		html, toc, err := r.RenderWithTOC(p.Content)
		if err != nil {
			return fmt.Errorf("render %s: %w", p.Slug, err)
		}
		p.HTML = html
		p.TOC = toc

		if !p.HasExcerpt() {
			p.Excerpt = render.AutoExcerpt(p.Content, 200)
		}
	}
	return nil
}

// computePermalinks sets the Permalink field on each post.
func computePermalinks(cfg *config.SiteConfig, posts []*content.Post) {
	pc := permalink.Config{
		Format:  permalink.Format(cfg.Permalink.Format),
		Pattern: cfg.Permalink.Pattern,
		Root:    cfg.Root,
	}
	for _, p := range posts {
		p.Permalink = permalink.Compute(pc, p.Slug, p.Date, p.CategorySet(), "")
	}
}

// computePagePermalinks sets the Permalink field on standalone pages.
// Pages use a flat /slug/ format regardless of date.
func computePagePermalinks(cfg *config.SiteConfig, pages []*content.Post) {
	for _, p := range pages {
		p.Permalink = cfg.Root + "/" + p.Slug + "/"
	}
}

// filterPublished returns only published posts.
func filterPublished(posts []*content.Post) []*content.Post {
	var result []*content.Post
	for _, p := range posts {
		if p.IsPublished() {
			result = append(result, p)
		}
	}
	return result
}
