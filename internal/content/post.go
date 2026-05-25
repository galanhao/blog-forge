// Package content defines the Post model and provides Markdown file loading.
package content

import (
	"strings"
	"time"
)

// Layout constants define standard page layouts.
const (
	LayoutPost     = "post"
	LayoutPage     = "page"
	LayoutArchive  = "archive"
	LayoutCategory = "category"
	LayoutTag      = "tag"
	LayoutIndex    = "index"
)

// Post represents a single blog post or page.
type Post struct {
	Title      string
	Slug       string
	Date       time.Time
	Updated    time.Time
	Tags       []string
	Categories []string
	Excerpt    string
	CoverImage string
	Author     string
	Layout     string // post, page, about, links, etc.
	Permalink  string // user-defined override, empty means auto-compute
	Content    string // raw markdown
	HTML       string // rendered HTML
	IsDraft    bool
	Comments   bool
	Lang       string
	Weight     int
	Featured   bool

	// computed after loading
	FilePath string
}

// IsPage reports whether the post is a standalone page (not a blog post).
func (p *Post) IsPage() bool {
	return p.Layout == "page"
}

// IsPublished reports whether the post is ready to be included in the build.
func (p *Post) IsPublished() bool {
	return !p.IsDraft
}

// HasExcerpt reports whether the post has a manually set excerpt.
func (p *Post) HasExcerpt() bool {
	return p.Excerpt != ""
}

// TagSet returns deduplicated, trimmed tags.
func (p *Post) TagSet() []string {
	seen := make(map[string]bool)
	var result []string
	for _, t := range p.Tags {
		t = strings.TrimSpace(t)
		if t != "" && !seen[t] {
			seen[t] = true
			result = append(result, t)
		}
	}
	return result
}

// CategorySet returns deduplicated, trimmed categories.
func (p *Post) CategorySet() []string {
	seen := make(map[string]bool)
	var result []string
	for _, c := range p.Categories {
		c = strings.TrimSpace(c)
		if c != "" && !seen[c] {
			seen[c] = true
			result = append(result, c)
		}
	}
	return result
}
