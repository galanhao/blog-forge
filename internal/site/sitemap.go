package site

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/galanhao/blog-forge/internal/content"
	"github.com/galanhao/blog-forge/internal/theme"
)

// urlsetXML defines the sitemap structure.
type urlsetXML struct {
	XMLName xml.Name  `xml:"urlset"`
	XMLNS   string    `xml:"xmlns,attr"`
	URLs    []urlXML  `xml:"url"`
}

// urlXML represents a single URL in the sitemap.
type urlXML struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

// writeSitemap generates a sitemap.xml with all published posts and pages.
func (b *Builder) writeSitemap(distDir string, siteCtx theme.SiteCtx, posts, pages []*content.Post) error {
	urls := make([]urlXML, 0, len(posts)+len(pages)+1)

	// homepage
	urls = append(urls, urlXML{
		Loc:     siteCtx.URL + "/",
		LastMod: time.Now().Format("2006-01-02"),
	})

	// posts
	for _, p := range posts {
		lastMod := p.Date.Format("2006-01-02")
		if !p.Updated.IsZero() {
			lastMod = p.Updated.Format("2006-01-02")
		}
		urls = append(urls, urlXML{
			Loc:     siteCtx.URL + p.Permalink,
			LastMod: lastMod,
		})
	}

	// pages
	for _, p := range pages {
		urls = append(urls, urlXML{
			Loc:     siteCtx.URL + p.Permalink,
			LastMod: p.Date.Format("2006-01-02"),
		})
	}

	urlset := urlsetXML{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}

	data, err := xml.MarshalIndent(urlset, "", "  ")
	if err != nil {
		return fmt.Errorf("sitemap: marshal: %w", err)
	}

	out := fmt.Sprintf("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n%s\n", data)
	return writeFile(fmt.Sprintf("%s/sitemap.xml", distDir), out)
}
