package content

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Loader reads Markdown files from a content directory.
type Loader struct {
	contentDir string
}

// NewLoader creates a Loader for the given content directory.
func NewLoader(contentDir string) *Loader {
	return &Loader{contentDir: contentDir}
}

// LoadPosts reads all Markdown files under contentDir/posts/.
func (l *Loader) LoadPosts() ([]*Post, error) {
	postsDir := filepath.Join(l.contentDir, "posts")
	return l.loadDir(postsDir, LayoutPost)
}

// LoadPages reads all Markdown files under contentDir/pages/.
func (l *Loader) LoadPages() ([]*Post, error) {
	pagesDir := filepath.Join(l.contentDir, "pages")
	return l.loadDir(pagesDir, LayoutPage)
}

// loadDir recursively reads .md files and converts them to Posts.
func (l *Loader) loadDir(dir string, fallbackLayout string) ([]*Post, error) {
	var posts []*Post

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		post, parseErr := l.parseFile(path, fallbackLayout)
		if parseErr != nil {
			return fmt.Errorf("parse %s: %w", path, parseErr)
		}
		posts = append(posts, post)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walk %s: %w", dir, err)
	}

	sortByDate(posts)
	return posts, nil
}

// parseFile reads a single Markdown file and converts it to a Post.
func (l *Loader) parseFile(path string, fallbackLayout string) (*Post, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fm, body, err := ParseFrontmatter(string(data))
	if err != nil {
		return nil, err
	}

	layout := fm.Layout
	if layout == "" {
		layout = fallbackLayout
	}

	post := &Post{
		Title:      fm.Title,
		Slug:       slugFromPath(path),
		Date:       parseTime(fm.Date),
		Updated:    parseTime(fm.Updated),
		Tags:       fm.Tags,
		Categories: fm.Categories,
		Excerpt:    fm.Excerpt,
		CoverImage: fm.CoverImage,
		Author:     fm.Author,
		Layout:     layout,
		Permalink:  fm.Permalink,
		Content:    body,
		IsDraft:    fm.IsDraft(),
		Comments:   fm.IsCommentsEnabled(),
		Lang:       fm.Lang,
		Weight:     fm.Weight,
		Featured:   fm.Featured,
		FilePath:   path,
	}

	if post.Date.IsZero() {
		post.Date = time.Now()
	}

	return post, nil
}

// slugFromPath extracts the slug from a file path.
// e.g. "content/posts/2026-01-15-hello-world.md" → "hello-world"
func slugFromPath(path string) string {
	name := filepath.Base(path)
	name = strings.TrimSuffix(name, filepath.Ext(name))

	// strip leading date prefix like "2026-01-15-"
	parts := strings.SplitN(name, "-", 4)
	if len(parts) == 4 && len(parts[0]) == 4 && len(parts[1]) == 2 && len(parts[2]) == 2 {
		return parts[3]
	}
	return name
}

// parseTime parses a time string in common formats.
func parseTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}

	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04",
		"2006-01-02T15:04",
		"2006-01-02",
	}

	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

// sortByDate sorts posts newest first.
func sortByDate(posts []*Post) {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})
}
