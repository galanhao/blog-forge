// Package theme loads and manages blog themes following our specification.
// A theme is a directory containing Go templates, assets, and a theme.yml manifest.
package theme

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Theme represents a loaded blog theme.
type Theme struct {
	Name        string
	Version     string
	Author      string
	Description string
	Config      map[string]any // theme.yml config schema (merged with user overrides)
	Dir         string         // absolute path to the theme directory
}

// Manifest defines the theme.yml structure.
type Manifest struct {
	Name        string         `yaml:"name"`
	Version     string         `yaml:"version"`
	Author      string         `yaml:"author"`
	Description string         `yaml:"description"`
	Config      map[string]any `yaml:"config"`
}

// RequiredTemplates lists template files that every theme must provide.
// Each layout uses {{template "page-start"}} and {{template "page-end"}}
// from base.html, avoiding the block/define collision issue.
var RequiredTemplates = []string{
	"base.html",
	"index.html",
	"post.html",
	"page.html",
	"archive.html",
	"tag.html",
	"category.html",
	"404.html",
}

// Load reads a theme directory and validates its structure.
func Load(themeDir string) (*Theme, error) {
	abs, err := filepath.Abs(themeDir)
	if err != nil {
		return nil, fmt.Errorf("theme: resolve path: %w", err)
	}

	manifest, err := loadManifest(filepath.Join(abs, "theme.yml"))
	if err != nil {
		return nil, err
	}

	t := &Theme{
		Name:        manifest.Name,
		Version:     manifest.Version,
		Author:      manifest.Author,
		Description: manifest.Description,
		Config:      manifest.Config,
		Dir:         abs,
	}

	if err := t.validate(); err != nil {
		return nil, err
	}

	return t, nil
}

// loadManifest reads and parses theme.yml.
func loadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("theme: read manifest: %w", err)
	}

	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("theme: parse manifest: %w", err)
	}

	return &m, nil
}

// validate checks that all required templates exist.
func (t *Theme) validate() error {
	if t.Name == "" {
		return fmt.Errorf("theme: name is required in theme.yml")
	}

	for _, tmpl := range RequiredTemplates {
		p := filepath.Join(t.Dir, "layouts", tmpl)
		if _, err := os.Stat(p); err != nil {
			return fmt.Errorf("theme: missing required template: %s", tmpl)
		}
	}

	return nil
}

// LayoutPath returns the absolute path to a layout template.
func (t *Theme) LayoutPath(name string) string {
	return filepath.Join(t.Dir, "layouts", name)
}

// AssetsDir returns the absolute path to the assets directory.
func (t *Theme) AssetsDir() string {
	return filepath.Join(t.Dir, "assets")
}

// StaticDir returns the absolute path to the static directory.
func (t *Theme) StaticDir() string {
	return filepath.Join(t.Dir, "static")
}

// MergeConfig deep-merges user overrides into the theme's default config.
// User values take precedence over defaults.
func (t *Theme) MergeConfig(overrides map[string]any) {
	t.Config = deepMerge(t.Config, overrides)
}

// deepMerge recursively merges src into dst. src values win on conflict.
func deepMerge(dst, src map[string]any) map[string]any {
	result := make(map[string]any)

	for k, v := range dst {
		result[k] = v
	}

	for k, v := range src {
		if srcMap, ok := v.(map[string]any); ok {
			if dstMap, ok := result[k].(map[string]any); ok {
				result[k] = deepMerge(dstMap, srcMap)
				continue
			}
		}
		result[k] = v
	}

	return result
}
