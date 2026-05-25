package site

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/galanhao/blog-builder/internal/content"
	"github.com/galanhao/blog-builder/internal/theme"
)

// rssXML defines the top-level RSS structure.
type rssXML struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	AtomNS  string   `xml:"xmlns:atom,attr"`
	Channel channel  `xml:"channel"`
}

// channel represents the RSS channel.
type channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Generator   string `xml:"generator"`
	LastBuild   string `xml:"lastBuildDate"`
	AtomLink    string `xml:"atom:link"`
	Items       []rssItem `xml:"item"`
}

// rssItem represents a single RSS entry.
type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	Description string `xml:"description"`
	Content     string `xml:"content:encoded"`
	PubDate     string `xml:"pubDate"`
}

// writeRSS generates the RSS feed with the latest posts.
func (b *Builder) writeRSS(distDir string, siteCtx theme.SiteCtx, posts []*content.Post) error {
	limit := 20
	if len(posts) < limit {
		limit = len(posts)
	}

	items := make([]rssItem, 0, limit)
	for i := 0; i < limit; i++ {
		p := posts[i]
		permalink := siteCtx.URL + p.Permalink
		items = append(items, rssItem{
			Title:       p.Title,
			Link:        permalink,
			GUID:        permalink,
			Description: p.Excerpt,
			Content:     p.HTML,
			PubDate:     p.Date.Format(time.RFC1123),
		})
	}

	rss := rssXML{
		Version: "2.0",
		AtomNS:  "http://www.w3.org/2005/Atom",
		Channel: channel{
			Title:       siteCtx.Title,
			Link:        siteCtx.URL,
			Description: siteCtx.Subtitle,
			Generator:   "blog-builder",
			LastBuild:   time.Now().Format(time.RFC1123),
			AtomLink:    fmt.Sprintf(`<atom:link href="%s/index.xml" rel="self" type="application/rss+xml"/>`, siteCtx.URL),
			Items:       items,
		},
	}

	data, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		return fmt.Errorf("rss: marshal: %w", err)
	}

	out := fmt.Sprintf("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n%s\n", data)
	return writeFile(fmt.Sprintf("%s/index.xml", distDir), out)
}
