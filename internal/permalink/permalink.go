// Package permalink computes permanent URLs for posts and pages.
// Supports the same 7 formats from blog-pilot plus custom patterns.
package permalink

import (
	"fmt"
	"path"
	"strings"
	"time"
)

// Format identifies a built-in permalink pattern.
type Format string

const (
	FormatDateFull      Format = "date-full"      // /2026/01/15/hello-world/
	FormatDateMonth     Format = "date-month"     // /2026/01/hello-world/
	FormatDateYear      Format = "date-year"      // /2026/hello-world/
	FormatPostsSlug     Format = "posts-slug"     // /posts/hello-world/
	FormatCategorySlug  Format = "category-slug"  // /tech/hello-world/
	FormatSlugOnly      Format = "slug-only"      // /hello-world/
	FormatPostID        Format = "post-id"        // /a1b2c3d4ef56/
	FormatCustom        Format = "custom"         // user-defined pattern
)

// Config holds permalink settings.
type Config struct {
	Format  Format
	Pattern string // only used when Format == FormatCustom
	Root    string // site root, e.g. "/" or "/blog/"
}

// Compute returns the full permalink path (relative to site root).
// The result always starts and ends with "/".
func Compute(cfg Config, slug string, date time.Time, categories []string, postID string) string {
	var p string

	switch cfg.Format {
	case FormatDateFull:
		p = fmt.Sprintf("/%s/%s/%s/%s/",
			date.Format("2006"), date.Format("01"), date.Format("02"), slug)

	case FormatDateMonth:
		p = fmt.Sprintf("/%s/%s/%s/",
			date.Format("2006"), date.Format("01"), slug)

	case FormatDateYear:
		p = fmt.Sprintf("/%s/%s/",
			date.Format("2006"), slug)

	case FormatPostsSlug:
		p = fmt.Sprintf("/posts/%s/", slug)

	case FormatCategorySlug:
		cat := firstCategory(categories)
		if cat == "" {
			cat = "uncategorized"
		}
		p = fmt.Sprintf("/%s/%s/", cat, slug)

	case FormatSlugOnly:
		p = fmt.Sprintf("/%s/", slug)

	case FormatPostID:
		id := postID
		if len(id) > 12 {
			id = id[:12]
		}
		p = fmt.Sprintf("/%s/", id)

	case FormatCustom:
		p = expandPattern(cfg.Pattern, slug, date, categories)

	default:
		p = fmt.Sprintf("/%s/%s/%s/%s/",
			date.Format("2006"), date.Format("01"), date.Format("02"), slug)
	}

	return joinWithPathRoot(cfg.Root, p)
}

// expandPattern replaces :year, :month, :day, :title, :category placeholders.
func expandPattern(pattern, slug string, date time.Time, categories []string) string {
	replacements := map[string]string{
		":year":    date.Format("2006"),
		":month":   date.Format("01"),
		":day":     date.Format("02"),
		":title":   slug,
		":category": firstCategory(categories),
	}

	result := pattern
	for key, val := range replacements {
		result = strings.ReplaceAll(result, key, val)
	}

	if !strings.HasPrefix(result, "/") {
		result = "/" + result
	}
	if !strings.HasSuffix(result, "/") {
		result = result + "/"
	}
	return result
}

// firstCategory returns the first top-level category, or empty string.
func firstCategory(categories []string) string {
	if len(categories) > 0 {
		return slugify(categories[0])
	}
	return ""
}

// slugify converts a category name to a URL-safe slug.
func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

// joinWithPathRoot prepends the site root to the permalink path.
func joinWithPathRoot(root, p string) string {
	if root == "" || root == "/" {
		return p
	}
	return path.Join(root, p) + "/"
}
