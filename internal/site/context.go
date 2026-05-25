package site

import (
	"fmt"
	"sort"
	"strings"

	"github.com/galanhao/blog-builder/internal/config"
	"github.com/galanhao/blog-builder/internal/content"
	"github.com/galanhao/blog-builder/internal/theme"
)

// buildSiteCtx creates the site-wide template context.
func buildSiteCtx(cfg *config.SiteConfig, t *theme.Theme) theme.SiteCtx {
	nav := make([]theme.NavItemCtx, 0, len(cfg.Nav))
	for _, item := range cfg.Nav {
		nav = append(nav, theme.NavItemCtx{
			Name: item.Name,
			URL:  item.URL,
			Icon: item.Icon,
		})
	}

	return theme.SiteCtx{
		Title:    cfg.Title,
		Subtitle: cfg.Subtitle,
		Author:   cfg.Author,
		URL:      cfg.URL,
		Root:     cfg.Root,
		Language: cfg.Language,
		Timezone: cfg.Timezone,
		Nav:      nav,
	}
}

// buildThemeCtx creates the theme template context.
func buildThemeCtx(t *theme.Theme) theme.ThemeCtx {
	return theme.ThemeCtx{
		Name:    t.Name,
		Version: t.Version,
		Config:  t.Config,
	}
}

// buildPostCtx creates a template context for a single post.
func buildPostCtx(siteCtx theme.SiteCtx, themeCtx theme.ThemeCtx, p *content.Post) theme.PostCtx {
	tags := make([]theme.TagCtx, 0, len(p.Tags))
	for _, t := range p.TagSet() {
		tags = append(tags, theme.TagCtx{Name: t, URL: "/tags/" + slugify(t) + "/"})
	}

	cats := make([]theme.CategoryCtx, 0, len(p.Categories))
	for _, c := range p.CategorySet() {
		cats = append(cats, theme.CategoryCtx{Name: c, URL: "/categories/" + slugify(c) + "/"})
	}

	dateStr := ""
	if !p.Date.IsZero() {
		dateStr = p.Date.Format("2006-01-02")
	}
	updatedStr := ""
	if !p.Updated.IsZero() {
		updatedStr = p.Updated.Format("2006-01-02")
	}

	return theme.PostCtx{
		Site:       siteCtx,
		Theme:      themeCtx,
		Title:      p.Title,
		Date:       dateStr,
		Updated:    updatedStr,
		Slug:       p.Slug,
		Permalink:  p.Permalink,
		Tags:       tags,
		Categories: cats,
		Excerpt:    p.Excerpt,
		Content:    p.HTML,
		CoverImage: p.CoverImage,
		Author:     p.Author,
		Comments:   p.Comments,
		Layout:     p.Layout,
	}
}

// buildPaginationCtx creates pagination context.
func buildPaginationCtx(current, total int, baseURL string) theme.PaginationCtx {
	pages := make([]theme.PageLinkCtx, 0, total)
	for i := 1; i <= total; i++ {
		pages = append(pages, theme.PageLinkCtx{
			Number: i,
			URL:    pageURL(baseURL, i),
			IsCurr: i == current,
		})
	}

	prevURL := ""
	if current > 1 {
		prevURL = pageURL(baseURL, current-1)
	}
	nextURL := ""
	if current < total {
		nextURL = pageURL(baseURL, current+1)
	}

	return theme.PaginationCtx{
		Current: current,
		Total:   total,
		PrevURL: prevURL,
		NextURL: nextURL,
		Pages:   pages,
	}
}

// pageURL returns the URL for a given page number.
func pageURL(baseURL string, page int) string {
	if page == 1 {
		if baseURL == "" {
			return "/"
		}
		return baseURL
	}
	return fmt.Sprintf("%s/page/%d/", baseURL, page)
}

// slugify converts a string to a URL-safe slug.
func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

// groupByYear groups posts by year.
func groupByYear(posts []*content.Post) map[int][]*content.Post {
	groups := make(map[int][]*content.Post)
	for _, p := range posts {
		year := p.Date.Year()
		groups[year] = append(groups[year], p)
	}
	return groups
}

// groupByTag groups posts by tag.
func groupByTag(posts []*content.Post) map[string][]*content.Post {
	groups := make(map[string][]*content.Post)
	for _, p := range posts {
		for _, t := range p.TagSet() {
			groups[t] = append(groups[t], p)
		}
	}
	return groups
}

// groupByCategory groups posts by category.
func groupByCategory(posts []*content.Post) map[string][]*content.Post {
	groups := make(map[string][]*content.Post)
	for _, p := range posts {
		for _, c := range p.CategorySet() {
			groups[c] = append(groups[c], p)
		}
	}
	return groups
}

// sortedYears returns years in descending order.
func sortedYears(groups map[int][]*content.Post) []int {
	years := make([]int, 0, len(groups))
	for y := range groups {
		years = append(years, y)
	}
	sort.Ints(years)
	for i, j := 0, len(years)-1; i < j; i, j = i+1, j-1 {
		years[i], years[j] = years[j], years[i]
	}
	return years
}
