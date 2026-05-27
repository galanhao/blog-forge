package theme

// SiteCtx contains site-wide data available to all templates.
type SiteCtx struct {
	Title     string
	Subtitle  string
	Author    string
	URL       string
	Root      string
	Language  string
	Timezone  string
	Favicon   string
	Avatar    string
	Banner    string
	Nav       []NavItemCtx
	Config    map[string]any
}

// NavItemCtx is a navigation menu entry.
type NavItemCtx struct {
	Name string
	URL  string
	Icon string
}

// ThemeCtx contains theme-specific data.
type ThemeCtx struct {
	Name    string
	Version string
	Config  map[string]any
}

// PostCtx is the template context for a single post or page.
type PostCtx struct {
	Site       SiteCtx
	Theme      ThemeCtx
	Title      string
	Date       string
	Updated    string
	Slug       string
	Permalink  string
	Tags       []TagCtx
	Categories []CategoryCtx
	Excerpt    string
	Content    string // rendered HTML
	TOC        string // table of contents HTML
	CoverImage string
	Author     string
	Comments   bool
	Prev       *PostCtx
	Next       *PostCtx
	PV         int
	Layout     string
}

// TagCtx is a tag with its URL.
type TagCtx struct {
	Name  string
	Count int
	URL   string
}

// CategoryCtx is a category with its URL.
type CategoryCtx struct {
	Name  string
	Count int
	URL   string
}

// IndexCtx is the template context for list pages (home, tag, category, archive).
type IndexCtx struct {
	Site       SiteCtx
	Theme      ThemeCtx
	Title      string // page title (e.g. "标签: Go" or "归档")
	Excerpt    string // page description (optional)
	Layout     string // layout name (index, tag, category, archive)
	Permalink  string // canonical URL path for this page (e.g. "/tags/go/")
	Posts      []PostCtx
	Pagination PaginationCtx
	Tag        *TagCtx      // only set on tag pages
	Category   *CategoryCtx // only set on category pages
	Year       int          // only set on archive pages
}

// PaginationCtx holds pagination state for list pages.
type PaginationCtx struct {
	Current  int
	Total    int
	PrevURL  string
	NextURL  string
	Pages    []PageLinkCtx
}

// PageLinkCtx is a single page number in the pagination bar.
type PageLinkCtx struct {
	Number int
	URL    string
	IsCurr bool
}
