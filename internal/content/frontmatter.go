package content

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Frontmatter holds the parsed YAML metadata from a Markdown file.
type Frontmatter struct {
	Title      string   `yaml:"title"`
	Date       string   `yaml:"date"`
	Updated    string   `yaml:"updated"`
	Tags       []string `yaml:"tags"`
	Categories []string `yaml:"categories"`
	Excerpt    string   `yaml:"excerpt"`
	CoverImage string   `yaml:"cover_image"`
	Author     string   `yaml:"author"`
	Layout     string   `yaml:"layout"`
	Permalink  string   `yaml:"permalink"`
	Published  *bool    `yaml:"published"` // pointer to distinguish false from missing
	Comments   *bool    `yaml:"comments"`
	Lang       string   `yaml:"lang"`
	Weight     int      `yaml:"weight"`
	Featured   bool     `yaml:"featured"`
}

// delimiter is the YAML frontmatter boundary.
const delimiter = "---"

// ParseFrontmatter splits raw markdown into frontmatter and body.
// Returns the parsed Frontmatter struct and the raw markdown body.
func ParseFrontmatter(raw string) (*Frontmatter, string, error) {
	raw = strings.TrimSpace(raw)

	if !strings.HasPrefix(raw, delimiter) {
		return &Frontmatter{}, raw, nil
	}

	// find closing ---
	closingIdx := strings.Index(raw[len(delimiter):], "\n"+delimiter)
	if closingIdx < 0 {
		return nil, "", fmt.Errorf("frontmatter: missing closing %s", delimiter)
	}

	yamlSection := raw[len(delimiter) : len(delimiter)+closingIdx]
	body := raw[len(delimiter)+closingIdx+len("\n"+delimiter):]

	var fm Frontmatter
	decoder := yaml.NewDecoder(bytes.NewReader([]byte(yamlSection)))
	decoder.KnownFields(true) // strict: error on unknown fields

	if err := decoder.Decode(&fm); err != nil {
		return nil, "", fmt.Errorf("frontmatter: %w", err)
	}

	return &fm, strings.TrimSpace(body), nil
}

// IsDraft reports whether the post should be excluded from the build.
func (fm *Frontmatter) IsDraft() bool {
	if fm.Published == nil {
		return false
	}
	return !*fm.Published
}

// IsCommentsEnabled reports whether comments are enabled for the post.
func (fm *Frontmatter) IsCommentsEnabled() bool {
	if fm.Comments == nil {
		return true // default enabled
	}
	return *fm.Comments
}
