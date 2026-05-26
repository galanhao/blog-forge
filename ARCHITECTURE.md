# ARCHITECTURE.md

## Overview

blog-forge is a static site generator written in Go. It reads Markdown files with YAML front matter, renders them to HTML using goldmark, applies Go templates from a pluggable theme system, and writes static files ready for deployment.

## Package Dependency Graph

```
                   cmd/builder
                       │
                    internal/site ──────────────────────────┐
                   ╱    │    ╲    ╲    ╲                     │
          config  content render permalink theme            │
            │                                              │
            └──────────── theme (SiteCtx/ThemeCtx) ────────┘

Legend: → means "imports"
```

**Key rule**: `internal/site/` is the ONLY package that imports other internal packages. Leaf packages (`content`, `render`, `permalink`) never import each other.

## Data Flow

```
1. config.Load("_config.yml")           → *SiteConfig
2. site.New(cfg, themesDir)             → *Builder
   ├── content.NewLoader("content")     → *Loader
   ├── render.NewMarkdownRenderer()     → *MarkdownRenderer
   ├── theme.Load(themeDir)             → *Theme (theme.yml + layouts)
   └── theme.NewEngine(theme)           → *Engine (parsed Go templates)
3. Builder.Build(distDir):
   a. loader.LoadPosts()               → []*Post
   b. loader.LoadPages()               → []*Post
   c. filterPublished()                → drop drafts
   d. renderContent()                  → Post.HTML rendered; AutoExcerpt
   e. computePermalinks()              → Post.Permalink
   f. writePosts()                     → one HTML file per post
   g. writePages()                     → one HTML file per page
   h. writeIndex()                     → paginated home page
   i. writeArchive()                   → /archives/ page
   j. writeTags()                      → /tags/<tag>/ pages
   k. writeCategories()                → /categories/<cat>/ pages
   l. copyAssets()                     → theme/assets + static/ → dist/
   m. writeRSS()                       → dist/index.xml
   n. writeSitemap()                   → dist/sitemap.xml
```

## Template Context Model

All templates receive EITHER `PostCtx` or `IndexCtx` as their top-level dot.

**PostCtx** (for post.html, page.html):
- `.Site` SiteCtx — site metadata + nav
- `.Theme` ThemeCtx — theme name + config
- `.Title`, `.Date`, `.Updated`, `.Slug`, `.Permalink`
- `.Tags []TagCtx`, `.Categories []CategoryCtx`
- `.Excerpt`, `.Content` (rendered HTML), `.TOC` (HTML)
- `.Prev *PostCtx`, `.Next *PostCtx`

**IndexCtx** (for index.html, archive.html, tag.html, category.html):
- `.Site` SiteCtx, `.Theme` ThemeCtx
- `.Title`, `.Excerpt`, `.Layout`, `.Permalink`
- `.Posts []PostCtx`, `.Pagination PaginationCtx`
- `.Tag *TagCtx` (tag pages only), `.Category *CategoryCtx` (category pages only)

## Theme Specification

Each theme lives in `themes/<name>/` and must contain:
- `theme.yml` — name, version, author, config map
- `layouts/*.html` — Go template files (base.html + page-specific layouts)

**Composition pattern**: `base.html` defines `{{define "page-start"}}` and `{{define "page-end"}}`; child templates call these partials to wrap content.

**Built-in themes**:
- `default` — Tailwind CSS, minimal style
- `papermod` — Hugo PaperMod adaptation with dark/light/auto toggle

## Permalink Formats

Supported `permalink.format` values in `_config.yml`:
- `date-full`: `/:year/:month/:day/:slug/`
- `date`: `/:year/:month/:slug/`
- `year-month`: `/:year-:month/:slug/`
- `year`: `/:year/:slug/`
- `slug`: `/:slug/`
- `category`: `/:category/:slug/`
- `custom`: use `permalink.pattern` with `:year`, `:month`, `:day`, `:slug`, `:category` tokens
