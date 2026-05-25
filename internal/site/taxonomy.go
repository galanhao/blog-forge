package site

import (
	"fmt"

	"github.com/galanhao/blog-builder/internal/content"
	"github.com/galanhao/blog-builder/internal/theme"
)

// writeArchive generates the archive page grouped by year.
func (b *Builder) writeArchive(distDir string, siteCtx theme.SiteCtx, themeCtx theme.ThemeCtx, posts []*content.Post) error {
	postCtxs := make([]theme.PostCtx, 0, len(posts))
	for _, p := range posts {
		postCtxs = append(postCtxs, buildPostCtx(siteCtx, themeCtx, p))
	}

	idxCtx := theme.IndexCtx{
		Site:      siteCtx,
		Theme:     themeCtx,
		Layout:    "archive",
		Title:     "归档",
		Permalink: "/archives/",
		Posts:     postCtxs,
		Pagination: theme.PaginationCtx{
			Current: 1,
			Total:   1,
			Pages: []theme.PageLinkCtx{
				{Number: 1, URL: "/archives/", IsCurr: true},
			},
		},
	}

	html, err := b.engine.Execute("archive.html", idxCtx)
	if err != nil {
		return fmt.Errorf("render archive: %w", err)
	}

	return writeFile(fmt.Sprintf("%s/archives/index.html", distDir), html)
}

// writeTags generates one page per tag.
func (b *Builder) writeTags(distDir string, siteCtx theme.SiteCtx, themeCtx theme.ThemeCtx, posts []*content.Post) error {
	groups := groupByTag(posts)

	for tag, tagPosts := range groups {
		postCtxs := make([]theme.PostCtx, 0, len(tagPosts))
		for _, p := range tagPosts {
			postCtxs = append(postCtxs, buildPostCtx(siteCtx, themeCtx, p))
		}

		tagCtx := theme.TagCtx{
			Name:  tag,
			Count: len(tagPosts),
			URL:   "/tags/" + slugify(tag) + "/",
		}

		idxCtx := theme.IndexCtx{
			Site:      siteCtx,
			Theme:     themeCtx,
			Layout:    "tag",
			Title:     "标签: " + tag,
			Permalink: "/tags/" + slugify(tag) + "/",
			Posts:     postCtxs,
			Tag:       &tagCtx,
			Pagination: theme.PaginationCtx{Current: 1, Total: 1},
		}

		html, err := b.engine.Execute("tag.html", idxCtx)
		if err != nil {
			return fmt.Errorf("render tag %s: %w", tag, err)
		}

		outPath := fmt.Sprintf("%s/tags/%s/index.html", distDir, slugify(tag))
		if err := writeFile(outPath, html); err != nil {
			return err
		}
	}

	return nil
}

// writeCategories generates one page per category.
func (b *Builder) writeCategories(distDir string, siteCtx theme.SiteCtx, themeCtx theme.ThemeCtx, posts []*content.Post) error {
	groups := groupByCategory(posts)

	for cat, catPosts := range groups {
		postCtxs := make([]theme.PostCtx, 0, len(catPosts))
		for _, p := range catPosts {
			postCtxs = append(postCtxs, buildPostCtx(siteCtx, themeCtx, p))
		}

		catCtx := theme.CategoryCtx{
			Name:  cat,
			Count: len(catPosts),
			URL:   "/categories/" + slugify(cat) + "/",
		}

		idxCtx := theme.IndexCtx{
			Site:      siteCtx,
			Theme:     themeCtx,
			Layout:    "category",
			Title:     "分类: " + cat,
			Permalink: "/categories/" + slugify(cat) + "/",
			Posts:     postCtxs,
			Category:  &catCtx,
			Pagination: theme.PaginationCtx{Current: 1, Total: 1},
		}

		html, err := b.engine.Execute("category.html", idxCtx)
		if err != nil {
			return fmt.Errorf("render category %s: %w", cat, err)
		}

		outPath := fmt.Sprintf("%s/categories/%s/index.html", distDir, slugify(cat))
		if err := writeFile(outPath, html); err != nil {
			return err
		}
	}

	return nil
}
