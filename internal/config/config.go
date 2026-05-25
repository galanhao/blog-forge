// Package config loads and validates the site configuration from _config.yml.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// SiteConfig holds all site-level settings from _config.yml.
type SiteConfig struct {
	Title      string        `yaml:"title"`
	Subtitle   string        `yaml:"subtitle"`
	Author     string        `yaml:"author"`
	Language   string        `yaml:"language"`
	Timezone   string        `yaml:"timezone"`
	URL        string        `yaml:"url"`
	Root       string        `yaml:"root"`
	PerPage    int           `yaml:"per_page"`
	Theme      string        `yaml:"theme"`
	Permalink  PermalinkConf `yaml:"permalink"`
	Deploy     DeployConf    `yaml:"deploy"`
	Integrate  IntegrateConf `yaml:"integrations"`
	Nav        []NavItem     `yaml:"nav"`
}

// PermalinkConf defines how URLs are generated.
type PermalinkConf struct {
	Format  string `yaml:"format"`  // date-full, date-month, slug-only, etc.
	Pattern string `yaml:"pattern"` // custom pattern, used when format=custom
}

// DeployConf defines deployment targets.
type DeployConf struct {
	Pages  bool       `yaml:"pages"`
	Server *ServerConf `yaml:"server"`
}

// ServerConf defines rsync deployment target.
type ServerConf struct {
	Host string `yaml:"host"`
	Dir  string `yaml:"dir"`
}

// IntegrateConf holds third-party integration settings.
type IntegrateConf struct {
	Comment   CommentConf    `yaml:"comment"`
	Analytics AnalyticsConf  `yaml:"analytics"`
	PV        PVConf         `yaml:"pv"`
	Wallpaper WallpaperConf  `yaml:"wallpaper"`
}

// CommentConf configures the comment system.
type CommentConf struct {
	Provider string         `yaml:"provider"` // giscus / none
	Giscus   *GiscusConf    `yaml:"giscus"`
}

// GiscusConf holds giscus-specific settings.
type GiscusConf struct {
	Repo     string `yaml:"repo"`
	Category string `yaml:"category"`
}

// AnalyticsConf configures analytics provider.
type AnalyticsConf struct {
	Provider string `yaml:"provider"` // cloudflare / none
}

// PVConf configures page-view counting.
type PVConf struct {
	Provider string `yaml:"provider"` // custom / busuanzi / none
	Endpoint string `yaml:"endpoint"`
}

// WallpaperConf configures random wallpaper.
type WallpaperConf struct {
	Provider string `yaml:"provider"` // custom / none
	Endpoint string `yaml:"endpoint"`
}

// NavItem is a single navigation menu entry.
type NavItem struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
	Icon string `yaml:"icon"`
}

// DefaultConfig returns a SiteConfig with sensible defaults.
func DefaultConfig() *SiteConfig {
	return &SiteConfig{
		Title:    "My Blog",
		Language: "zh-CN",
		Timezone: "Asia/Shanghai",
		Root:     "/",
		PerPage:  10,
		Theme:    "default",
		Permalink: PermalinkConf{
			Format: "date-full",
		},
		Integrate: IntegrateConf{
			Comment: CommentConf{Provider: "none"},
			Analytics: AnalyticsConf{Provider: "none"},
			PV: PVConf{Provider: "none"},
			Wallpaper: WallpaperConf{Provider: "none"},
		},
	}
}

// Load reads _config.yml from the given path and returns a SiteConfig.
func Load(path string) (*SiteConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks required fields and normalizes values.
func validate(cfg *SiteConfig) error {
	if cfg.Title == "" {
		return fmt.Errorf("config: title is required")
	}

	if cfg.URL == "" {
		return fmt.Errorf("config: url is required")
	}

	validFormats := map[string]bool{
		"date-full": true, "date-month": true, "date-year": true,
		"posts-slug": true, "category-slug": true,
		"slug-only": true, "post-id": true, "custom": true,
	}
	if !validFormats[cfg.Permalink.Format] {
		return fmt.Errorf("config: invalid permalink format %q", cfg.Permalink.Format)
	}

	if cfg.Permalink.Format == "custom" && cfg.Permalink.Pattern == "" {
		return fmt.Errorf("config: permalink pattern is required when format is custom")
	}

	return nil
}
